package crd

import (
	lbv1a2 "github.com/caicloud/clientset/pkg/apis/loadbalance/v1alpha2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	"github.com/caicloud/dashboard-admin/pkg/errors"
	"github.com/caicloud/dashboard-admin/pkg/kubernetes"
)

const (
	CacheNameLoadBalancer = "LoadBalancer"
)

func (scc *subClusterCaches) GetLoadBalancerCache() (*LoadBalancersCache, bool) {
	return scc.GetAsLoadBalancerCache(CacheNameLoadBalancer)
}
func (scc *subClusterCaches) GetAsLoadBalancerCache(name string) (*LoadBalancersCache, bool) {
	c, ok := scc.m[name]
	if ok {
		return &LoadBalancersCache{lwCache: c, kc: scc.kc}, true
	}
	return nil, false
}

type LoadBalancersCache struct {
	lwCache *ListWatchCache
	kc      kubernetes.Interface
}

func NewLoadBalancersCache(kc kubernetes.Interface) (*LoadBalancersCache, error) {
	listWatcher, objType := GetLoadBalancerCacheConfig(kc)
	c, e := NewListWatchCache(listWatcher, objType)
	if e != nil {
		return nil, e
	}
	return &LoadBalancersCache{
		lwCache: c,
		kc:      kc,
	}, nil
}

func (tc *LoadBalancersCache) Run(stopCh chan struct{}) {
	tc.lwCache.Run(stopCh)
}

func (tc *LoadBalancersCache) Get(namespace, key string) (*lbv1a2.LoadBalancer, error) {
	return CacheGetLoadBalancer(namespace, key, tc.lwCache.indexer, tc.kc)
}
func (tc *LoadBalancersCache) List(namespace string) ([]lbv1a2.LoadBalancer, error) {
	return CacheListLoadBalancers(namespace, tc.lwCache.indexer, tc.kc)
}
func (tc *LoadBalancersCache) ListCachePointer(namespace string) (re []*lbv1a2.LoadBalancer) {
	return CacheListLoadBalancersPointer(namespace, tc.lwCache.indexer, tc.kc)
}

func (tc *LoadBalancersCache) Indexes() cache.Indexer {
	return tc.lwCache.indexer
}

func GetLoadBalancerCacheConfig(kc kubernetes.Interface) (cache.ListerWatcher, runtime.Object) {
	return &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			options.FieldSelector = fields.Everything().String()
			return kc.LoadbalanceV1alpha2().LoadBalancers(metav1.NamespaceAll).List(options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			options.FieldSelector = fields.Everything().String()
			options.Watch = true
			return kc.LoadbalanceV1alpha2().LoadBalancers(metav1.NamespaceAll).Watch(options)
		},
	}, &lbv1a2.LoadBalancer{}
}

func CacheGetLoadBalancer(namespace, key string, indexer cache.Indexer, kc kubernetes.Interface) (*lbv1a2.LoadBalancer, error) {
	if indexer != nil {
		if obj, exist, e := indexer.GetByKey(key); exist && obj != nil && e == nil {
			if loadBalancer, _ := obj.(*lbv1a2.LoadBalancer); CheckNamespace(loadBalancer, namespace) && loadBalancer.Name == key {
				return loadBalancer, nil
			}
		}
	}
	if kc == nil {
		return nil, errors.ErrVarKubeClientNil
	}
	loadBalancer, e := kc.LoadbalanceV1alpha2().LoadBalancers(namespace).Get(key, metav1.GetOptions{})
	if e != nil {
		return nil, e
	}
	return loadBalancer, nil
}

func CacheListLoadBalancers(namespace string, indexer cache.Indexer, kc kubernetes.Interface) ([]lbv1a2.LoadBalancer, error) {
	if items := indexer.List(); len(items) > 0 {
		re := make([]lbv1a2.LoadBalancer, 0, len(items))
		for _, obj := range items {
			loadBalancer, _ := obj.(*lbv1a2.LoadBalancer)
			if CheckNamespace(loadBalancer, namespace) {
				re = append(re, *loadBalancer)
			}
		}
		if len(re) > 0 {
			return re, nil
		}
	}
	loadBalancerList, e := kc.LoadbalanceV1alpha2().LoadBalancers(namespace).List(metav1.ListOptions{})
	if e != nil {
		return nil, e
	}
	return loadBalancerList.Items, nil
}

func CacheListLoadBalancersPointer(namespace string, indexer cache.Indexer, kc kubernetes.Interface) (re []*lbv1a2.LoadBalancer) {
	// from cache
	items := indexer.List()
	if len(items) > 0 {
		re = make([]*lbv1a2.LoadBalancer, 0, len(items))
		for _, obj := range items {
			loadBalancer, _ := obj.(*lbv1a2.LoadBalancer)
			if CheckNamespace(loadBalancer, namespace) {
				re = append(re, loadBalancer)
			}
		}
	}
	if len(re) > 0 {
		return re
	}
	// from source
	loadBalancerList, e := kc.LoadbalanceV1alpha2().LoadBalancers(namespace).List(metav1.ListOptions{})
	if e != nil || len(loadBalancerList.Items) == 0 {
		return nil
	}
	re = make([]*lbv1a2.LoadBalancer, len(loadBalancerList.Items))
	for i := range loadBalancerList.Items {
		re[i] = &loadBalancerList.Items[i]
	}
	return re
}
