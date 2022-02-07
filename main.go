package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/arpsch/signoz/optel"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	serviceName  = os.Getenv("SERVICE_NAME")
	signozToken  = os.Getenv("SIGNOZ_ACCESS_TOKEN")
	collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	insecure     = os.Getenv("INSECURE_MODE")
)

func initTracer() func(context.Context) error {

	// headers := map[string]string{
	// 	"signoz-access-token": signozToken,
	// }

	//	secureOption := otlptracehttp.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	//	if len(insecure) > 0 {
	secureOption := otlptracehttp.WithInsecure()
	//	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			secureOption,
			otlptracehttp.WithEndpoint(collectorURL),
			//otlptracegrpc.WithHeaders(headers),
		),
	)

	if err != nil {
		log.Fatal(err)
	}
	resources, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
			attribute.String("library.language", "go"),
		),
	)
	if err != nil {
		log.Printf("Could not set resources: ", err)
	}

	otel.SetTracerProvider(
		sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sdktrace.AlwaysSample()),
			sdktrace.WithSpanProcessor(sdktrace.NewBatchSpanProcessor(exporter)),
			sdktrace.WithSyncer(exporter),
			sdktrace.WithResource(resources),
		),
	)
	return exporter.Shutdown
}

func main() {

	cleanup := initTracer()
	defer cleanup(context.Background())

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	api.Use(&optel.NewOptelMiddleware{
		Service: "test-opentelemetry",
	})
	router, err := rest.MakeRouter(
		rest.Get("/lookup/#host", func(w rest.ResponseWriter, req *rest.Request) {
			ip, err := net.LookupIP(req.PathParam("host"))
			if err != nil {
				rest.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteJson(&ip)
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)

	fmt.Println("starting the server on :8080")
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}
