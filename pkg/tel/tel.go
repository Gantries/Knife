package tel

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"os"
	"strings"
	"time"

	"github.com/gantries/knife/pkg/easy"
	"github.com/gantries/knife/pkg/kube"
	"github.com/gantries/knife/pkg/lists"
	"github.com/gantries/knife/pkg/ver"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/resource"

	"go.opentelemetry.io/otel"
	ilog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
)

var loggerProvider ilog.LoggerProvider
var meterProvider *metric.MeterProvider
var tracerProvider *trace.TracerProvider
var hostname, ip, mac, ns = Fingerprint()
var builtinAttributes = lists.List[attribute.KeyValue]{
	semconv.ServiceVersion(ver.Version),
	semconv.HostIP(ip[0]),
	semconv.HostName(hostname),
	semconv.HostMac(mac[0]),
	semconv.K8SNamespaceName(ns),
}
var builtinAttributeStrings = lists.Collects(builtinAttributes, func(kv attribute.KeyValue) []string {
	if kv.Value.Type() == attribute.STRING || len(kv.Value.AsString()) > 0 {
		return []string{string(kv.Key), kv.Value.AsString()}
	}
	if b, err := json.Marshal(kv.Value.AsInterface()); err == nil {
		return []string{string(kv.Key), string(b)}
	}
	logger.Warn("Unable to marshall attribute", "key", kv.Key, "value", kv.Value)
	return []string{}
})
var builtinAttributeFlatStrings = lists.Collects(builtinAttributes, func(kv attribute.KeyValue) []string {
	if kv.Value.Type() == attribute.STRING || len(kv.Value.AsString()) > 0 {
		return []string{strings.ReplaceAll(string(kv.Key), ".", "_"), kv.Value.AsString()}
	}
	if b, err := json.Marshal(kv.Value.AsInterface()); err == nil {
		return []string{strings.ReplaceAll(string(kv.Key), ".", "_"), string(b)}
	}
	logger.Warn("Unable to marshall attribute", "key", kv.Key, "value", kv.Value)
	return []string{}
})

type Closable func(context.Context) error

// SetupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
// reference https://opentelemetry.io/docs/languages/go/getting-started/
func SetupOTelSDK(ctx context.Context, name string, config OpenTelemetry) (
	shutdown func(context.Context) error, err error) {
	var shutdownFunctions []Closable

	// shutdown calls cleanup functions registered via shutdownFunctions.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var shutdownError error
		for _, fn := range shutdownFunctions {
			e := fn(ctx)
			shutdownError = errors.Join(shutdownError, e)
		}
		shutdownFunctions = nil
		return shutdownError
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(context.Background())) // using background task to prevent from deadlock
	}

	res, err := resource.Merge(resource.Default(), resource.NewWithAttributes(semconv.SchemaURL,
		append(BuiltinAttributes(), semconv.ServiceName(name))...),
	)
	easy.Panic(err, handleErr)

	newTracer(name)
	newMeter(name)

	// Setup propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// Setup trace provider.
	tracerProvider, err = newTracerProvider(newTraceExport(ctx, config.Trace), res)
	easy.Panic(err, handleErr)
	shutdownFunctions = append(shutdownFunctions, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// Setup meter provider.
	meterProvider, err = newMeterProvider(newMetricExport(ctx, config.Metric), res)
	easy.Panic(err, handleErr)
	shutdownFunctions = append(shutdownFunctions, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	// Setup logger provider
	p, err := SetupOTelLog(ctx, res, config.Log, map[string]string{})
	easy.Panic(err, handleErr)
	shutdownFunctions = append(shutdownFunctions, p.Shutdown)
	loggerProvider = p

	// add necessary instrumentation
	easy.Panic(host.Start(host.WithMeterProvider(meterProvider)), handleErr)
	easy.Panic(runtime.Start(runtime.WithMeterProvider(meterProvider)), handleErr)

	return
}

func SetupOTelLog(ctx context.Context, res *resource.Resource, remote Remote, headers map[string]string) (l *log.LoggerProvider, err error) {
	l, err = newLoggerProvider(newLoggerExporter(ctx, remote, headers), res)
	if err != nil {
		return nil, err
	}
	global.SetLoggerProvider(l)
	SetupLoggersCreated(l)
	return l, nil
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTracerProvider(traceExportGenerator func() (trace.SpanExporter, error), res *resource.Resource) (*trace.TracerProvider, error) {
	tracerExporter, err := traceExportGenerator()
	if err != nil {
		return nil, err
	}

	return trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithBatcher(tracerExporter,
			// Default is 5s. Set to 1s for demonstrative purposes.
			trace.WithBatchTimeout(5*time.Second)),
	), nil
}

func newMeterProvider(metricExportGenerator func() (metric.Exporter, error), res *resource.Resource) (*metric.MeterProvider, error) {
	metricExporter, err := metricExportGenerator()
	if err != nil {
		return nil, err
	}

	return metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(metricExporter,
			// Default is 1m. Set to 3s for demonstrative purposes.
			metric.WithInterval(3*time.Second))),
	), nil
}

func newLoggerProvider(loggerExporterBuilder func() (log.Exporter, error), res *resource.Resource) (*log.LoggerProvider, error) {
	exporter, err := loggerExporterBuilder()
	if err != nil {
		return nil, err
	}
	batch := log.NewBatchProcessor(exporter)
	return log.NewLoggerProvider(log.WithResource(res), log.WithProcessor(batch)), nil
}

func newMetricExport(ctx context.Context, config Remote) func() (metric.Exporter, error) {
	return func() (metric.Exporter, error) {
		switch config.Protocol {
		case HTTP:
			return otlpmetrichttp.New(ctx, otlpmetrichttp.WithEndpoint(config.Endpoint),
				otlpmetrichttp.WithInsecure())
		case GRPC:
			return otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithEndpoint(config.Endpoint),
				otlpmetricgrpc.WithInsecure())
		default:
			return stdoutmetric.New(stdoutmetric.WithPrettyPrint())
		}
	}
}

func newTraceExport(ctx context.Context, config Remote) func() (trace.SpanExporter, error) {
	return func() (trace.SpanExporter, error) {
		switch config.Protocol {
		case HTTP:
			return otlptracehttp.New(ctx, otlptracehttp.WithEndpoint(config.Endpoint), otlptracehttp.WithInsecure())
		case GRPC:
			return otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint(config.Endpoint), otlptracegrpc.WithInsecure())
		default:
			return stdouttrace.New(stdouttrace.WithPrettyPrint())
		}
	}
}

func newLoggerExporter(ctx context.Context, config Remote, headers map[string]string) func() (log.Exporter, error) {
	return func() (log.Exporter, error) {
		switch config.Protocol {
		case HTTP:
			return otlploghttp.New(ctx, otlploghttp.WithEndpoint(config.Endpoint), otlploghttp.WithInsecure(), otlploghttp.WithHeaders(headers))
		case GRPC:
			return otlploggrpc.New(ctx, otlploggrpc.WithEndpoint(config.Endpoint), otlploggrpc.WithInsecure(), otlploggrpc.WithHeaders(headers))
		default:
			return stdoutlog.New(stdoutlog.WithPrettyPrint())
		}
	}
}

func Fingerprint() (hostname string, ip, mac []string, ns string) {
	hostname, _ = os.Hostname()
	ifs, err := net.Interfaces()
	if err == nil {
		for _, i := range ifs {
			if i.Flags&net.FlagLoopback != 0 || i.Flags&net.FlagUp == 0 || i.Flags&net.FlagPointToPoint != 0 {
				continue
			}
			if addrs, err := i.Addrs(); err == nil {
				mac = append(mac, i.HardwareAddr.String())
				for _, addr := range addrs {
					network := addr.Network()
					address := addr.String()
					switch network {
					case "ip+net":
						ip = append(ip, address)
					default:
						logger.Warn("Ignoring network address in telemetry", "network", network, "address", address)
					}
				}
			}
		}
	}

	ns, _ = kube.Namespace()
	return
}

func BuiltinAttributes() lists.List[attribute.KeyValue] {
	return builtinAttributes
}

func BuiltinAttributeStrings() lists.List[string] {
	return builtinAttributeStrings
}

func BuiltinAttributeFlatStrings() lists.List[string] {
	return builtinAttributeFlatStrings
}
