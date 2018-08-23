package crd

import (
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	"github.com/caicloud/dashboard-admin/pkg/errors"
	"github.com/caicloud/dashboard-admin/pkg/kubernetes"
)

const (
	CacheNameStorageClass = "StorageClass"
)

func (scc *subClusterCaches) GetStorageClassCache() (*StorageClassesCache, bool) {
	return scc.GetAsStorageClassCache(CacheNameStorageClass)
}
func (scc *subClusterCaches) GetAsStorageClassCache(name string) (*StorageClassesCache, bool) {
	c, ok := scc.m[name]
	if ok {
		return &StorageClassesCache{lwCache: c, kc: scc.kc}, true
	}
	return nil, false
}

type StorageClassesCache struct {
	lwCache *ListWatchCache
	kc      kubernetes.Interface
}

func NewStorageClassesCache(kc kubernetes.Interface) (*StorageClassesCache, error) {
	listWatcher, objType := GetStorageClassCacheConfig(kc)
	c, e := NewListWatchCache(listWatcher, objType)
	if e != nil {
		return nil, e
	}
	return &StorageClassesCache{
		lwCache: c,
		kc:      kc,
	}, nil
}

func (tc *StorageClassesCache) Run(stopCh chan struct{}) {
	tc.lwCache.Run(stopCh)
}

func (tc *StorageClassesCache) Get(namespace, key string) (*storagev1.StorageClass, error) {
	return CacheGetStorageClass(namespace, key, tc.lwCache.indexer, tc.kc)
}
func (tc *StorageClassesCache) List(namespace string) ([]storagev1.StorageClass, error) {
	return CacheListStorageClasses(namespace, tc.lwCache.indexer, tc.kc)
}
func (tc *StorageClassesCache) ListCachePointer(namespace string) (re []*storagev1.StorageClass) {
	return CacheListStorageClassesPointer(namespace, tc.lwCache.indexer, tc.kc)
}

func (tc *StorageClassesCache) Indexes() cache.Indexer {
	return tc.lwCache.indexer
}

func GetStorageClassCacheConfig(kc kubernetes.Interface) (cache.ListerWatcher, runtime.Object) {
	return &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			options.FieldSelector = fields.Everything().String()
			return kc.StorageV1().StorageClasses(metav1.NamespaceAll).List(options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			options.FieldSelector = fields.Everything().String()
			options.Watch = true
			return kc.StorageV1().StorageClasses(metav1.NamespaceAll).Watch(options)
		},
	}, &storagev1.StorageClass{}
}

func CacheGetStorageClass(namespace, key string, indexer cache.Indexer, kc kubernetes.Interface) (*storagev1.StorageClass, error) {
	if indexer != nil {
		if obj, exist, e := indexer.GetByKey(key); exist && obj != nil && e == nil {
			if storageClass, _ := obj.(*storagev1.StorageClass); CheckNamespace(storageClass, namespace) && storageClass.Name == key {
				return storageClass, nil
			}
		}
	}
	if kc == nil {
		return nil, errors.ErrVarKubeClientNil
	}
	storageClass, e := kc.StorageV1().StorageClasses(namespace).Get(key, metav1.GetOptions{})
	if e != nil {
		return nil, e
	}
	return storageClass, nil
}

func CacheListStorageClasses(namespace string, indexer cache.Indexer, kc kubernetes.Interface) ([]storagev1.StorageClass, error) {
	if items := indexer.List(); len(items) > 0 {
		re := make([]storagev1.StorageClass, 0, len(items))
		for _, obj := range items {
			storageClass, _ := obj.(*storagev1.StorageClass)
			if CheckNamespace(storageClass, namespace) {
				re = append(re, *storageClass)
			}
		}
		if len(re) > 0 {
			return re, nil
		}
	}
	storageClassList, e := kc.StorageV1().StorageClasses(namespace).List(metav1.ListOptions{})
	if e != nil {
		return nil, e
	}
	return storageClassList.Items, nil
}

func CacheListStorageClassesPointer(namespace string, indexer cache.Indexer, kc kubernetes.Interface) (re []*storagev1.StorageClass) {
	// from cache
	items := indexer.List()
	if len(items) > 0 {
		re = make([]*storagev1.StorageClass, 0, len(items))
		for _, obj := range items {
			storageClass, _ := obj.(*storagev1.StorageClass)
			if CheckNamespace(storageClass, namespace) {
				re = append(re, storageClass)
			}
		}
	}
	if len(re) > 0 {
		return re
	}
	// from source
	storageClassList, e := kc.StorageV1().StorageClasses(namespace).List(metav1.ListOptions{})
	if e != nil || len(storageClassList.Items) == 0 {
		return nil
	}
	re = make([]*storagev1.StorageClass, len(storageClassList.Items))
	for i := range storageClassList.Items {
		re[i] = &storageClassList.Items[i]
	}
	return re
}
