package crd

import (
	resv1b1 "github.com/caicloud/clientset/pkg/apis/resource/v1beta1"
)

const (
	ClusterStatusNew           resv1b1.ClusterPhase = "New"
	ClusterStatusInstallMaster resv1b1.ClusterPhase = "InstallMaster"
	ClusterStatusInstallAddon  resv1b1.ClusterPhase = "InstallAddon"
	ClusterStatusReady         resv1b1.ClusterPhase = "Ready"
	ClusterStatusNotReady      resv1b1.ClusterPhase = "NotReady"
	ClusterStatusFailed        resv1b1.ClusterPhase = "Failed"
	ClusterStatusDeleting      resv1b1.ClusterPhase = "Deleting"
)

var defaultConfig = []Config{
	{Name: CacheNameNode, Initializer: GetNodeCacheConfig},
	{Name: CacheNameRelease, Initializer: GetReleaseCacheConfig},
	{Name: CacheNamePod, Initializer: GetPodCacheConfig},
	{Name: CacheNameClusterQuota, Initializer: GetClusterQuotaCacheConfig},
	{Name: CacheNameTenant, Initializer: GetTenantCacheConfig},
	{Name: CacheNamePartition, Initializer: GetPartitionCacheConfig},
	{Name: CacheNameStorageClass, Initializer: GetStorageClassCacheConfig},
	{Name: CacheNameLoadBalancer, Initializer: GetLoadBalancerCacheConfig},
}

func GetDefaultConfig() []Config {
	re := make([]Config, len(defaultConfig))
	for i := range defaultConfig {
		re[i].Name = defaultConfig[i].Name
		re[i].Initializer = defaultConfig[i].Initializer
	}
	return re
}
