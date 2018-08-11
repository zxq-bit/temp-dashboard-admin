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

func (scc *subClusterCaches) GetClusterQuotaCache() (*ClusterQuotasCache, bool) {
	return scc.GetAsClusterQuotaCache(CacheNameClusterQuota)
}
func (scc *subClusterCaches) GetAsClusterQuotaCache(name string) (*ClusterQuotasCache, bool) {
	c, ok := scc.m[name]
	if ok {
		return &ClusterQuotasCache{lwCache: c, kc: scc.kc}, true
	}
	return nil, false
}

type ClusterQuotasCache struct {
	lwCache *ListWatchCache
	kc      kubernetes.Interface
}

func NewClusterQuotasCache(kc kubernetes.Interface) (*ClusterQuotasCache, error) {
	listWatcher, objType := GetClusterQuotaCacheConfig(kc)
	c, e := NewListWatchCache(listWatcher, objType)
	if e != nil {
		return nil, e
	}
	return &ClusterQuotasCache{
		lwCache: c,
		kc:      kc,
	}, nil
}

func (tc *ClusterQuotasCache) Run(stopCh chan struct{}) {
	tc.lwCache.Run(stopCh)
}

func (tc *ClusterQuotasCache) Get(key string) (*tntv1al.ClusterQuota, error) {
	return CacheGetClusterQuota(key, tc.lwCache.indexer, tc.kc)
}
func (tc *ClusterQuotasCache) List() ([]tntv1al.ClusterQuota, error) {
	return CacheListClusterQuotas(tc.lwCache.indexer, tc.kc)
}
func (tc *ClusterQuotasCache) ListCachePointer() (re []*tntv1al.ClusterQuota) {
	return CacheListClusterQuotasPointer(tc.lwCache.indexer, tc.kc)
}

func (tc *ClusterQuotasCache) Indexes() cache.Indexer {
	return tc.lwCache.indexer
}

func GetClusterQuotaCacheConfig(kc kubernetes.Interface) (cache.ListerWatcher, runtime.Object) {
	return &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			options.FieldSelector = fields.Everything().String()
			return kc.TenantV1alpha1().ClusterQuotas().List(options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			options.FieldSelector = fields.Everything().String()
			options.Watch = true
			return kc.TenantV1alpha1().ClusterQuotas().Watch(options)
		},
	}, &tntv1al.ClusterQuota{}
}

func CacheGetClusterQuota(key string, indexer cache.Indexer, kc kubernetes.Interface) (*tntv1al.ClusterQuota, error) {
	if indexer != nil {
		if obj, exist, e := indexer.GetByKey(key); exist && obj != nil && e == nil {
			if clusterQuota, _ := obj.(*tntv1al.ClusterQuota); clusterQuota != nil && clusterQuota.Name == key {
				return clusterQuota.DeepCopy(), nil
			}
		}
	}
	if kc == nil {
		return nil, errors.ErrVarKubeClientNil
	}
	clusterQuota, e := kc.TenantV1alpha1().ClusterQuotas().Get(key, metav1.GetOptions{})
	if e != nil {
		return nil, e
	}
	return clusterQuota, nil
}

func CacheListClusterQuotas(indexer cache.Indexer, kc kubernetes.Interface) ([]tntv1al.ClusterQuota, error) {
	if items := indexer.List(); len(items) > 0 {
		re := make([]tntv1al.ClusterQuota, 0, len(items))
		for _, obj := range items {
			clusterQuota, _ := obj.(*tntv1al.ClusterQuota)
			if clusterQuota != nil {
				re = append(re, *clusterQuota)
			}
		}
		if len(re) > 0 {
			return re, nil
		}
	}
	clusterQuotaList, e := kc.TenantV1alpha1().ClusterQuotas().List(metav1.ListOptions{})
	if e != nil {
		return nil, e
	}
	return clusterQuotaList.Items, nil
}

func CacheListClusterQuotasPointer(indexer cache.Indexer, kc kubernetes.Interface) (re []*tntv1al.ClusterQuota) {
	// from cache
	items := indexer.List()
	if len(items) > 0 {
		re = make([]*tntv1al.ClusterQuota, 0, len(items))
		for _, obj := range items {
			clusterQuota, _ := obj.(*tntv1al.ClusterQuota)
			if clusterQuota != nil {
				re = append(re, clusterQuota)
			}
		}
	}
	if len(re) > 0 {
		return re
	}
	// from source
	clusterQuotaList, e := kc.TenantV1alpha1().ClusterQuotas().List(metav1.ListOptions{})
	if e != nil || len(clusterQuotaList.Items) == 0 {
		return nil
	}
	re = make([]*tntv1al.ClusterQuota, len(clusterQuotaList.Items))
	for i := range clusterQuotaList.Items {
		re[i] = &clusterQuotaList.Items[i]
	}
	return re
}
