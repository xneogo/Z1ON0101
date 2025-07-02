package typealias

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

// 仅用于解除xlog与xtrace的循环引用

// SpanContext represents propagated span identity and state
type SpanContext = jaeger.SpanContext

func SpanFromContext(ctx context.Context) opentracing.Span {
	return opentracing.SpanFromContext(ctx)
}
