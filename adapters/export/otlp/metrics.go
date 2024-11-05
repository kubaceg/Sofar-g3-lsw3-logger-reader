package otlp

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/kubaceg/sofar_g3_lsw3_logger_reader/adapters/devices/sofar"
	grpc "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	http "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
)

const (
	appName       = "sofar.logger"
	defaultPrefix = "sofar.logger"
)

type Config struct {
	Http struct {
		Url string `yaml:"url"`
	} `yaml:"http"`
	Grpc struct {
		Url string `yaml:"url"`
	} `yaml:"grpc"`
	Prefix string `yaml:"prefix"`
}

type Service struct {
	m            metric.Meter
	prefix       string
	measurements map[string]interface{}
	reader       sdk.Reader
	exporters    []sdk.Exporter
}

func New(c *Config) (*Service, error) {
	reader := sdk.NewManualReader()
	mp := sdk.NewMeterProvider(
		sdk.WithReader(reader),
		sdk.WithResource(newResource()),
	)

	prefix := defaultPrefix
	if c.Prefix != "" {
		prefix = c.Prefix
	}

	global.SetMeterProvider(mp)
	m := global.Meter(prefix)

	exporters := make([]sdk.Exporter, 0)
	if url := c.Grpc.Url; url != "" {
		e, err := grpc.New(context.Background(), grpc.WithEndpoint(url), grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		exporters = append(exporters, e)
	}

	if url := c.Http.Url; url != "" {
		e, err := http.New(context.Background(), http.WithEndpoint(url), http.WithInsecure())
		if err != nil {
			return nil, err
		}
		exporters = append(exporters, e)
	}

	s := Service{
		m:            m,
		prefix:       prefix,
		measurements: make(map[string]interface{}),
		reader:       reader,
		exporters:    exporters,
	}

	err := s.initGauges()
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// initGauges creates Int64 gauges for all reply fields that will be read
func (s *Service) initGauges() error {
	for _, name := range sofar.GetAllRegisterNames() {
		lookup := name // creating locally scoped variable for use in callback function
		g := s.createGauge(lookup)
		_, err := s.m.RegisterCallback(
			// this function is called when a collection is triggered
			func(ctx context.Context, o metric.Observer) error {
				if v, ok := s.measurements[lookup]; ok {
					o.ObserveInt64(*g, convertToInt64(v))
				} else {
					slog.Debug(fmt.Sprintf("could not find measurement for %s\n", name))
				}
				return nil
			}, *g)
		if err != nil {
			slog.Error(fmt.Sprintf("error registering gauge callback: %s", err))
			return err
		}
	}
	return nil
}

func (s *Service) createGauge(n string) *instrument.Int64ObservableGauge {
	newGauge, _ := s.m.Int64ObservableGauge(
		appName+"."+n,
		instrument.WithUnit("1"),
	)
	return &newGauge
}

// CollectAndPushMetrics triggers the collection and export of metrics over OTLP
func (s *Service) CollectAndPushMetrics(ctx context.Context, measurements map[string]interface{}) error {
	s.measurements = measurements
	rm := metricdata.ResourceMetrics{}
	if err := s.reader.Collect(ctx, &rm); err != nil {
		return err
	}

	for _, e := range s.exporters {
		if err := e.Export(ctx, rm); err != nil {
			return err
		}
	}

	return nil
}

func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(appName),
			semconv.ServiceVersion("v0.1.0"),
		),
	)
	return r
}

func convertToInt64(v interface{}) int64 {
	switch i := v.(type) {
	case uint32:
		return int64(i)
	case uint16:
		return int64(i)
	case int16:
		return int64(i)
	default:
		fmt.Println("unexpected type encountered")
		return 0
	}
}
