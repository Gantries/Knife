package search

import (
	"fmt"

	"github.com/gantries/knife/pkg/errors"
	"github.com/gantries/knife/pkg/maps"
	"github.com/gantries/knife/pkg/national"
	"github.com/gantries/knife/pkg/types"
)

type Total struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

type Hit[T interface{}, R interface{}] struct {
	Index  string  `json:"_index"`
	Id     string  `json:"_id"`
	Score  float64 `json:"_score"`
	Source T       `json:"_source"`
	Sort   *R      `json:"sort"`
}

type Hits[T interface{}, R interface{}] struct {
	Total Total       `json:"total"`
	Hits  []Hit[T, R] `json:"hits"`
}

type Records[T, R interface{}] struct {
	Hits Hits[T, R] `json:"hits"`
}

type Query[R interface{}] struct {
	Size   int64              `json:"size,omitempty"`
	Query  FunctionScoreQuery `json:"query,omitempty"`
	Source []string           `json:"_source,omitempty"`
	Sort   []any              `json:"sort,omitempty"`
	After  *R                 `json:"search_after,omitempty"`
}

type FunctionScoreQuery struct {
	FunctionScore FunctionScore `json:"function_score"`
}

type FunctionScore struct {
	Query     BoolQuery        `json:"query"`
	Functions []FunctionSyntax `json:"functions"`
	BoostMode string           `json:"boost_mode"`
	MinScore  float64          `json:"min_score"`
}

type FunctionSyntax struct {
	ScriptScore ScriptScore `json:"script_score"`
}

type BoolQuery struct {
	Bool map[string]interface{} `json:"bool"`
}

type ScriptScore struct {
	Script Script `json:"script"`
}

type Script struct {
	ID     string         `json:"id,omitempty"`
	Source string         `json:"source,omitempty"`
	Params map[string]any `json:"params"`
	Lang   string         `json:"lang,omitempty"`
}

func (r *Records[T, R]) Last() *R {
	hits := r.Hits.Hits
	if len(hits) <= 0 {
		return nil
	}
	last := len(hits) - 1
	return hits[last].Sort
}

type Shards struct {
	Total      int `yaml:"total,omitempty"`
	Failed     int `yaml:"failed,omitempty"`
	Successful int `yaml:"successful,omitempty"`
}

type Result struct {
	Shards  Shards `yaml:"_shards,omitempty"`
	Index   string `yaml:"_index,omitempty"`
	Id      string `yaml:"_id,omitempty"`
	Version int    `yaml:"_version,omitempty"`
	Seq     int    `yaml:"_seq_no,omitempty"`
	Primary int    `yaml:"_primary_term,omitempty"`
	Result  string `yaml:"result"`
}

func (r *Result) GetId() string {
	return r.Id
}

type Doc map[types.ColumnName]interface{}

// docId id for indexing document
const docId = "id"

var _builtins = maps.Map[types.ColumnName, interface{}]{
	docId: true,
}

func NewDoc(id string) *Doc {
	return &Doc{
		docId: id,
	}
}

func (d *Doc) GetId() *string {
	if v, ok := (*d)[docId].(string); ok {
		return &v
	}
	return nil
}

func (d *Doc) set(key types.ColumnName, value any) national.Sentence {
	if _, ok := (*d)[key]; ok {
		return errors.OverwriteIsForbiddenError
	} else {
		(*d)[key] = value
		return national.Ok
	}
}

func (d *Doc) SetId(id string) national.Sentence {
	return d.set(docId, id)
}

func (d *Doc) Set(key types.ColumnName, value any, builtins maps.Set[types.ColumnName]) national.Sentence {
	if _builtins.Has(key) {
		return errors.OverwriteInternalBuiltinError
	}
	if builtins.Has(key) {
		return errors.OverwriteBuiltinError
	}
	return d.set(key, value)
}

func (d *Doc) SetNX(key types.ColumnName, value any) {
	_ = d.set(key, value)
}

func (d *Doc) Del(key types.ColumnName) (any, error) {
	if old, ok := (*d)[key]; ok {
		delete(*d, key)
		return old, nil
	} else {
		return nil, fmt.Errorf("field %s not exist", key)
	}
}

type UpdateDocument struct {
	Doc Doc `json:"doc"`
}

type Warn *string

func Warning(message string) Warn {
	return &message
}

func New(c Configurator) Searcher {
	switch c.SearchEngine() {
	case ElasticSearch:
		return NewElasticSearchClient(c)
	default:
		return nil
	}
}
