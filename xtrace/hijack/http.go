package hijack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/opentracing/opentracing-go"
)

// TraceHTTPRequest is a helper function for injecting trace-related context into HTTP headers.
func TraceHTTPRequest(ctx context.Context, req *http.Request) error {
	fun := "TraceHttpRequest-->"
	if ctx == nil {
		return fmt.Errorf("%s got nil context", fun)
	}

	if req == nil {
		return fmt.Errorf("%s got nil request", fun)
	}

	if span := opentracing.SpanFromContext(ctx); span != nil {
		return opentracing.GlobalTracer().Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(req.Header))
	}
	return fmt.Errorf("%s got nil span", fun)
}
