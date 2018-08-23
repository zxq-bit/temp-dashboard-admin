package crd

import (
	"sync"

	resv1b1 "github.com/caicloud/clientset/pkg/apis/resource/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	"github.com/caicloud/dashboard-admin/pkg/errors"
	"github.com/caicloud/dashboard-admin/pkg/kubernetes"
)

const (
	CacheNameCluster = "Cluster"
)

type ClustersCache struct {
	lwCache *ListWatchCache
	kc      kubernetes.Interface
	kcCache *sync.Map
}

func NewClustersCache(kc kubernetes.Interface) (*ClustersCache, error) {
	cc := &ClustersCache{
		kc:      kc,
		kcCache: new(sync.Map),
	}
	listWatcher, objType := GetClusterCacheConfig(kc)
	c, e := NewListWatchCacheWithEventHandler(listWatcher, objType, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			cluster := obj.(*resv1b1.Cluster)
			ForceUpdateKubeClientCache(cc.kcCache, cluster)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			cluster := newObj.(*resv1b1.Cluster)
			ForceUpdateKubeClientCache(cc.kcCache, cluster)
		},
		DeleteFunc: func(obj interface{}) {
			cluster := obj.(*resv1b1.Cluster)
			if cluster != nil {
				DeleteInKubeClientCache(cc.kcCache, cluster.Name)
			}
		},
	})
	if e != nil {
		return nil, e
	}
	cc.lwCache = c
	return cc, nil
}

func (cc *ClustersCache) Run(stopCh chan struct{}) {
	cc.lwCache.Run(stopCh)
}

func (cc *ClustersCache) Get(key string) (*resv1b1.Cluster, error) {
	return CacheGetCluster(key, cc.lwCache.indexer, cc.kc)
}
func (cc *ClustersCache) List() ([]resv1b1.Cluster, error) {
	return CacheListClusters(cc.lwCache.indexer, cc.kc)
}
func (cc *ClustersCache) ListCachePointer() []*resv1b1.Cluster {
	return CacheListClustersPointer(cc.lwCache.indexer, cc.kc)
}

func (cc *ClustersCache) GetKubeClient(key string) (kubernetes.Interface, error) {
	return CacheGetKubeClient(key, cc.kcCache, cc.Get)
}

func GetClusterCacheConfig(kc kubernetes.Interface) (cache.ListerWatcher, runtime.Object) {
	return &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			options.FieldSelector = fields.Everything().String()
			return kc.ResourceV1beta1().Clusters().List(options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			options.FieldSelector = fields.Everything().String()
			options.Watch = true
			return kc.ResourceV1beta1().Clusters().Watch(options)
		},
	}, &resv1b1.Cluster{}
}

func CacheGetCluster(key string, indexer cache.Indexer, kc kubernetes.Interface) (*resv1b1.Cluster, error) {
	if indexer != nil {
		if obj, exist, e := indexer.GetByKey(key); exist && obj != nil && e == nil {
			if cluster, _ := obj.(*resv1b1.Cluster); cluster != nil && cluster.Name == key {
				return cluster, nil
			}
		}
	}
	if kc == nil {
		return nil, errors.ErrVarKubeClientNil
	}
	cluster, e := kc.ResourceV1beta1().Clusters().Get(key, metav1.GetOptions{})
	if e != nil {
		return nil, e
	}
	return cluster, nil
}

func CacheGetKubeClient(clusterName string, syncMap *sync.Map, clusterGetter func(clusterName string) (*resv1b1.Cluster, error)) (kc kubernetes.Interface, e error) {
	if syncMap != nil {
		obj, ok := syncMap.Load(clusterName)
		if ok && obj != nil {
			kc, ok = obj.(kubernetes.Interface)
			if ok && kc != nil {
				return kc, nil
			}
		}
	}

	cluster, e := clusterGetter(clusterName)
	if e != nil {
		return nil, e
	}
	// if cluster.Status.Phase != kubernetes.ClusterStatusReady { // TODO is OK to remove this check?
	// 	return nil, errors.ErrVarClusterNotReady
	// }
	kc, e = kubernetes.NewClientFromRestConfig(GetKubeConfigFromClusterAuth(&cluster.Spec.Auth))
	if e != nil {
		return nil, e
	}
	if syncMap != nil {
		syncMap.Store(clusterName, kc)
	}
	return kc, nil
}

func CacheListClusters(indexer cache.Indexer, kc kubernetes.Interface) ([]resv1b1.Cluster, error) {
	if items := indexer.List(); len(items) > 0 {
		re := make([]resv1b1.Cluster, 0, len(items))
		for _, obj := range items {
			cluster, _ := obj.(*resv1b1.Cluster)
			if cluster != nil {
				re = append(re, *cluster)
			}
		}
		if len(re) > 0 {
			return re, nil
		}
	}
	clusterList, e := kc.ResourceV1beta1().Clusters().List(metav1.ListOptions{})
	if e != nil {
		return nil, e
	}
	return clusterList.Items, nil
}

func CacheListClustersPointer(indexer cache.Indexer, kc kubernetes.Interface) (re []*resv1b1.Cluster) {
	// from cache
	items := indexer.List()
	if len(items) > 0 {
		re = make([]*resv1b1.Cluster, 0, len(items))
		for _, obj := range items {
			cluster, _ := obj.(*resv1b1.Cluster)
			if cluster != nil {
				re = append(re, cluster)
			}
		}
	}
	if len(re) > 0 {
		return re
	}
	// from source
	clusterList, e := kc.ResourceV1beta1().Clusters().List(metav1.ListOptions{})
	if e != nil || len(clusterList.Items) == 0 {
		return nil
	}
	re = make([]*resv1b1.Cluster, len(clusterList.Items))
	for i := range clusterList.Items {
		re[i] = &clusterList.Items[i]
	}
	return re
}
