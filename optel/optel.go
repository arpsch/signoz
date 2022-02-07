package optel

import (
	"fmt"

	"github.com/ant0ine/go-json-rest/rest"
	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

const (
	tracerKey  = "otel-go-contrib-tracer"
	tracerName = "TestTracer"
)

type NewRelicMiddleware struct {
	Service        string
	Name           string
	Verbose        bool
	TracerProvider oteltrace.TracerProvider
	Propagators    propagation.TextMapPropagator
}

func (mw *NewRelicMiddleware) MiddlewareFunc(handler rest.HandlerFunc) rest.HandlerFunc {
	fmt.Println("opent-telementry mw.....")
	if mw.TracerProvider == nil {
		mw.TracerProvider = otel.GetTracerProvider()
	}
	tracer := mw.TracerProvider.Tracer(
		tracerName,
		oteltrace.WithInstrumentationVersion("0.27.0"),
	)
	if mw.Propagators == nil {
		mw.Propagators = otel.GetTextMapPropagator()
	}

	return func(writer rest.ResponseWriter, request *rest.Request) {
		//	c.Set(tracerKey, tracer)

		fmt.Println("opent-telementry closure  mw.....")
		savedCtx := request.Context()
		defer func() {
			request.Request = request.WithContext(savedCtx)
		}()
		ctx := mw.Propagators.Extract(savedCtx, propagation.HeaderCarrier(request.Header))
		opts := []oteltrace.SpanStartOption{
			oteltrace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", request.Request)...),
			oteltrace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(request.Request)...),
			oteltrace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(mw.Service, request.URL.Path, request.Request)...),
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		}
		spanName := request.URL.Path
		if spanName == "" {
			spanName = fmt.Sprintf("HTTP %s route not found", request.Method)
		}
		ctx, span := tracer.Start(ctx, spanName, opts...)
		defer span.End()

		// pass the span through the request context
		request.Request = request.WithContext(ctx)

		// serve the request to the next middleware
		handler(writer, request)

		fmt.Println("opent-telementry closure  span Name ..... ", spanName)
		status := 200
		attrs := semconv.HTTPAttributesFromHTTPStatusCode(status)
		spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCode(status)
		span.SetAttributes(attrs...)
		span.SetStatus(spanStatus, spanMessage)
		/*
			if len(c.Errors) > 0 {
				span.SetAttributes(attribute.String("gin.errors", c.Errors.String()))
			}*/

		// the timer middleware keeps track of the time
		//startTime := request.Env["START_TIME"].(*time.Time)
		//mw.agent.HTTPTimer.UpdateSince(*startTime)
	}
}
