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
	CacheNameNode = "Node"
)

func (scc *subClusterCaches) GetNodeCache() (*NodesCache, bool) {
	return scc.GetAsNodeCache(CacheNameNode)
}
func (scc *subClusterCaches) GetAsNodeCache(name string) (*NodesCache, bool) {
	c, ok := scc.m[name]
	if ok {
		return &NodesCache{lwCache: c, kc: scc.kc}, true
	}
	return nil, false
}

type NodesCache struct {
	lwCache *ListWatchCache
	kc      kubernetes.Interface
}

func NewNodesCache(kc kubernetes.Interface) (*NodesCache, error) {
	listWatcher, objType := GetNodeCacheConfig(kc)
	c, e := NewListWatchCache(listWatcher, objType)
	if e != nil {
		return nil, e
	}
	return &NodesCache{
		lwCache: c,
		kc:      kc,
	}, nil
}

func (tc *NodesCache) Run(stopCh chan struct{}) {
	tc.lwCache.Run(stopCh)
}

func (tc *NodesCache) Get(key string) (*corev1.Node, error) {
	return CacheGetNode(key, tc.lwCache.indexer, tc.kc)
}
func (tc *NodesCache) List() ([]corev1.Node, error) {
	return CacheListNodes(tc.lwCache.indexer, tc.kc)
}
func (tc *NodesCache) ListCachePointer() (re []*corev1.Node) {
	return CacheListNodesPointer(tc.lwCache.indexer, tc.kc)
}

func (tc *NodesCache) Indexes() cache.Indexer {
	return tc.lwCache.indexer
}

func GetNodeCacheConfig(kc kubernetes.Interface) (cache.ListerWatcher, runtime.Object) {
	return &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			options.FieldSelector = fields.Everything().String()
			return kc.CoreV1().Nodes().List(options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			options.FieldSelector = fields.Everything().String()
			options.Watch = true
			return kc.CoreV1().Nodes().Watch(options)
		},
	}, &corev1.Node{}
}

func CacheGetNode(key string, indexer cache.Indexer, kc kubernetes.Interface) (*corev1.Node, error) {
	if indexer != nil {
		if obj, exist, e := indexer.GetByKey(key); exist && obj != nil && e == nil {
			if node, _ := obj.(*corev1.Node); node != nil && node.Name == key {
				return node, nil
			}
		}
	}
	if kc == nil {
		return nil, errors.ErrVarKubeClientNil
	}
	node, e := kc.CoreV1().Nodes().Get(key, metav1.GetOptions{})
	if e != nil {
		return nil, e
	}
	return node, nil
}

func CacheListNodes(indexer cache.Indexer, kc kubernetes.Interface) ([]corev1.Node, error) {
	if items := indexer.List(); len(items) > 0 {
		re := make([]corev1.Node, 0, len(items))
		for _, obj := range items {
			node, _ := obj.(*corev1.Node)
			if node != nil {
				re = append(re, *node)
			}
		}
		if len(re) > 0 {
			return re, nil
		}
	}
	nodeList, e := kc.CoreV1().Nodes().List(metav1.ListOptions{})
	if e != nil {
		return nil, e
	}
	return nodeList.Items, nil
}

func CacheListNodesPointer(indexer cache.Indexer, kc kubernetes.Interface) (re []*corev1.Node) {
	// from cache
	items := indexer.List()
	if len(items) > 0 {
		re = make([]*corev1.Node, 0, len(items))
		for _, obj := range items {
			node, _ := obj.(*corev1.Node)
			if node != nil {
				re = append(re, node)
			}
		}
	}
	if len(re) > 0 {
		return re
	}
	// from source
	nodeList, e := kc.CoreV1().Nodes().List(metav1.ListOptions{})
	if e != nil || len(nodeList.Items) == 0 {
		return nil
	}
	re = make([]*corev1.Node, len(nodeList.Items))
	for i := range nodeList.Items {
		re[i] = &nodeList.Items[i]
	}
	return re
}
