package boilerplate

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/bridges/otellogrus"
	"go.opentelemetry.io/otel"
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
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func setupOtel(ctx context.Context, conf OtelConfig, serviceName string) (shutdown func(context.Context) error, err error) {

	var shutdownFuncs []func(context.Context) error

	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(prop)

	if conf.Metrics.Enabled {
		var meterProvider *sdkmetric.MeterProvider
		meterProvider, err = newMeterProvider(ctx, conf, serviceName)
		if err != nil {
			handleErr(err)
			return
		}
		shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
		otel.SetMeterProvider(meterProvider)
	}

	if conf.Tracing.Enabled {
		var tracerProvider *sdktrace.TracerProvider
		tracerProvider, err = newTraceProvider(ctx, conf, serviceName)
		if err != nil {
			handleErr(err)
			return
		}
		shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
		otel.SetTracerProvider(tracerProvider)
	}

	if conf.Logging.Enabled {
		var loggerProvider *sdklog.LoggerProvider
		loggerProvider, err = newLoggerProvider(ctx, conf, serviceName)
		if err != nil {
			handleErr(err)
			return
		}
		shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
		global.SetLoggerProvider(loggerProvider)
		hook := otellogrus.NewHook(conf.LoggerName, otellogrus.WithLoggerProvider(loggerProvider))
		logrus.AddHook(hook)
	}

	return
}

func defaultResource(serviceName string) (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			// TODO: inject the service version
			// semconv.ServiceVersion(serviceVersion)
		),
	)
}

func newTraceProvider(ctx context.Context, conf OtelConfig, serviceName string) (*sdktrace.TracerProvider, error) {
	traceExporter, err := newTraceExporter(ctx, conf)
	if err != nil {
		return nil, err
	}

	resource, err := defaultResource(serviceName)
	if err != nil {
		return nil, err
	}

	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithResource(resource),
		sdktrace.WithBatcher(traceExporter,
			sdktrace.WithBatchTimeout(
				time.Duration(conf.Interval)*time.Second)),
	)
	return traceProvider, nil
}

func newTraceExporter(ctx context.Context, conf OtelConfig) (sdktrace.SpanExporter, error) {
	var exporter sdktrace.SpanExporter
	var err error

	switch conf.TracingProtocol() {

	case "http", "https":
		var options []otlptracehttp.Option

		options = append(options, otlptracehttp.WithEndpoint(conf.TracingAddr()))

		if conf.TracingInsecure() {
			options = append(options, otlptracehttp.WithInsecure())
		}

		exporter, err = otlptracehttp.New(ctx, options...)

	case "grpc":
		var options []otlptracegrpc.Option
		options = append(options, otlptracegrpc.WithEndpoint(conf.TracingAddr()))
		if conf.TracingInsecure() {
			options = append(options, otlptracegrpc.WithInsecure())
		}
		exporter, err = otlptracegrpc.New(ctx, options...)

	default:
		exporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
	}

	logrus.WithContext(ctx).Debugf("otlp tracing exporter: %s, %s", conf.TracingProtocol(), conf.TracingAddr())
	return exporter, err
}

func newMeterProvider(ctx context.Context, conf OtelConfig, serviceName string) (*sdkmetric.MeterProvider, error) {
	exporter, err := newMetricExporter(ctx, conf)
	if err != nil {
		return nil, err
	}

	resource, err := defaultResource(serviceName)
	if err != nil {
		return nil, err
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter,
			sdkmetric.WithInterval(time.Duration(conf.Interval)*time.Second))),
	)

	return meterProvider, nil
}

func newMetricExporter(ctx context.Context, conf OtelConfig) (sdkmetric.Exporter, error) {
	var exporter sdkmetric.Exporter
	var err error

	switch conf.MetricsProtocol() {

	case "http", "https":
		var options []otlpmetrichttp.Option
		options = append(options, otlpmetrichttp.WithEndpoint(conf.MetricsAddr()))
		if conf.MetricsInsecure() {
			options = append(options, otlpmetrichttp.WithInsecure())
		}
		exporter, err = otlpmetrichttp.New(ctx, options...)

	case "grpc":
		var options []otlpmetricgrpc.Option
		options = append(options, otlpmetricgrpc.WithEndpoint(conf.MetricsAddr()))
		if conf.MetricsInsecure() {
			options = append(options, otlpmetricgrpc.WithInsecure())
		}
		exporter, err = otlpmetricgrpc.New(ctx, options...)

	default:
		exporter, err = stdoutmetric.New()
	}

	logrus.WithContext(ctx).Debugf("otlp metrics exporter: %s, %s", conf.MetricsProtocol(), conf.MetricsAddr())
	return exporter, err
}

func newLoggerProvider(ctx context.Context, conf OtelConfig, serviceName string) (*sdklog.LoggerProvider, error) {
	exporter, err := newLoggingExporter(ctx, conf)
	if err != nil {
		return nil, err
	}

	resource, err := defaultResource(serviceName)
	if err != nil {
		return nil, err
	}

	processor := sdklog.NewBatchProcessor(exporter, sdklog.WithExportInterval(conf.LoggingInterval()))
	loggingProvider := sdklog.NewLoggerProvider(
		sdklog.WithResource(resource),
		sdklog.WithProcessor(processor),
	)
	return loggingProvider, nil
}

func newLoggingExporter(ctx context.Context, conf OtelConfig) (sdklog.Exporter, error) {
	var exporter sdklog.Exporter
	var err error

	switch conf.LoggingProtocol() {

	case "http", "https":
		var options []otlploghttp.Option
		options = append(options, otlploghttp.WithEndpoint(conf.LoggingAddr()))
		if conf.LoggingInsecure() {
			options = append(options, otlploghttp.WithInsecure())
		}
		exporter, err = otlploghttp.New(ctx, options...)

	case "grpc":
		var options []otlploggrpc.Option
		options = append(options, otlploggrpc.WithEndpoint(conf.LoggingAddr()))
		if conf.LoggingInsecure() {
			options = append(options, otlploggrpc.WithInsecure())
		}
		exporter, err = otlploggrpc.New(ctx, options...)

	default:
		exporter, err = stdoutlog.New()
	}

	logrus.WithContext(ctx).Debugf("otlp logger exporter: %s, %s", conf.LoggingProtocol(), conf.LoggingAddr())
	return exporter, err
}
