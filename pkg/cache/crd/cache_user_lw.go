package crd

import (
	"fmt"
	"log"
	"sync"

	resv1b1 "github.com/caicloud/clientset/pkg/apis/resource/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"

	"github.com/caicloud/dashboard-admin/pkg/errors"
	"github.com/caicloud/dashboard-admin/pkg/kubernetes"
)

// control cluster

type Config struct {
	Name        string
	Initializer func(kc kubernetes.Interface) (ListWatcher cache.ListerWatcher, ObjType runtime.Object)
}

type ClusterResourcesCache struct {
	// cluster cache
	cc *ListWatchCache
	ec *ClustersCache // export cluster cache

	// cluster:caches
	m     map[string]*subClusterCaches
	mLock sync.RWMutex

	kc      kubernetes.Interface
	kcCache *sync.Map // cluster:kc

	configs []Config
}

func NewDefaultClusterResourcesCache(kc kubernetes.Interface) (rc *ClusterResourcesCache, e error) {
	return NewClusterResourcesCache(kc, GetDefaultConfig())
}

func NewClusterResourcesCache(kc kubernetes.Interface, configs []Config) (rc *ClusterResourcesCache, e error) {
	if e = checkCacheCreateConfigs(kc, configs); e != nil {
		return nil, e
	}
	rc = &ClusterResourcesCache{
		m:       make(map[string]*subClusterCaches),
		kc:      kc,
		kcCache: new(sync.Map),
		configs: configs,
	}
	listWatcher, objType := GetClusterCacheConfig(kc)
	rc.cc, e = NewListWatchCacheWithEventHandler(listWatcher, objType,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    rc.handleClusterAdd,
			UpdateFunc: rc.handleClusterUpdate,
			DeleteFunc: rc.handleClusterDelete,
		})
	if e != nil {
		return nil, e
	}
	rc.ec = &ClustersCache{lwCache: rc.cc, kc: rc.kc, kcCache: rc.kcCache}
	return rc, nil
}

func (rc *ClusterResourcesCache) handleClusterAdd(obj interface{}) {
	cluster := obj.(*resv1b1.Cluster)
	ForceUpdateKubeClientCache(rc.kcCache, cluster)
	rc.updateClusterCache(cluster)
}

func (rc *ClusterResourcesCache) handleClusterUpdate(oldObj, newObj interface{}) {
	cluster := newObj.(*resv1b1.Cluster)
	ForceUpdateKubeClientCache(rc.kcCache, cluster)
	rc.updateClusterCache(cluster)
}

func (rc *ClusterResourcesCache) handleClusterDelete(obj interface{}) {
	cluster := obj.(*resv1b1.Cluster)
	if cluster != nil {
		DeleteInKubeClientCache(rc.kcCache, cluster.Name)
	}
	rc.deleteClusterCache(cluster)
}

func (rc *ClusterResourcesCache) updateClusterCache(cluster *resv1b1.Cluster) {
	// check
	if cluster == nil {
		return
	}
	switch cluster.Status.Phase {
	case ClusterStatusNew:
		return
	case ClusterStatusInstallMaster:
		return
	case ClusterStatusInstallAddon:
	case ClusterStatusReady:
	case ClusterStatusFailed:
		return
	case ClusterStatusDeleting:
		rc.deleteClusterCache(cluster)
		return
	default:
		return
	}
	// try client
	kc, e := rc.ec.GetKubeClient(cluster.Name)
	if e != nil {
		return
	}
	// check exist
	rc.mLock.RLock()
	c := rc.m[cluster.Name]
	rc.mLock.RUnlock()
	if c != nil {
		return
	}
	// try set
	rc.mLock.Lock()
	defer rc.mLock.Unlock()
	c = rc.m[cluster.Name]
	if c != nil {
		return
	}
	c, e = NewSubClusterCaches(kc, rc.configs, cluster.Name)
	if e != nil {
		return
	}
	rc.m[cluster.Name] = c
	go c.Start()
}
func (rc *ClusterResourcesCache) deleteClusterCache(cluster *resv1b1.Cluster) {
	var (
		cleanCaches []*subClusterCaches
	)
	// lock area
	rc.mLock.Lock()
	if cluster != nil {
		// specific cluster
		cleanCaches = append(cleanCaches, rc.m[cluster.Name])
		delete(rc.m, cluster.Name)
	} else {
		// check all cluster
		var cleanList []string
		for k, v := range rc.m {
			item, _, _ := rc.cc.indexer.GetByKey(k)
			if item != nil && item.(*resv1b1.Cluster) != nil {
				cleanList = append(cleanList, k)
				cleanCaches = append(cleanCaches, v)
			}
		}
		for _, cName := range cleanList {
			delete(rc.m, cName)
		}
	}
	rc.mLock.Unlock()
	// stop
	for _, c := range cleanCaches {
		if c != nil {
			c.Stop()
		}
	}
}

func (rc *ClusterResourcesCache) Run(stopCh chan struct{}) {
	rc.cc.Run(stopCh)
	// cleanup
	rc.mLock.Lock()
	for _, c := range rc.m {
		c.Stop()
	}
	rc.mLock.Unlock()
}

func (rc *ClusterResourcesCache) GetAsClusterCache() *ClustersCache {
	return rc.ec
}

func (rc *ClusterResourcesCache) GetSubClusterCaches(clusterName string) (*subClusterCaches, *errors.FormatError) {
	rc.mLock.RLock()
	c := rc.m[clusterName]
	rc.mLock.RUnlock()
	if c != nil {
		return c, nil
	}
	item, _, _ := rc.cc.indexer.GetByKey(clusterName)
	if item != nil && item.(*resv1b1.Cluster) != nil {
		cluster := item.(*resv1b1.Cluster)
		return nil, errors.NewError().SetErrorClusterNotReady(clusterName, string(cluster.Status.Phase))
	}
	return nil, errors.NewError().SetErrorObjectNotFound(clusterName, nil)
}

// sub cluster

type subClusterCaches struct {
	name   string
	kc     kubernetes.Interface
	m      map[string]*ListWatchCache
	stopCh chan struct{}

	hasSynced []func() bool
}

func NewSubClusterCaches(kc kubernetes.Interface, configs []Config, clusterName string) (*subClusterCaches, error) {
	e := checkCacheCreateConfigs(kc, configs)
	if e != nil {
		return nil, e
	}
	scc := &subClusterCaches{
		name:   clusterName,
		kc:     kc,
		m:      make(map[string]*ListWatchCache, len(configs)),
		stopCh: make(chan struct{}),
	}
	for i := range configs {
		name := configs[i].Name
		listWatcher, objType := configs[i].Initializer(kc)
		c, e := NewListWatchCache(listWatcher, objType)
		if e != nil {
			return nil, e
		}
		scc.m[name] = c
		scc.hasSynced = append(scc.hasSynced, c.HasSynced)
	}
	return scc, nil
}

func (scc *subClusterCaches) Start() {
	for k, c := range scc.m {
		if c != nil {
			go func(name string) {
				logPrefix := fmt.Sprintf("[cluster=%s][cache=%s]", scc.name, name)
				log.Printf("%s start", logPrefix)
				c.Run(scc.stopCh)
				log.Printf("%s stopped", logPrefix)
			}(k)
		}
	}
}
func (scc *subClusterCaches) HasSynced() bool {
	for _, f := range scc.hasSynced {
		if !f() {
			return false
		}
	}
	return true
}

func (scc *subClusterCaches) Stop() {
	close(scc.stopCh)
}

func (scc *subClusterCaches) GetCoreCache(name string) (*ListWatchCache, bool) {
	c, ok := scc.m[name]
	return c, ok
}

// config

func checkCacheCreateConfigs(kc kubernetes.Interface, configs []Config) error {
	if kc == nil {
		return errors.ErrVarKubeClientNil
	}
	m := make(map[string]struct{}, len(configs))
	for i := range configs {
		if configs[i].Initializer == nil {
			return errors.ErrVarBadConfig
		}
		if _, ok := m[configs[i].Name]; ok {
			return errors.ErrVarDuplicatedConfig
		}
	}
	return nil
}
