package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/gantries/knife/pkg/log"
)

var logger = log.New("pkg/log")

type elasticSearchClient struct {
	client *es.Client
}

func NewElasticSearchClient(configurator Configurator) Searcher {
	cfg := es.Config{
		Addresses: configurator.Addresses(),
	}
	client, err := es.NewClient(cfg)
	if err != nil {
		client = nil
	}
	return elasticSearchClient{
		client,
	}
}

func (c elasticSearchClient) SearchRaw(ctxt context.Context, index string, query any, result any) error {
	var buf *bytes.Buffer

	if s, ok := query.(string); ok {
		buf = bytes.NewBufferString(s)
	} else if s, ok := query.(*string); ok {
		buf = bytes.NewBufferString(*s)
	} else if b, ok := query.([]byte); ok {
		buf = bytes.NewBuffer(b)
	} else if b, ok := query.(*[]byte); ok {
		buf = bytes.NewBuffer(*b)
	}

	if buf == nil {
		buf = &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(query); err != nil {
			return err
		}
	}

	res, err := c.client.Search(c.client.Search.WithContext(ctxt),
		c.client.Search.WithIndex(index), c.client.Search.WithBody(buf))
	if err != nil {
		return err
	}

	if err := json.NewDecoder(res.Body).Decode(result); err != nil {
		logger.ErrorContext(ctxt, "Unable to parse search results", "error", err)
		return err
	}

	return nil
}

type Action func(client *es.Client) (*esapi.Response, error)

func (c elasticSearchClient) exec(ctxt context.Context, action Action, result *Result) (error, Warn) {
	res, err := action(c.client)
	if err != nil {
		return err, nil
	}

	if res.IsError() {
		return fmt.Errorf("elasticsearch error: %s", res.Status()), nil
	}

	if result != nil {
		buf := bytes.NewBuffer([]byte{})
		tr := io.TeeReader(res.Body, buf)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				logger.ErrorContext(ctxt, "Unable to read response body", "error", err)
			}
		}(res.Body)

		data := bytes.NewBuffer([]byte{})
		if _, err := io.Copy(data, tr); err != nil {
			logger.ErrorContext(ctxt, "Unable to copy response body", "error", err)
			return nil, Warning(err.Error())
		}
		err := json.Unmarshal(data.Bytes(), result)
		if err != nil {
			logger.ErrorContext(ctxt, "Unable to deserialize response body", "error", err)
			return nil, Warning(err.Error())
		}
	}
	return nil, nil
}

func (c elasticSearchClient) Create(ctxt context.Context, index string, doc *Doc, result *Result) (error, Warn) {
	return c.exec(ctxt, func(client *es.Client) (*esapi.Response, error) {
		body, err := json.Marshal(doc)
		if err != nil {
			return nil, err
		}
		return c.client.Create(index, *doc.GetId(), io.NopCloser(bytes.NewReader(body)),
			c.client.Create.WithContext(ctxt))
	}, result)
}

func (c elasticSearchClient) Update(ctxt context.Context, index string, doc *Doc, result *Result) (error, Warn) {
	return c.exec(ctxt, func(client *es.Client) (*esapi.Response, error) {
		body, err := json.Marshal(&UpdateDocument{Doc: *doc})
		if err != nil {
			return nil, err
		}
		return c.client.Update(index, *doc.GetId(), io.NopCloser(bytes.NewReader(body)),
			c.client.Update.WithContext(ctxt))
	}, result)
}

func (c elasticSearchClient) Delete(ctxt context.Context, index string, doc *Doc, result *Result) (error, Warn) {
	return c.exec(ctxt, func(client *es.Client) (*esapi.Response, error) {
		return c.client.Delete(index, *doc.GetId(), c.client.Delete.WithContext(ctxt))
	}, result)
}

func (c elasticSearchClient) BulkDelete(ctxt context.Context, index string, ids []*string, result *Result) (error, Warn) {
	return c.exec(ctxt, func(client *es.Client) (*esapi.Response, error) {
		var buf bytes.Buffer
		for _, id := range ids {
			meta := []byte(fmt.Sprintf(`{"delete":{ "_id" : "%s", "_index": "%s" }}%s`, *id, index, "\n"))
			buf.Write(meta)
		}
		req := esapi.BulkRequest{
			Index:   index,
			Body:    &buf,
			Refresh: "true",
		}
		res, err := req.Do(ctxt, c.client)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		if res.IsError() {
			return nil, fmt.Errorf("bulk delete request failed: %s", res.String())
		}
		return res, nil
	}, result)
}

func (c elasticSearchClient) CreateIndex(ctxt context.Context, index string, schame string, result *Result) (error, Warn) {
	return c.exec(ctxt, func(client *es.Client) (*esapi.Response, error) {
		req := esapi.IndicesCreateRequest{
			Index: index,
			Body:  io.NopCloser(bytes.NewReader([]byte(schame))),
		}
		return req.Do(ctxt, client)
	}, result)
}
