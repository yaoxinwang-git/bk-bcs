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

package main

import (
	"context"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-monitor/pkg/config"
	"github.com/TencentBlueKing/bkmonitor-kits/logger"
	"github.com/gin-gonic/gin"
	"github.com/oklog/run"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// APIServerCmd :
func APIServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Short: "Monitor api server",
	}

	cmd.Run = func(cmd *cobra.Command, args []string) {
		runCmd(cmd, runAPIServer)
	}

	cmd.Flags().StringVar(&httpAddress, "http-address", "0.0.0.0:8089", "API listen http ip")
	cmd.Flags().StringVar(&advertiseAddress, "advertise-address", "", "The IP address on which to advertise the server")

	return cmd
}

// runAPIServer apiserver 子服务
func runAPIServer(ctx context.Context, g *run.Group, opt *option) error {
	logger.Infow("listening for requests and metrics", "address", httpAddress)
	//	addrIPv6 := getIPv6AddrFromEnv(httpAddress)

	tp, err := config.InitTracingInstance(config.G.TracingConf)
	if err != nil {
		logger.Errorf("initTracingInstance failed: %v", err.Error())
	}

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logger.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	r := gin.New()
	r.Use(otelgin.Middleware("my-server"))
	//tmplName := "user"
	//tmplStr := "user {{ .name }} (id {{ .id }})\n"
	//tmpl := template.Must(template.New(tmplName).Parse(tmplStr))
	//r.SetHTMLTemplate(tmpl)
	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		name := getUser(c, id)
		//otelgin.HTML(c, http.StatusOK, tmplName, gin.H{
		//	"name": name,
		//	"id":   id,
		//})
		c.JSON(200, name)
	})
	_ = r.Run(":19999")
	//newCtx, span := otel.Tracer("bcs-monitor1").Start(context.Background(), "bcs-monitor-api")
	//ctx := gin.Context{}
	//span := trace.SpanFromContext(ctx)
	//span.SetStatus(codes.Error, "operationThatCouldFail failed")
	//span.SetAttributes(attribute.Bool("isTrue", true), attribute.String("stringAttr", "hi!"))
	//defer span.End()

	//server, err := api.NewAPIServer(newCtx, httpAddress, addrIPv6)
	//if err != nil {
	//	return errors.Wrap(err, "apiserver")
	//}
	//
	//sdName := fmt.Sprintf("%s-%s", appName, "api")
	//sd, err := discovery.NewServiceDiscovery(newCtx, sdName, version.BcsVersion, httpAddress, advertiseAddress, addrIPv6)
	//if err != nil {
	//	return err
	//}
	//// 启动 apiserver
	//g.Add(server.Run, func(err error) {
	//	server.Close()
	//})
	//g.Add(sd.Run, func(error) {})
	//span.End()
	return nil
}

var tracer = otel.Tracer("gin-server")

func getUser(c *gin.Context, id string) string {
	// Pass the built-in `context.Context` object from http.Request to OpenTelemetry APIs
	// where required. It is available from gin.Context.Request.Context()
	_, span := tracer.Start(c.Request.Context(), "getUser", oteltrace.WithAttributes(attribute.String("id", id)))
	defer span.End()
	if id == "123" {
		return "otelgin tester"
	}
	return "unknown"
}
