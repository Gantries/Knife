package tel

import (
	"context"

	"github.com/gantries/knife/pkg/maps"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var meter metric.Meter

var (
	counts     = maps.Map[string, SimpleCounter]{}
	histograms = maps.Map[string, SimpleHistogram]{}
	gauges     = maps.Map[string, SimpleGauge]{}
)

func newMeter(name string) {
	meter = otel.Meter(name)
}

type SimpleCounter interface {
	Add(ctx context.Context, incr int64, options ...metric.AddOption) SimpleCounter
	Error() error
}

type SimpleHistogram interface {
	Record(ctx context.Context, incr float64, options ...metric.RecordOption) SimpleHistogram
	Error() error
}

type SimpleGauge interface {
	Error() error
}

func Counter(name string) SimpleCounter {
	if counts.Has(name) {
		return *(counts.Get(name))
	}
	s := innerCounter{}
	c, err := meter.Int64Counter(name)
	if err != nil {
		s.err = err
	} else {
		s.Instance = c
	}
	counts.Put(name, s)
	return s
}

func Histogram(name string) SimpleHistogram {
	if histograms.Has(name) {
		return *(histograms.Get(name))
	}

	s := innerHistogram{}

	h, err := meter.Float64Histogram(name)
	if err != nil {
		s.err = err
	} else {
		s.Instance = h
	}

	histograms.Put(name, s)
	return s
}

func Gauge(name string) SimpleGauge {
	if gauges.Has(name) {
		return *(gauges.Get(name))
	}

	i := innerGauge{}

	g, err := meter.Int64ObservableGauge(name)
	if err != nil {
		i.err = err
	} else {
		i.Instance = g
	}
	gauges.Put(name, i)
	return i
}
