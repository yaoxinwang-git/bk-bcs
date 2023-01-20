/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package config

import (
	//"github.com/Tencent/bk-bcs/bcs-common/pkg/otel/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.5.0"
	//semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

// TracingConf web 相关配置
type TracingConf struct {
	TracingSwitch string `yaml:"tracingSwitch" usage:"tracing switch"`
	TracingType   string `yaml:"tracingType" usage:"tracing type(default jaeger)"`

	ServiceName string `yaml:"serviceName" usage:"tracing serviceName"`

	ExporterURL string `yaml:"exporterURL" usage:"url of exporter"`

	ResourceAttrs []attribute.KeyValue `yaml:"resourceAttrs" usage:"attributes of traced service"`
}

// init 初始化
func InitTracingInstance(c *TracingConf) (*sdktrace.TracerProvider, error) {
	//opts := []trace.Option{}
	//if c.TracingSwitch != "" {
	//	opts = append(opts, trace.TracerSwitch(c.TracingSwitch))
	//}
	//if c.TracingType != "" {
	//	opts = append(opts, trace.TracerType(c.TracingType))
	//}
	//
	//if c.ExporterURL != "" {
	//	opts = append(opts, trace.ExporterURL(c.ExporterURL))
	//}
	//tracer, err := trace.InitTracerProvider(c.ServiceName, opts...)
	url := "http://192.168.37.160:14268/api/traces"
	// 创建导出器
	exp, err := newExporter(url)
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exp),
		//sdktrace.WithResource(newResource()),
	)
	if err != nil {
		return nil, err
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}

// defaultTracingConf 默认配置
func defaultTracingConf() *TracingConf {
	values := make([]attribute.KeyValue, 0)
	c := &TracingConf{
		TracingSwitch: "on",
		TracingType:   "jaeger",
		ServiceName:   "bcs-monitor",
		ExporterURL:   "http://192.168.37.160:14268/api/traces",
		ResourceAttrs: values,
	}
	return c
}
func newExporter(url string) (*jaeger.Exporter, error) {
	return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
}
func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("fib"),
			semconv.ServiceVersionKey.String("v0.1.0"),
			attribute.String("environment", "demo"),
		),
	)
	return r
}
