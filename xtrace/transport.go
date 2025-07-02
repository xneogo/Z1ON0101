package xtrace

// This group of variables defines the standard HTTP header keys
//   used for tracing
const (
	TraceDebugHeader         = "trace-debug-id"
	TraceBaggageHeader       = "trace-baggage"
	TraceContextHeaderName   = "bamboo-trace-id"
	TraceBaggageHeaderPrefix = "bamboo-ctx-"
)
