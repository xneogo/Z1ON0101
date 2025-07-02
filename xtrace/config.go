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
 @Time    : 2024/10/18 -- 16:22
 @Author  : 亓官竹 ❤️ MONEY
 @Copyright 2024 亓官竹
 @Description: config.go
*/

package xtrace

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/xneogo/Z1ON0101/xconfig/xapollo"
	"github.com/xneogo/Z1ON0101/xlog"
)

const (
	defaultAgentHost = "127.0.0.1"
	defaultAgentPort = "6831"
)

// Environment variables for jaeger agent
const (
	EnvJaegerAgentHost = "JAEGER_AGENT_HOST"
	EnvJaegerAgentPort = "JAEGER_AGENT_PORT"
)

const (
	apolloKeyJaegerSamplerTypePrefix  = "jaeger_sampler_type."
	apolloKeyJaegerSamplerParamPrefix = "jaeger_sampler_param."
	apolloJaegerSamplerGlobalName     = "default_global"

	defaultValueJaegerSamplerType          = jaeger.SamplerTypeRateLimiting
	defaultValueJaegerSamplerParam float64 = 1

	defaultGetApolloTimeout = 3 * time.Second
)

// TracerConfig keeps metadata for tracer
type TracerConfig struct {
	Payload interface{}
}

// TracerConfigManager defines the interface for an concrete implementation of config manager.
type TracerConfigManager interface {
	GetConfig(serviceName string, tracerType TracerType) TracerConfig
}

func newTracerConfigManager() TracerConfigManager {
	return newSimpleTracerConfigManager()
}

type simpleManager struct{}

type apolloManager struct {
}

func newSimpleTracerConfigManager() *simpleManager {
	return &simpleManager{}
}

func newApolloTracerConfigManager() *apolloManager {
	return &apolloManager{}
}

func (s *simpleManager) GetConfig(serviceName string, tracerType TracerType) TracerConfig {
	if tracerType != TracerTypeJaeger {
		xlog.Panicf(context.Background(), "unknown tracer type %s for simpleManager", tracerType)
	}

	agentHost, agentPort := defaultAgentHost, defaultAgentPort

	if h, ok := os.LookupEnv(EnvJaegerAgentHost); ok {
		agentHost = h
	}

	if p, ok := os.LookupEnv(EnvJaegerAgentPort); ok {
		agentPort = p
	}

	return TracerConfig{
		Payload: config.Configuration{
			ServiceName: serviceName,
			Disabled:    false,
			RPCMetrics:  false,
			Sampler:     defaultSamplerConfig(),
			Reporter:    defaultReporterConfig(agentHost, agentPort),
			Headers:     defaultHeadersConfig(),
		},
	}
}

func (a *apolloManager) GetConfig(serviceName string, tracerType TracerType) TracerConfig {
	f := "apolloManager.GetConfig --> "
	if tracerType != TracerTypeJaeger {
		xlog.Panicf(context.Background(), "%s unknown tracer type %s for simpleManager", f, tracerType)
	}

	agentHost, agentPort := defaultAgentHost, defaultAgentPort

	if h, ok := os.LookupEnv(EnvJaegerAgentHost); ok {
		agentHost = h
	}

	if p, ok := os.LookupEnv(EnvJaegerAgentPort); ok {
		agentPort = p
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultGetApolloTimeout)
	defer cancel()

	samplerCfg := a.getSamplerConfig(ctx, serviceName)

	return TracerConfig{
		Payload: config.Configuration{
			ServiceName: serviceName,
			Disabled:    false,
			RPCMetrics:  false,
			Sampler:     samplerCfg,
			Reporter:    defaultReporterConfig(agentHost, agentPort),
			Headers:     defaultHeadersConfig(),
		},
	}
}

func (a *apolloManager) getSamplerConfig(ctx context.Context, serviceName string) *config.SamplerConfig {
	f := "applloManager.getSamplerConfig --> "

	cfg, err := a.getServiceSamplerConfigFromApollo(ctx, serviceName)
	if err == nil {
		xlog.Infof(ctx, "%s get service sampler config from apollo success, name: %s, cfg: %v", f, serviceName, cfg)
		return cfg
	}
	xlog.Warnf(ctx, "%s get service sampler config from apollo error, name: %s, err: %v", f, serviceName, err)

	cfg, err = a.getGlobalSamplerConfigFromApollo(ctx)
	if err == nil {
		xlog.Infof(ctx, "%s get global sampler config from apollo success, cfg: %v", f, cfg)
		return cfg
	}
	xlog.Errorf(ctx, "%s get global sampler config from apollo error, err: %v", f, err)

	cfg = defaultSamplerConfig()
	xlog.Infof(ctx, "%s use default sampler config: %v", f, cfg)
	return cfg
}

func (a *apolloManager) getServiceSamplerConfigFromApollo(ctx context.Context, serviceName string) (*config.SamplerConfig, error) {
	return a.getSamplerConfigFromApolloByName(ctx, serviceName)
}

func (a *apolloManager) getGlobalSamplerConfigFromApollo(ctx context.Context) (*config.SamplerConfig, error) {
	return a.getSamplerConfigFromApolloByName(ctx, apolloJaegerSamplerGlobalName)
}

func (a *apolloManager) getSamplerConfigFromApolloByName(ctx context.Context, name string) (*config.SamplerConfig, error) {
	keyType, keyParam := getApolloKeysJaegerSamplerWithName(name)
	samplerTypeVal, ok := apolloCenter.GetStringWithNamespace(ctx, xapollo.DefaultApolloTraceNamespace, keyType)
	if !ok {
		return nil, fmt.Errorf("key not found: %s", keyType)
	}
	samplerParamVal, ok := apolloCenter.GetStringWithNamespace(ctx, xapollo.DefaultApolloTraceNamespace, keyParam)
	if !ok {
		return nil, fmt.Errorf("key not found: %s", keyParam)
	}
	samplerParam, err := strconv.ParseFloat(samplerParamVal, 64)
	if err != nil {
		return nil, fmt.Errorf("parse key error, key: %s, value: %s, err: %v", keyParam, samplerParamVal, err)
	}

	cfg := &config.SamplerConfig{
		Type:  samplerTypeVal,
		Param: samplerParam,
	}
	return cfg, nil
}

func getApolloKeysJaegerSamplerWithName(name string) (string, string) {
	keySamplerType := apolloKeyJaegerSamplerTypePrefix + name
	keySamplerParam := apolloKeyJaegerSamplerParamPrefix + name
	return keySamplerType, keySamplerParam
}

func isMyJaegerSamplerConfigKey(key string, serviceName string) bool {
	serviceSamplerType, serviceSamplerParam := getApolloKeysJaegerSamplerWithName(serviceName)
	globalSamplerType, globalSamplerParam := getApolloKeysJaegerSamplerWithName(apolloJaegerSamplerGlobalName)
	return key == serviceSamplerType || key == serviceSamplerParam || key == globalSamplerType || key == globalSamplerParam
}

func defaultSamplerConfig() *config.SamplerConfig {
	return &config.SamplerConfig{
		Type:  defaultValueJaegerSamplerType,
		Param: defaultValueJaegerSamplerParam,
	}
}

func defaultReporterConfig(agentHost, agentPort string) *config.ReporterConfig {
	return &config.ReporterConfig{
		LocalAgentHostPort: getAgentAddr(agentHost, agentPort),
	}
}

func getAgentAddr(agentHost, agentPort string) string {
	return fmt.Sprintf("%s:%s", agentHost, agentPort)
}

func defaultHeadersConfig() *jaeger.HeadersConfig {
	return &jaeger.HeadersConfig{
		JaegerDebugHeader:        TraceDebugHeader,
		JaegerBaggageHeader:      TraceBaggageHeader,
		TraceContextHeaderName:   TraceContextHeaderName,
		TraceBaggageHeaderPrefix: TraceBaggageHeaderPrefix,
	}
}
