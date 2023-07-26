package httpclient

import (
	"github.com/go-resty/resty/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func NewResty(baseURL string) *resty.Client {
	return resty.New().SetBaseURL(baseURL).OnBeforeRequest(injectOtelHeaders)
}

func injectOtelHeaders(c *resty.Client, r *resty.Request) error {
	otel.GetTextMapPropagator().Inject(r.Context(), propagation.HeaderCarrier(r.Header))
	return nil
}
