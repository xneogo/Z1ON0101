/*
 *  ┏┓      ┏┓
 *┏━┛┻━━━━━━┛┻┓
 *┃　　　━　　  ┃
 *┃   ┳┛ ┗┳   ┃
 *┃           ┃
 *┃     ┻     ┃
 *┗━━━┓     ┏━┛
 *　　 ┃　　　┃神兽保佑
 *　　 ┃　　　┃代码无BUG！
 *　　 ┃　　　┗━━━┓
 *　　 ┃         ┣┓
 *　　 ┃         ┏┛
 *　　 ┗━┓┓┏━━┳┓┏┛
 *　　   ┃┫┫  ┃┫┫
 *      ┗┻┛　 ┗┻┛
 @Time    : 2024/10/12 -- 15:30
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: tracer.go
*/

package xtrace

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/xneogo/Z1ON0101/xconfig/xapollo"
	"github.com/xneogo/Z1ON0101/xlog"
	"github.com/xneogo/Z1ON0101/xtrace/typealias"
	"github.com/xneogo/matrix/mconfig/mobserver"
)

// TracerType denotes the underlining implementation of opentracing-compatible tracer
type TracerType string

// Tracer is a simple, thin interface for Span creation and SpanContext
// propagation
type Tracer = opentracing.Tracer

// SpanContext represents propagated span identity and state
type SpanContext = typealias.SpanContext

// StartSpanOption instances (zero or more) may be passed to Tracer.StartSpan.
//
// StartSpanOption borrows from the "functional options" pattern, per
// http://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
type StartSpanOption = opentracing.StartSpanOption

// Span represents an active, un-finished span in the OpenTracing system.
//
// Spans are created by the Tracer interface.
type Span = opentracing.Span

// TracerTypeJaeger identity the Jaeger's tracer implementation
const TracerTypeJaeger TracerType = "jaeger"

// InitDefaultTracer will initialize the default tracer, which is now the Jaeger tracer.
func InitDefaultTracer(serviceName string) error {
	return InitTracer(TracerTypeJaeger, serviceName)
}

// InitTracer provides a way of initialize a customized tracer, which support only the Jaeger tracer currently
func InitTracer(tracerType TracerType, serviceName string) error {
	if bt != nil {
		return nil
	}

	gServiceName = serviceName
	switch tracerType {
	case TracerTypeJaeger:
		return initJaeger(serviceName)
	default:
		return fmt.Errorf("unknown tracer type %v", tracerType)
	}
}

// CloseTracer stop a tracer from collecting trace information, usually this function
//   should be invoked in an graceful exit/handling.
func CloseTracer() error {
	if bt != nil {
		bt.Close()
	}
	return nil
}

func initJaeger(serviceName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := initApolloCenter(ctx); err != nil {
		return err
	}

	cfg, err := getJaegerTracerConfiguration(ctx, serviceName)
	if err != nil {
		return err
	}

	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		return err
	}

	bt = newBackTracer(serviceName)
	bt.InitTracer(tracer, closer)

	observer := mobserver.NewConfigObserver(bt.handleJaegerSamplerChangeEvent)
	// 注意这里不能传入带cancel的ctx, 否则observer会被关闭
	apolloCenter.RegisterObserver(context.Background(), observer)

	xlog.Infof(ctx, "init jaeger for %s [done]", serviceName)
	return nil
}

func getJaegerTracerConfiguration(ctx context.Context, serviceName string) (*config.Configuration, error) {
	var configManager TracerConfigManager

	configManager = newApolloTracerConfigManager()
	tracerConfig := configManager.GetConfig(serviceName, TracerTypeJaeger)

	cfg, ok := tracerConfig.Payload.(config.Configuration)
	if !ok {
		return nil, fmt.Errorf("imcompatible tracer config %v for jaeger", cfg)
	}

	return &cfg, nil
}

var bt *backTracer

type backTracer struct {
	serviceName string
	tracer      opentracing.Tracer
	closer      io.Closer

	mu sync.Mutex
}

func newBackTracer(serviceName string) *backTracer {
	return &backTracer{
		serviceName: serviceName,
	}
}

func (b *backTracer) InitTracer(tracer opentracing.Tracer, closer io.Closer) {
	fun := "backTracer.InitTracer --> "
	b.mu.Lock()
	defer b.mu.Unlock()

	// TODO: 这里是否存在并发安全问题?
	opentracing.SetGlobalTracer(tracer)

	originCloser := b.closer
	b.tracer = tracer
	b.closer = closer

	go func() {
		if originCloser != nil {
			if err := originCloser.Close(); err != nil {
				xlog.Errorf(context.Background(), "%s async close origin tracer error: %v", fun, err)
			} else {
				xlog.Infof(context.Background(), "%s async close origin tracer success", fun)
			}
		}
	}()
}

func (b *backTracer) Close() {
	f := "backTracer.Close --> "
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closer != nil {
		if err := b.closer.Close(); err != nil {
			xlog.Errorf(context.Background(), "%s close closer error: %v", f, err)
		} else {
			xlog.Infof(context.Background(), "%s close closer success", f)
		}
	}
}

func (b *backTracer) handleJaegerSamplerChangeEvent(ctx context.Context, event *mobserver.ChangeEvent) {
	fun := "backTracer.HandleChangeEvent --> "

	if event.Namespace != xapollo.DefaultApolloTraceNamespace {
		return
	}

	needReloadTracer := false
	for key, change := range event.Changes {
		if isMyJaegerSamplerConfigKey(key, b.serviceName) {
			xlog.Infof(ctx, "%s get key %s from apollo, old value: %s, new value: %s", fun, key, change.OldValue, change.NewValue)
			needReloadTracer = true
		}
	}

	if !needReloadTracer {
		return
	}

	cfg, err := getJaegerTracerConfiguration(ctx, b.serviceName)
	if err != nil {
		xlog.Errorf(ctx, "%s reload tracer failed, get jaeger tracer configuration error: %v", fun, err)
		return
	}

	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		xlog.Errorf(ctx, "%s reload tracer failed, new tracer error: %v", fun, err)
		return
	}

	b.InitTracer(tracer, closer)
	xlog.Infof(ctx, "%s reload tracer success", fun)
}

// String adds a string-valued key:value pair to a Span.LogFields() record
func String(key, value string) log.Field {
	return log.String(key, value)
}

// Int adds an int-valued key:value pair to a Span.LogFields() record
func Int(key string, value int) log.Field {
	return log.Int(key, value)
}

// SpanFromContext returns the `Span` previously associated with `ctx`, or
// `nil` if no such `Span` could be found.
//
// NOTE: context.Context != SpanContext: the former is Go's intra-process
// context propagation mechanism, and the latter houses OpenTracing's per-Span
// identity and baggage information.
func SpanFromContext(ctx context.Context) opentracing.Span {
	return typealias.SpanFromContext(ctx)
}

// StartSpanFromContext starts and returns a Span with `operationName`, using
// any Span found within `ctx` as a ChildOfRef. If no such parent could be
// found, StartSpanFromContext creates a root (parentless) Span.
//
// The second return value is a context.Context object built around the
// returned Span.
//
// Example usage:
//
//    SomeFunction(ctx context.Context, ...) {
//        sp, ctx := opentracing.StartSpanFromContext(ctx, "SomeFunction")
//        defer sp.Finish()
//        ...
//    }
func StartSpanFromContext(ctx context.Context, operationName string, opts ...StartSpanOption) (Span, context.Context) {
	return opentracing.StartSpanFromContext(ctx, operationName, opts...)
}

// GlobalTracer returns the global singleton `Tracer` implementation.
// Before `SetGlobalTracer()` is called, the `GlobalTracer()` is a noop
// implementation that drops all data handed to it.
func GlobalTracer() Tracer {
	return opentracing.GlobalTracer()
}

func ExtractTraceID(ctx context.Context) string {
	span := SpanFromContext(ctx)
	if span == nil {
		return ""
	}

	sctx, ok := span.Context().(typealias.SpanContext)
	if !ok {
		return ""
	}

	return sctx.TraceID().String()
}

func ExtractSpanOperationName(ctx context.Context) string {
	span := SpanFromContext(ctx)
	if span == nil {
		return ""
	}

	jspan, ok := span.(*jaeger.Span)
	if !ok {
		return ""
	}

	return jspan.OperationName()
}

var gServiceName string

func ServiceName() string {
	return gServiceName
}
