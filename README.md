Added the Optel middleware for  "github.com/ant0ine/go-json-rest/rest"  lightweight go framework.

This is adopted from - https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/github.com/gin-gonic/gin/otelgin

SERVICE_NAME=openTelemetryTest INSECURE_MODE=true OTEL_METRICS_EXPORTER=none OTEL_EXPORTER_OTLP_ENDPOINT=10.27.15.85:4317 go run main.go
