package tel

import (
	"context"
	"go.opentelemetry.io/otel/metric"
)

type Metric[T any] struct {
	Instance T
	err      error
}

type innerCounter Metric[metric.Int64Counter]
type innerHistogram Metric[metric.Float64Histogram]
type innerGauge Metric[metric.Int64ObservableGauge]

func (m innerCounter) Error() error {
	return m.err
}
func (m innerCounter) Add(ctx context.Context, incr int64, options ...metric.AddOption) SimpleCounter {
	if m.Instance != nil {
		m.Instance.Add(ctx, incr, options...)
	}
	return m
}

func (m innerHistogram) Record(ctx context.Context, incr float64, options ...metric.RecordOption) SimpleHistogram {
	if m.Instance != nil {
		m.Instance.Record(ctx, incr, options...)
	}
	return m
}
func (m innerHistogram) Error() error {
	return m.err
}

func (m innerGauge) Error() error {
	return m.err
}
