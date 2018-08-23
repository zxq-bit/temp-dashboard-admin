package crd

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	"github.com/caicloud/dashboard-admin/pkg/errors"
	"github.com/caicloud/dashboard-admin/pkg/kubernetes"
)

const (
	CacheNamePod = "Pod"
)

func (scc *subClusterCaches) GetPodCache() (*PodsCache, bool) {
	return scc.GetAsPodCache(CacheNamePod)
}
func (scc *subClusterCaches) GetAsPodCache(name string) (*PodsCache, bool) {
	c, ok := scc.m[name]
	if ok {
		return &PodsCache{lwCache: c, kc: scc.kc}, true
	}
	return nil, false
}

type PodsCache struct {
	lwCache *ListWatchCache
	kc      kubernetes.Interface
}

func NewPodsCache(kc kubernetes.Interface) (*PodsCache, error) {
	listWatcher, objType := GetPodCacheConfig(kc)
	c, e := NewListWatchCache(listWatcher, objType)
	if e != nil {
		return nil, e
	}
	return &PodsCache{
		lwCache: c,
		kc:      kc,
	}, nil
}

func (tc *PodsCache) Run(stopCh chan struct{}) {
	tc.lwCache.Run(stopCh)
}

func (tc *PodsCache) Get(namespace, key string) (*corev1.Pod, error) {
	return CacheGetPod(namespace, key, tc.lwCache.indexer, tc.kc)
}
func (tc *PodsCache) List(namespace string) ([]corev1.Pod, error) {
	return CacheListPods(namespace, tc.lwCache.indexer, tc.kc)
}
func (tc *PodsCache) ListCachePointer(namespace string) (re []*corev1.Pod) {
	return CacheListPodsPointer(namespace, tc.lwCache.indexer, tc.kc)
}

func (tc *PodsCache) Indexes() cache.Indexer {
	return tc.lwCache.indexer
}

func GetPodCacheConfig(kc kubernetes.Interface) (cache.ListerWatcher, runtime.Object) {
	return &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			options.FieldSelector = fields.Everything().String()
			return kc.CoreV1().Pods(metav1.NamespaceAll).List(options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			options.FieldSelector = fields.Everything().String()
			options.Watch = true
			return kc.CoreV1().Pods(metav1.NamespaceAll).Watch(options)
		},
	}, &corev1.Pod{}
}

func CacheGetPod(namespace, key string, indexer cache.Indexer, kc kubernetes.Interface) (*corev1.Pod, error) {
	if indexer != nil {
		if obj, exist, e := indexer.GetByKey(key); exist && obj != nil && e == nil {
			if pod, _ := obj.(*corev1.Pod); CheckNamespace(pod, namespace) && pod.Name == key {
				return pod, nil
			}
		}
	}
	if kc == nil {
		return nil, errors.ErrVarKubeClientNil
	}
	pod, e := kc.CoreV1().Pods(namespace).Get(key, metav1.GetOptions{})
	if e != nil {
		return nil, e
	}
	return pod, nil
}

func CacheListPods(namespace string, indexer cache.Indexer, kc kubernetes.Interface) ([]corev1.Pod, error) {
	if items := indexer.List(); len(items) > 0 {
		re := make([]corev1.Pod, 0, len(items))
		for _, obj := range items {
			pod, _ := obj.(*corev1.Pod)
			if CheckNamespace(pod, namespace) {
				re = append(re, *pod)
			}
		}
		if len(re) > 0 {
			return re, nil
		}
	}
	podList, e := kc.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if e != nil {
		return nil, e
	}
	return podList.Items, nil
}

func CacheListPodsPointer(namespace string, indexer cache.Indexer, kc kubernetes.Interface) (re []*corev1.Pod) {
	// from cache
	items := indexer.List()
	if len(items) > 0 {
		re = make([]*corev1.Pod, 0, len(items))
		for _, obj := range items {
			pod, _ := obj.(*corev1.Pod)
			if CheckNamespace(pod, namespace) {
				re = append(re, pod)
			}
		}
	}
	if len(re) > 0 {
		return re
	}
	// from source
	podList, e := kc.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if e != nil || len(podList.Items) == 0 {
		return nil
	}
	re = make([]*corev1.Pod, len(podList.Items))
	for i := range podList.Items {
		re[i] = &podList.Items[i]
	}
	return re
}
