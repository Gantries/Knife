package search

const (
	_ = iota
	ElasticSearch
	OpenSearch
)

type Configurator interface {
	Addresses() []string
	SearchEngine() int
}
