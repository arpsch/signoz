Added the Optel middleware for  "github.com/ant0ine/go-json-rest/rest"  lightweight go framework.

This is adopted from - https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/github.com/gin-gonic/gin/otelgin


## Run command
SERVICE_NAME=goApp INSECURE_MODE=true OTEL_EXPORTER_OTLP_ENDPOINT=<IP of SigNoz backend>:4317 go run main.go

