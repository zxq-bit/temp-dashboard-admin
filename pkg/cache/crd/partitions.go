package crd

import (
	tntv1al "github.com/caicloud/clientset/pkg/apis/tenant/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	"github.com/caicloud/dashboard-admin/pkg/errors"
	"github.com/caicloud/dashboard-admin/pkg/kubernetes"
)

func (scc *subClusterCaches) GetPartitionCache() (*PartitionsCache, bool) {
	return scc.GetAsPartitionCache(CacheNamePartition)
}
func (scc *subClusterCaches) GetAsPartitionCache(name string) (*PartitionsCache, bool) {
	c, ok := scc.m[name]
	if ok {
		return &PartitionsCache{lwCache: c, kc: scc.kc}, true
	}
	return nil, false
}

type PartitionsCache struct {
	lwCache *ListWatchCache
	kc      kubernetes.Interface
}

func NewPartitionsCache(kc kubernetes.Interface) (*PartitionsCache, error) {
	listWatcher, objType := GetPartitionCacheConfig(kc)
	c, e := NewListWatchCache(listWatcher, objType)
	if e != nil {
		return nil, e
	}
	return &PartitionsCache{
		lwCache: c,
		kc:      kc,
	}, nil
}

func (tc *PartitionsCache) Run(stopCh chan struct{}) {
	tc.lwCache.Run(stopCh)
}

func (tc *PartitionsCache) Get(key string) (*tntv1al.Partition, error) {
	return CacheGetPartition(key, tc.lwCache.indexer, tc.kc)
}
func (tc *PartitionsCache) List() ([]tntv1al.Partition, error) {
	return CacheListPartitions(tc.lwCache.indexer, tc.kc)
}
func (tc *PartitionsCache) ListCachePointer() (re []*tntv1al.Partition) {
	return CacheListPartitionsPointer(tc.lwCache.indexer, tc.kc)
}

func (tc *PartitionsCache) Indexes() cache.Indexer {
	return tc.lwCache.indexer
}

func GetPartitionCacheConfig(kc kubernetes.Interface) (cache.ListerWatcher, runtime.Object) {
	return &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			options.FieldSelector = fields.Everything().String()
			return kc.TenantV1alpha1().Partitions().List(options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			options.FieldSelector = fields.Everything().String()
			options.Watch = true
			return kc.TenantV1alpha1().Partitions().Watch(options)
		},
	}, &tntv1al.Partition{}
}

func CacheGetPartition(key string, indexer cache.Indexer, kc kubernetes.Interface) (*tntv1al.Partition, error) {
	if indexer != nil {
		if obj, exist, e := indexer.GetByKey(key); exist && obj != nil && e == nil {
			if partition, _ := obj.(*tntv1al.Partition); partition != nil && partition.Name == key {
				return partition.DeepCopy(), nil
			}
		}
	}
	if kc == nil {
		return nil, errors.ErrVarKubeClientNil
	}
	partition, e := kc.TenantV1alpha1().Partitions().Get(key, metav1.GetOptions{})
	if e != nil {
		return nil, e
	}
	return partition, nil
}

func CacheListPartitions(indexer cache.Indexer, kc kubernetes.Interface) ([]tntv1al.Partition, error) {
	if items := indexer.List(); len(items) > 0 {
		re := make([]tntv1al.Partition, 0, len(items))
		for _, obj := range items {
			partition, _ := obj.(*tntv1al.Partition)
			if partition != nil {
				re = append(re, *partition)
			}
		}
		if len(re) > 0 {
			return re, nil
		}
	}
	partitionList, e := kc.TenantV1alpha1().Partitions().List(metav1.ListOptions{})
	if e != nil {
		return nil, e
	}
	return partitionList.Items, nil
}

func CacheListPartitionsPointer(indexer cache.Indexer, kc kubernetes.Interface) (re []*tntv1al.Partition) {
	// from cache
	items := indexer.List()
	if len(items) > 0 {
		re = make([]*tntv1al.Partition, 0, len(items))
		for _, obj := range items {
			partition, _ := obj.(*tntv1al.Partition)
			if partition != nil {
				re = append(re, partition)
			}
		}
	}
	if len(re) > 0 {
		return re
	}
	// from source
	partitionList, e := kc.TenantV1alpha1().Partitions().List(metav1.ListOptions{})
	if e != nil || len(partitionList.Items) == 0 {
		return nil
	}
	re = make([]*tntv1al.Partition, len(partitionList.Items))
	for i := range partitionList.Items {
		re[i] = &partitionList.Items[i]
	}
	return re
}
