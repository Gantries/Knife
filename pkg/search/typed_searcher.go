package search

import (
	"context"
)

type typedSearchClient[T, R interface{}] struct {
	client Searcher
}

func NewTypedSearcher[T, R interface{}](searcher Searcher) TypedSearcher[T, R] {
	return typedSearchClient[T, R]{
		client: searcher,
	}
}

func (c typedSearchClient[T, R]) Search(ctxt context.Context, index string, query any) (*Records[T, R], error) {
	result := new(Records[T, R])
	if err := c.client.SearchRaw(ctxt, index, query, result); err != nil {
		return nil, err
	}
	return result, nil
}

// ContinuousSearchByFunctionWithScriptID performs a continuous search on an index using a function score query
// with a script ID. It allows for scrolling through a large set of results and sorting them.
//
// Parameters:
//   - ctxt: The context for the search operation.
//   - index: The name of the index to search.
//   - scrollSize: The number of results to return in each batch. It will be capped at 1000 if it exceeds.
//   - source: An array of fields to include in the search results.
//   - sorts: An array of sorting criteria for the search results.
//   - scriptId: A pointer to the ID of the script to use for scoring.
//   - cond: A pointer to a map representing the boolean conditions for the query.
//   - params: A pointer to a map of parameters to pass to the script.
//   - result: A pointer to a Records object where the search results will be stored.
//
// Returns:
//   - A pointer to the Records object containing the search results.
//   - An error if the search operation fails.
func (c typedSearchClient[T, R]) ContinuousSearchByFunctionWithScriptID(ctxt context.Context,
	index string, scrollSize int64, source []string, sorts []any,
	scriptId *string, cond *map[string]interface{}, params *map[string]interface{}, result *Records[T, R]) (*Records[T, R], error) {
	if scrollSize <= 0 {
		scrollSize = 10
	}
	if scrollSize > 1000 {
		scrollSize = 1000
	}

	query := &Query[R]{
		Size:   scrollSize,
		Source: source,
		Query: FunctionScoreQuery{
			FunctionScore: FunctionScore{
				Query: BoolQuery{
					Bool: *cond,
				},
				Functions: []FunctionSyntax{
					{
						ScriptScore: ScriptScore{
							Script: Script{
								ID:     *scriptId,
								Params: *params,
							},
						},
					},
				},
				BoostMode: "replace",
				MinScore:  0.01,
			},
		},
		Sort: sorts,
	}
	if result != nil {
		if after := result.Last(); after != nil {
			query.After = after
		}
	}
	return c.Search(ctxt, index, query)
}

func (c typedSearchClient[T, R]) SearchFunctionContinuous(ctxt context.Context, index string, fn *string,
	cond *map[string]interface{}, params *map[string]interface{}, source []string, sorts []any, result *Records[T, R]) (*Records[T, R], error) {
	query := &Query[R]{
		Size:   2,
		Source: source,
		Query: FunctionScoreQuery{
			FunctionScore: FunctionScore{
				Query: BoolQuery{
					Bool: *cond,
				},
				Functions: []FunctionSyntax{
					{
						ScriptScore: ScriptScore{
							Script: Script{
								Source: *fn,
								Params: *params,
								Lang:   "painless",
							},
						},
					},
				},
				BoostMode: "replace",
				MinScore:  0.01,
			},
		},
		Sort: sorts,
	}
	if result != nil {
		if after := result.Last(); after != nil {
			query.After = after
		}
	}
	return c.Search(ctxt, index, query)
}

func (c typedSearchClient[T, R]) SearchFunction(ctxt context.Context, index string, fn *string,
	cond *map[string]interface{}, params *map[string]interface{}, source []string, sorts []any) (*Records[T, R], error) {
	return c.SearchFunctionContinuous(ctxt, index, fn, cond, params, source, sorts, nil)
}
