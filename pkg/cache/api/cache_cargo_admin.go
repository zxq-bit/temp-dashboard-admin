package api

import (
	"net/http"
	"sync"

	"github.com/caicloud/nirvana/log"
)

const (
	CacheNameCargo = "cargo admin"
)

type CargoCache struct {
	lock       sync.RWMutex
	registries map[string]*Registry
}

func NewCargoCache() (*CargoCache, error) {
	c := &CargoCache{
		registries: make(map[string]*Registry),
	}
	return c, nil
}

func (c *CargoCache) Name() string {
	return CacheNameCargo
}

func (c *CargoCache) Refresh(client *http.Client, host string) error {
	registries, e := GetRegistriesMap(client, host)
	if e != nil {
		log.Errorf("refresh list registry failed, %v")
		return e
	}

	c.lock.Lock()
	c.registries = registries
	c.lock.Unlock()

	return nil
}

func (c *CargoCache) GetRegistriesMap() map[string]*Registry {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.registries
}
