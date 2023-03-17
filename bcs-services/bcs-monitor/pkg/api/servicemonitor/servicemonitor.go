/**
 * @Author xinwang_yao
 * @Date 2023/3/16 14:26
 **/

package service_monitor

import (
	"fmt"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-monitor/pkg/component/k8sclient"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-monitor/pkg/component/service_monitor"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-monitor/pkg/rest"
)

// ListServiceMonitors 获取ServiceMonitor列表数据
// @Summary ServiceMonitor列表数据
// @Tags    Metrics
// @Success 200 {string} string
// @Router  /service_monitors [get]
func ListServiceMonitors(c *rest.Context) (interface{}, error) {
	projectId := c.Param("projectId")
	clusterId := c.Param("clusterId")
	limit := c.Param("limit")
	offset := c.Param("offset")
	logQuery := &k8sclient.LogQuery{}
	if err := c.ShouldBindQuery(logQuery); err != nil {
		rest.AbortWithBadRequestError(c, err)
		return
	}

	// 下载参数
	clusterMap, err := service_monitor.ListClusters(projectId, limit, offset)
	if err != nil {
		return nil, err
	}
	namespaceMap, err := service_monitor.ListNameSpaces(projectId, limit, offset)
	if err != nil {
		return nil, err
	}
	if _, ok := clusterMap[clusterId]; !ok {
		return nil, fmt.Errorf("'集群 ID %s 不合法'", clusterId)
	}

	service_monitor.ListServiceMonitor("")
}
