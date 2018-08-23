package crd

import (
	rlsv1a1 "github.com/caicloud/clientset/pkg/apis/release/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	"github.com/caicloud/dashboard-admin/pkg/errors"
	"github.com/caicloud/dashboard-admin/pkg/kubernetes"
)

const (
	CacheNameRelease = "Release"
)

func (scc *subClusterCaches) GetReleaseCache() (*ReleasesCache, bool) {
	return scc.GetAsReleaseCache(CacheNameRelease)
}
func (scc *subClusterCaches) GetAsReleaseCache(name string) (*ReleasesCache, bool) {
	c, ok := scc.m[name]
	if ok {
		return &ReleasesCache{lwCache: c, kc: scc.kc}, true
	}
	return nil, false
}

type ReleasesCache struct {
	lwCache *ListWatchCache
	kc      kubernetes.Interface
}

func NewReleasesCache(kc kubernetes.Interface) (*ReleasesCache, error) {
	listWatcher, objType := GetReleaseCacheConfig(kc)
	c, e := NewListWatchCache(listWatcher, objType)
	if e != nil {
		return nil, e
	}
	return &ReleasesCache{
		lwCache: c,
		kc:      kc,
	}, nil
}

func (tc *ReleasesCache) Run(stopCh chan struct{}) {
	tc.lwCache.Run(stopCh)
}

func (tc *ReleasesCache) Get(namespace, key string) (*rlsv1a1.Release, error) {
	return CacheGetRelease(namespace, key, tc.lwCache.indexer, tc.kc)
}
func (tc *ReleasesCache) List(namespace string) ([]rlsv1a1.Release, error) {
	return CacheListReleases(namespace, tc.lwCache.indexer, tc.kc)
}
func (tc *ReleasesCache) ListCachePointer(namespace string) (re []*rlsv1a1.Release) {
	return CacheListReleasesPointer(namespace, tc.lwCache.indexer, tc.kc)
}

func (tc *ReleasesCache) Indexes() cache.Indexer {
	return tc.lwCache.indexer
}

func GetReleaseCacheConfig(kc kubernetes.Interface) (cache.ListerWatcher, runtime.Object) {
	return &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			options.FieldSelector = fields.Everything().String()
			return kc.ReleaseV1alpha1().Releases(metav1.NamespaceAll).List(options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			options.FieldSelector = fields.Everything().String()
			options.Watch = true
			return kc.ReleaseV1alpha1().Releases(metav1.NamespaceAll).Watch(options)
		},
	}, &rlsv1a1.Release{}
}

func CacheGetRelease(namespace, key string, indexer cache.Indexer, kc kubernetes.Interface) (*rlsv1a1.Release, error) {
	if indexer != nil {
		if obj, exist, e := indexer.GetByKey(key); exist && obj != nil && e == nil {
			if release, _ := obj.(*rlsv1a1.Release); CheckNamespace(release, namespace) && release.Name == key {
				return release, nil
			}
		}
	}
	if kc == nil {
		return nil, errors.ErrVarKubeClientNil
	}
	release, e := kc.ReleaseV1alpha1().Releases(namespace).Get(key, metav1.GetOptions{})
	if e != nil {
		return nil, e
	}
	return release, nil
}

func CacheListReleases(namespace string, indexer cache.Indexer, kc kubernetes.Interface) ([]rlsv1a1.Release, error) {
	if items := indexer.List(); len(items) > 0 {
		re := make([]rlsv1a1.Release, 0, len(items))
		for _, obj := range items {
			release, _ := obj.(*rlsv1a1.Release)
			if CheckNamespace(release, namespace) {
				re = append(re, *release)
			}
		}
		if len(re) > 0 {
			return re, nil
		}
	}
	releaseList, e := kc.ReleaseV1alpha1().Releases(namespace).List(metav1.ListOptions{})
	if e != nil {
		return nil, e
	}
	return releaseList.Items, nil
}

func CacheListReleasesPointer(namespace string, indexer cache.Indexer, kc kubernetes.Interface) (re []*rlsv1a1.Release) {
	// from cache
	items := indexer.List()
	if len(items) > 0 {
		re = make([]*rlsv1a1.Release, 0, len(items))
		for _, obj := range items {
			release, _ := obj.(*rlsv1a1.Release)
			if CheckNamespace(release, namespace) {
				re = append(re, release)
			}
		}
	}
	if len(re) > 0 {
		return re
	}
	// from source
	releaseList, e := kc.ReleaseV1alpha1().Releases(namespace).List(metav1.ListOptions{})
	if e != nil || len(releaseList.Items) == 0 {
		return nil
	}
	re = make([]*rlsv1a1.Release, len(releaseList.Items))
	for i := range releaseList.Items {
		re[i] = &releaseList.Items[i]
	}
	return re
}
