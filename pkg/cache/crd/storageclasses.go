package crd

import (
	storagev1 "k8s.io/api/storage/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	"github.com/caicloud/dashboard-admin/pkg/errors"
	"github.com/caicloud/dashboard-admin/pkg/kubernetes"
)

func (scc *subClusterCaches) GetStorageClassCache() (*StorageClassCache, bool) {
	return scc.GetAsStorageClassCache(CacheNameStorageClass)
}
func (scc *subClusterCaches) GetAsStorageClassCache(name string) (*StorageClassCache, bool) {
	c, ok := scc.m[name]
	if ok {
		return &StorageClassCache{lwCache: c, kc: scc.kc}, true
	}
	return nil, false
}

type StorageClassCache struct {
	lwCache *ListWatchCache
	kc      kubernetes.Interface
}

func NewStorageClassCache(kc kubernetes.Interface) (*StorageClassCache, error) {
	listWatcher, objType := GetStorageClassCacheConfig(kc)
	c, e := NewListWatchCache(listWatcher, objType)
	if e != nil {
		return nil, e
	}
	return &StorageClassCache{
		lwCache: c,
		kc:      kc,
	}, nil
}

func (tc *StorageClassCache) Run(stopCh chan struct{}) {
	tc.lwCache.Run(stopCh)
}

func (tc *StorageClassCache) Get(key string) (*storagev1.StorageClass, error) {
	return CacheGetStorageClass(key, tc.lwCache.indexer, tc.kc)
}
func (tc *StorageClassCache) List() ([]storagev1.StorageClass, error) {
	return CacheListStorageClasses(tc.lwCache.indexer, tc.kc)
}
func (tc *StorageClassCache) ListCachePointer() (re []*storagev1.StorageClass) {
	return CacheListStorageClassesPointer(tc.lwCache.indexer, tc.kc)
}

func GetStorageClassCacheConfig(kc kubernetes.Interface) (cache.ListerWatcher, runtime.Object) {
	return &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			options.FieldSelector = fields.Everything().String()
			return kc.StorageV1beta1().StorageClasses().List(options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			options.FieldSelector = fields.Everything().String()
			options.Watch = true
			return kc.StorageV1beta1().StorageClasses().Watch(options)
		},
	}, &storagev1.StorageClass{}
}

func CacheGetStorageClass(key string, indexer cache.Indexer, kc kubernetes.Interface) (*storagev1.StorageClass, error) {
	if indexer != nil {
		if obj, exist, e := indexer.GetByKey(key); exist && obj != nil && e == nil {
			if sc, _ := obj.(*storagev1.StorageClass); sc != nil && sc.Name == key {
				return sc.DeepCopy(), nil
			}
		}
	}
	if kc == nil {
		return nil, errors.ErrVarKubeClientNil
	}
	sc, e := kc.StorageV1beta1().StorageClasses().Get(key, metav1.GetOptions{})
	if e != nil {
		return nil, e
	}
	return sc, nil
}

func CacheListStorageClasses(indexer cache.Indexer, kc kubernetes.Interface) ([]storagev1.StorageClass, error) {
	if items := indexer.List(); len(items) > 0 {
		re := make([]storagev1.StorageClass, 0, len(items))
		for _, obj := range items {
			sc, _ := obj.(*storagev1.StorageClass)
			if sc != nil {
				re = append(re, *sc)
			}
		}
		if len(re) > 0 {
			return re, nil
		}
	}
	scList, e := kc.StorageV1beta1().StorageClasses().List(metav1.ListOptions{})
	if e != nil {
		return nil, e
	}
	return scList.Items, nil
}

func CacheListStorageClassesPointer(indexer cache.Indexer, kc kubernetes.Interface) (re []*storagev1.StorageClass) {
	// from cache
	items := indexer.List()
	if len(items) > 0 {
		re = make([]*storagev1.StorageClass, 0, len(items))
		for _, obj := range items {
			sc, _ := obj.(*storagev1.StorageClass)
			if sc != nil {
				re = append(re, sc)
			}
		}
	}
	if len(re) > 0 {
		return re
	}
	// from source
	scList, e := kc.StorageV1beta1().StorageClasses().List(metav1.ListOptions{})
	if e != nil || len(scList.Items) == 0 {
		return nil
	}
	re = make([]*storagev1.StorageClass, len(scList.Items))
	for i := range scList.Items {
		re[i] = &scList.Items[i]
	}
	return re
}
