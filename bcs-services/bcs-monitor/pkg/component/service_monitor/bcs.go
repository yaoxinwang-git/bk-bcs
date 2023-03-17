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

// Package service_monitor xxx
package service_monitor

import (
	"fmt"
	"k8s.io/klog/v2"
	"time"

	"github.com/Tencent/bk-bcs/bcs-common/pkg/bcsapi/clustermanager"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-monitor/pkg/component"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-monitor/pkg/config"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-monitor/pkg/storage"
)

// Cluster 集群信息
type Cluster struct {
	ProjectId   string `json:"projectID"`
	ClusterId   string `json:"clusterID"`
	ClusterName string `json:"clusterName"`
	BKBizID     string `json:"businessID"`
	Status      string `json:"status"`
	IsShared    bool   `json:"is_shared"`
}

// ServiceMonitor
type ServiceMonitor struct {
	ApiVersion  string    `json:"apiVersion"`
	ClusterId   string    `json:"clusterID"`
	ClusterName string    `json:"clusterName"`
	CreateTime  time.Time `json:"create_time"`
	Environment string    `json:"environment"`
	IAMNSId     string    `json:"iam_ns_id"`
	InstanceId  string    `json:"instance_id"`
	IsSystem    bool      `json:"is_system"`
	Kind        string    `json:"kind"`
	Metadata    *Metadata `json:"metadata"`
	Name        string    `json:"name"`
	NameSpace   string    `json:"namespace"`
	NameSpaceId string    `json:"namespace_id"`
	ServiceName string    `json:"service_name"`

	ServiceMonitorSpec *ServiceMonitorSpec `json:"spec,omitempty"` // 类型是 NamespaceList 等有值
}

// ServiceMonitorSpec defines the desired state of ServiceMonitor
type ServiceMonitorSpec struct {
	Endpoints   []Endpoint    `json:"endpoints,omitempty"`
	SampleLimit uint32        `json:"sampleLimit,omitempty"`
	Selector    LabelSelector `json:"selector,omitempty"`
}

// Metadata :
type Metadata struct {
	Labels      Set    `json:"labels"`
	Annotations Set    `json:"annotations"`
	Name        string `json:"name"`
	NameSpace   string `json:"namespace"`
	ServiceName string `json:"service_name"`
}

//// Item :
//type Item struct {
//	Metadata *Metadata `json:"metadata"`
//}

// Set is a map of label:value. It implements Labels.
// https://github.com/kubernetes/apimachinery/blob/master/pkg/labels/labels.go
type Set map[string]string

// String :
func (c *Cluster) String() string {
	return fmt.Sprintf("cluster<%s, %s>", c.ClusterName, c.ClusterId)
}

//// CacheListClusters 定时同步 cluster 列表
//func CacheListClusters() {
//	go func() {
//		ListClusters()
//		for range time.Tick(time.Minute * 10) {
//			klog.Infof("list clusters running")
//			ListClusters()
//			klog.Infof("list clusters end")
//		}
//	}()
//}

const listClustersCacheKey = "bcs.ListClusters"
const listNameSpacesCacheKey = "bcs.ListNameSpaces"

// ListClusters 获取集群列表
func ListClusters(projectId, limit, offset string) (map[string]*Cluster, error) {
	url := fmt.Sprintf("%s/bcsapi/v4/clustermanager/v1/cluster", config.G.BCS.Host)

	queryParams := map[string]string{
		"projectId": projectId,
		"limit":     limit,
		"offset":    offset,
	}
	resp, err := component.GetClient().R().
		SetAuthToken(config.G.BCS.Token).
		SetQueryParams(queryParams).
		Get(url)

	if err != nil {
		klog.Infof("list clusters error, %s", err.Error())
		return nil, err
	}

	var result []*Cluster
	if err := component.UnmarshalBKResult(resp, &result); err != nil {
		klog.Infof("unmarshal clusters error, %s", err.Error())
		return nil, err
	}

	clusterMap := map[string]*Cluster{}
	for _, cluster := range result {
		// 过滤掉共享集群
		if cluster.IsShared {
			continue
		}
		// 集群状态 https://github.com/Tencent/bk-bcs/blob/master/bcs-services/bcs-cluster-manager/api/clustermanager/clustermanager.proto#L1003
		if cluster.Status != "RUNNING" {
			continue
		}
		clusterMap[cluster.ClusterId] = cluster
	}

	storage.LocalCache.Slot.Set(listClustersCacheKey, clusterMap, -1)
	return clusterMap, nil
}

// GetClusterMap 获取全部集群数据, map格式
func GetClusterMap() (map[string]*Cluster, error) {
	if cacheResult, ok := storage.LocalCache.Slot.Get(listClustersCacheKey); ok {
		return cacheResult.(map[string]*Cluster), nil
	}
	return nil, fmt.Errorf("not found clusters")
}

// GetCluster 获取集群详情
func GetCluster(clusterID string) (*Cluster, error) {
	var cacheResult interface{}
	var ok bool
	if cacheResult, ok = storage.LocalCache.Slot.Get(listClustersCacheKey); !ok {
		return nil, fmt.Errorf("not found cluster")
	}
	if clusterMap, ok := cacheResult.(map[string]*Cluster); ok {
		return clusterMap[clusterID], nil
	}
	return nil, fmt.Errorf("cluster cache is invalid")
}

// ListNameSpaces 获取NameSpace列表
func ListNameSpaces(projectId, limit, offset string) (map[string]*clustermanager.Namespace, error) {
	url := fmt.Sprintf("%s/bcsapi/v4/clustermanager/v1/namespace", config.G.BCS.Host)

	queryParams := map[string]string{
		"projectId": projectId,
		//"with_lb":   limit,
		"limit":  limit,
		"offset": offset,
	}
	resp, err := component.GetClient().R().
		SetAuthToken(config.G.BCS.Token).
		SetQueryParams(queryParams).
		Get(url)

	if err != nil {
		klog.Infof("list namespace error, %s", err.Error())
		return nil, err
	}

	var result []*clustermanager.Namespace
	if err := component.UnmarshalBKResult(resp, &result); err != nil {
		klog.Infof("unmarshal namespace error, %s", err.Error())
		return nil, err
	}

	namespaceMap := map[string]*clustermanager.Namespace{}
	for _, namespace := range result {
		// namespaces = getitems(resp, 'data.results', []) or []
		// return {(i['cluster_id'], i['name']): i['id'] for i in namespaces}
		namespaceMap[namespace.Name] = namespace
	}

	storage.LocalCache.Slot.Set(listNameSpacesCacheKey, namespaceMap, -1)
	return namespaceMap, nil
}

// ListClusters 获取集群列表
func ListServiceMonitor(namespace string) {
	var url string
	if namespace == "" {
		url = fmt.Sprintf("%s/apis/monitoring.coreos.com/v1/namespaces/%s/servicemonitors", config.G.BCS.Host, namespace)
	} else {
		url = fmt.Sprintf("%s/apis/monitoring.coreos.com/v1/servicemonitors", config.G.BCS.Host)
	}

	resp, err := component.GetClient().R().
		SetAuthToken(config.G.BCS.Token).
		Get(url)

	if err != nil {
		klog.Infof("list servicemonitors error, %s", err.Error())
		return
	}

	var result []*Cluster
	if err := component.UnmarshalBKResult(resp, &result); err != nil {
		klog.Infof("unmarshal servicemonitors error, %s", err.Error())
		return
	}

	clusterMap := map[string]*Cluster{}
	for _, cluster := range result {
		// 过滤掉共享集群
		if cluster.IsShared {
			continue
		}
		// 集群状态 https://github.com/Tencent/bk-bcs/blob/master/bcs-services/bcs-cluster-manager/api/clustermanager/clustermanager.proto#L1003
		if cluster.Status != "RUNNING" {
			continue
		}
		clusterMap[cluster.ClusterId] = cluster
	}

	storage.LocalCache.Slot.Set(listClustersCacheKey, clusterMap, -1)
	return
}
