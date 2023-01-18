/**
 * @Author xinwang_yao
 * @Date 2023/1/18 9:57
 **/

package config

import (
	"github.com/Tencent/bk-bcs/bcs-common/pkg/otel/trace"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var serviceName string = "bcs-monitor"

// TracingConf 链路追踪配置
type TracingConf struct {
	TracingEnabled bool `yaml:"tracingEnabled" usage:"tracing enabled"`

	ExporterURL string `yaml:"exporterURL" usage:"url of exporter"`

	ResourceAttrs []attribute.KeyValue `yaml:"resourceAttrs" usage:"attributes of traced service"`
}

func (c *TracingConf) init() (*sdktrace.TracerProvider, error) {

	opts := []trace.Option{}
	c = defaultTracingConf()
	if c.ExporterURL != "" {
		opts = append(opts, trace.ExporterURL(c.ExporterURL))
	}
	tracer, err := trace.InitTracerProvider(serviceName, opts...)
	if err != nil {
		return nil, err
	}

	return tracer, nil
}

// defaultTracingConf 默认配置
func defaultTracingConf() *TracingConf {
	m := make([]attribute.KeyValue, 0)
	c := &TracingConf{
		TracingEnabled: true,
		ExporterURL:    "http://192.168.37.160:14268/api/traces",
		ResourceAttrs:  m,
	}
	return c
}
