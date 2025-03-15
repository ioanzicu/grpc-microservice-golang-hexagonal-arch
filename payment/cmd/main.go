package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ioanzicu/microservices/payment/config"
	"github.com/ioanzicu/microservices/payment/internal/adapters/db"
	"github.com/ioanzicu/microservices/payment/internal/adapters/grpc"
	"github.com/ioanzicu/microservices/payment/internal/application/core/api"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/credentials/insecure"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.11.0"
	"go.opentelemetry.io/otel/trace"
	grpcLib "google.golang.org/grpc"
)

func tracerProvider(ctx context.Context, url string) (*tracesdk.TracerProvider, error) {
	conn, err := grpcLib.NewClient(url,
		// 	// Note the use of insecure transport here. TLS is recommended in production.
		grpcLib.WithTransportCredentials(insecure.NewCredentials()),
		grpcLib.WithBlock(),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(traceExporter),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(config.GetServiceName()),
			attribute.String("environment", config.GetEnvironmentType()),
			attribute.Int64("ID", config.GetServiceID()),
		)),
	)
	return tp, nil
}

func init() {
	log.SetFormatter(customLogger{
		formatter: log.JSONFormatter{
			FieldMap: log.FieldMap{
				"msg": "message",
			}},
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

type customLogger struct {
	formatter log.JSONFormatter
}

func (l customLogger) Format(entry *log.Entry) ([]byte, error) {
	span := trace.SpanFromContext(entry.Context)
	entry.Data["trace_id"] = span.SpanContext().TraceID().String()
	entry.Data["span_id"] = span.SpanContext().SpanID().String()

	// Below injection is Just to understand what Context has
	entry.Data["Context"] = span.SpanContext()
	return l.formatter.Format(entry)
}

func main() {
	tp, err := tracerProvider(context.Background(), config.GetTracerProviderURL())
	if err != nil {
		log.Fatal(err)
	}

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}))

	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatalf("Failed to connect to database. Error: %v", err)
	}

	application := api.NewApplication(dbAdapter)
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}
