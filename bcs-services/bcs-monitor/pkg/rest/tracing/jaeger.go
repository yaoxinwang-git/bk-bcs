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

package tracing

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel/attribute"
)

func NewHandlerWrapper() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now().UnixNano() / 1000000

		//tid, _ := trace.TraceIDFromHex("c1e47e4be8d44407be9dd303108c9cc4")
		//sid, _ := trace.SpanIDFromHex("c1e47e4be8d44407be9dd303108c9cc4")
		//sc := trace.NewSpanContext(trace.SpanContextConfig{
		//	TraceID:    tid,
		//	SpanID:     sid,
		//	TraceFlags: trace.FlagsSampled,
		//	Remote:     true,
		//})
		//ctx = trace.ContextWithSpanContext(ctx, sc)

		name := fmt.Sprintf("%s.%s", c.Request.RequestURI, c.Request.UserAgent())

		tracer := otel.Tracer(c.Request.UserAgent())
		commonAttrs := []attribute.KeyValue{
			attribute.String("component", "http"),
			attribute.String("method", c.Request.Method),
			attribute.String("url", c.Request.UserAgent()),
		}
		_, span := tracer.Start(c, name, trace.WithSpanKind(trace.SpanKindServer), trace.WithAttributes(commonAttrs...))
		defer span.End()

		reqData, _ := json.Marshal(c.Request.Body)

		c.Next()

		rspData, _ := json.Marshal(c.Request.Response)
		// 结束时间
		endTime := time.Now().UnixNano() / 1000000
		costTime := fmt.Sprintf("%vms", endTime-startTime)
		// 设置额外标签

		span.SetAttributes(attribute.Key("req").String(string(reqData)))
		span.SetAttributes(attribute.Key("cost_time").String(costTime))
		span.SetAttributes(attribute.Key("rsp").String(string(rspData)))
	}

}
