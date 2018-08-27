package cache

import (
	"fmt"

	"github.com/caicloud/dashboard-admin/pkg/cache/api"
	"github.com/caicloud/dashboard-admin/pkg/cache/crd"
	"github.com/caicloud/dashboard-admin/pkg/config"
	"github.com/caicloud/dashboard-admin/pkg/errors"
	"github.com/caicloud/dashboard-admin/pkg/kubernetes"
)

type Cache struct {
	*crd.ClusterResourcesCache
	*api.Cache
}

func NewCache(cfg *config.Config) (*Cache, error) {
	if cfg == nil {
		return nil, errors.ErrVarKubeClientNil
	}
	e := cfg.Validate()
	if e != nil {
		return nil, e
	}
	kc, e := kubernetes.NewClientFromFlags(cfg.KubeHost, cfg.KubeConfig)
	if e != nil {
		return nil, fmt.Errorf("NewClientFromFlags failed, %v", e)
	}

	cc, e := crd.NewClusterResourcesCache(kc, crd.GetDefaultConfig())
	if e != nil {
		return nil, e
	}
	ac, e := api.NewCache(cfg)
	if e != nil {
		return nil, e
	}

	return &Cache{
		ClusterResourcesCache: cc,
		Cache: ac,
	}, nil
}

func (c *Cache) Run(stopCh chan struct{}) {
	go c.ClusterResourcesCache.Run(stopCh)
	go c.Cache.Run(stopCh)
	<-stopCh
}
