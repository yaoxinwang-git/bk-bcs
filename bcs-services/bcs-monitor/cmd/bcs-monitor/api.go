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
	"fmt"
	"github.com/Tencent/bk-bcs/bcs-common/common/version"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-monitor/pkg/api"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-monitor/pkg/config"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-monitor/pkg/discovery"
	"github.com/TencentBlueKing/bkmonitor-kits/logger"
	"github.com/oklog/run"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
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
	addrIPv6 := getIPv6AddrFromEnv(httpAddress)

	tp, err := config.InitTracingInstance(config.G.TracingConf)
	if err != nil {
		logger.Errorf("initTracingInstance failed: %v", err.Error())
	}

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logger.Fatal(err)
		}
	}()
	otel.SetTracerProvider(tp)
	newCtx, span := otel.Tracer("bcs-monitor1").Start(context.Background(), "bcs-monitor-api")
	//ctx := gin.Context{}
	//span := trace.SpanFromContext(ctx)
	//span.SetStatus(codes.Error, "operationThatCouldFail failed")
	//span.SetAttributes(attribute.Bool("isTrue", true), attribute.String("stringAttr", "hi!"))
	defer span.End()
	server, err := api.NewAPIServer(newCtx, httpAddress, addrIPv6)
	if err != nil {
		return errors.Wrap(err, "apiserver")
	}

	sdName := fmt.Sprintf("%s-%s", appName, "api")
	sd, err := discovery.NewServiceDiscovery(newCtx, sdName, version.BcsVersion, httpAddress, advertiseAddress, addrIPv6)
	if err != nil {
		return err
	}

	// 启动 apiserver
	g.Add(server.Run, func(err error) { server.Close() })
	g.Add(sd.Run, func(error) {})
	span.End()
	return nil
}
