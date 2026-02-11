package search

import "context"

type Searcher interface {
	SearchRaw(ctxt context.Context, index string, query any, result any) error
	Create(ctxt context.Context, index string, doc *Doc, result *Result) (error, Warn)
	Update(ctxt context.Context, index string, doc *Doc, result *Result) (error, Warn)
	Delete(ctxt context.Context, index string, doc *Doc, result *Result) (error, Warn)
	BulkDelete(ctxt context.Context, index string, ids []*string, result *Result) (error, Warn)
	CreateIndex(ctxt context.Context, index string, schema string, result *Result) (error, Warn)
}

type TypedSearcher[T interface{}, R interface{}] interface {
	Search(ctxt context.Context, index string, query any) (*Records[T, R], error)
	SearchFunction(context.Context, string, *string, *map[string]interface{}, *map[string]interface{}, []string, []any) (*Records[T, R], error)
	SearchFunctionContinuous(context.Context, string, *string, *map[string]interface{}, *map[string]interface{}, []string, []any, *Records[T, R]) (*Records[T, R], error)
	ContinuousSearchByFunctionWithScriptID(ctxt context.Context, index string, scrollSize int64, source []string, sorts []any, scriptId *string, cond *map[string]interface{}, params *map[string]interface{}, result *Records[T, R]) (*Records[T, R], error)
}
