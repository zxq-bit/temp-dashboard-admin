package db

import (
	"fmt"
	"sync"
	"time"

	"github.com/caicloud/nirvana/log"
)

type DexCache struct {
	dexHost     string
	timeout     time.Duration
	refreshTime time.Duration

	httpClt *Client
	lock    sync.RWMutex
	users   map[string]*User
	teams   map[string]*Team
	tenants map[string]*Tenant
	roles   map[string]*Role
}

func NewDexCache(dexHost string, timeout, refreshTime time.Duration) (*DexCache, error) {
	if len(dexHost) == 0 {
		return nil, fmt.Errorf("nil dex host")
	}
	if timeout < 0 {
		return nil, fmt.Errorf("illegal timeout: %v", timeout)
	}
	if refreshTime < 1 {
		return nil, fmt.Errorf("illegal refresh seconds: %v", refreshTime)
	}
	hc, e := NewClient(timeout)
	if e != nil {
		return nil, fmt.Errorf("NewClient failed, %v", e)
	}
	c := &DexCache{
		dexHost:     dexHost,
		timeout:     timeout,
		refreshTime: refreshTime,
		httpClt:     hc,
		users:       make(map[string]*User),
		teams:       make(map[string]*Team),
		tenants:     make(map[string]*Tenant),
		roles:       make(map[string]*Role),
	}
	return c, nil
}

func (c *DexCache) Run(stopCh chan struct{}) {
	log.Infof("dex cache start in refresh time: %v", c.refreshTime)
	c.refresh()

	tk := time.NewTicker(c.refreshTime)
	for {
		select {
		case <-stopCh:
			log.Warningf("dex cache stopped")
			tk.Stop()
			return
		case <-tk.C:
			log.Infof("dex cache refresh")
			c.refresh()
		}
	}
}

func (c *DexCache) refresh() {
	var (
		users   map[string]*User
		teams   map[string]*Team
		tenants map[string]*Tenant
		roles   map[string]*Role
		wg      sync.WaitGroup
		e       error
	)
	wg.Add(4)
	go func() {
		users, e = c.httpClt.GetUsersMap(c.dexHost)
		if e != nil {
			log.Errorf("refresh get user map failed, %v")
		}
		wg.Done()
	}()
	go func() {
		teams, e = c.httpClt.GetTeamsMap(c.dexHost)
		if e != nil {
			log.Errorf("refresh get team map failed, %v")
		}
		wg.Done()
	}()
	go func() {
		tenants, e = c.httpClt.GetTenantMap(c.dexHost)
		if e != nil {
			log.Errorf("refresh get tenant map failed, %v")
		}
		wg.Done()
	}()
	go func() {
		roles, e = c.httpClt.GetRolesMap(c.dexHost)
		if e != nil {
			log.Errorf("refresh get role map failed, %v")
		}
		wg.Done()
	}()
	wg.Wait()

	c.lock.Lock()
	defer c.lock.Unlock()

	if users != nil {
		c.users = users
	}
	if teams != nil {
		c.teams = teams
	}
	if tenants != nil {
		c.tenants = tenants
	}
	if roles != nil {
		c.roles = roles
	}
}

func (c *DexCache) GetUsersMap() map[string]*User {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.users
}
func (c *DexCache) GetTeamsMap() map[string]*Team {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.teams
}
func (c *DexCache) GetTenantsMap() map[string]*Tenant {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.tenants
}
func (c *DexCache) GetRolesMap() map[string]*Role {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.roles
}
