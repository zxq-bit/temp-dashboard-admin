package crd

import (
	"fmt"
	"sync"

	resv1b1 "github.com/caicloud/clientset/pkg/apis/resource/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"

	"github.com/caicloud/dashboard-admin/pkg/kubernetes"
)

// cache core

type ListWatchCache struct {
	indexer  cache.Indexer
	informer cache.Controller
}

func NewListWatchCache(listWatcher cache.ListerWatcher, objType runtime.Object) (*ListWatchCache, error) {
	return NewListWatchCacheWithEventHandler(listWatcher, objType, cache.ResourceEventHandlerFuncs{})
}

func NewListWatchCacheWithEventHandler(listWatcher cache.ListerWatcher, objType runtime.Object,
	evHandler cache.ResourceEventHandler) (*ListWatchCache, error) {
	if listWatcher == nil {
		return nil, fmt.Errorf("nil ListerWatcher for ListWatchCache")
	}
	if objType == nil {
		return nil, fmt.Errorf("nil runtime.Object for type")
	}
	indexer, informer := cache.NewIndexerInformer(listWatcher, objType, 0,
		evHandler, cache.Indexers{})
	return &ListWatchCache{
		indexer:  indexer,
		informer: informer,
	}, nil
}

func (c *ListWatchCache) Run(stopCh chan struct{}) {
	defer utilruntime.HandleCrash()

	go c.informer.Run(stopCh)

	// Wait for all involved caches to be synced, before processing items from the queue is started
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return
	}

	<-stopCh
}

func (c *ListWatchCache) Get(key string) (item interface{}, exists bool, err error) {
	return c.indexer.GetByKey(key)
}

func (c *ListWatchCache) GetInNamespace(namespace, key string) (item interface{}, exists bool, err error) {
	return c.indexer.GetByKey(namespace + "/" + key)
}

func (c *ListWatchCache) List() (items []interface{}) {
	return c.indexer.List()
}

func (c *ListWatchCache) HasSynced() bool {
	return c.informer.HasSynced()
}

// kube client

func ForceUpdateKubeClientCache(syncMap *sync.Map, cluster *resv1b1.Cluster) {
	if syncMap == nil || cluster == nil {
		return
	}
	kc, e := kubernetes.NewClientFromRestConfig(GetKubeConfigFromClusterAuth(&cluster.Spec.Auth))
	if e != nil {
		return
	}
	syncMap.Store(cluster.Name, kc)
}

func GetKubeConfigFromClusterAuth(clusterAuth *resv1b1.ClusterAuth) *rest.Config {
	return &rest.Config{
		Username: clusterAuth.KubeUser,
		Password: clusterAuth.KubePassword,
		Host:     fmt.Sprintf("https://%s:%s", clusterAuth.EndpointIP, clusterAuth.EndpointPort),
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}
}

func DeleteInKubeClientCache(syncMap *sync.Map, key string) {
	if syncMap == nil || len(key) == 0 {
		return
	}
	syncMap.Delete(key)
}

// common

func CheckNamespace(obj metav1.Object, namespace string) bool {
	if obj == nil {
		return false
	}
	objNs := obj.GetNamespace()
	if objNs == namespace ||
		len(objNs) == 0 && namespace == metav1.NamespaceDefault ||
		len(namespace) == 0 && objNs == metav1.NamespaceDefault {
		return true
	}
	return false
}
