package api

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/caicloud/nirvana/log"
)

const (
	CacheNameCauth = "cauth"
)

type CauthCache struct {
	lock    sync.RWMutex
	users   map[string]*User
	teams   map[string]*Team
	tenants map[string]*Tenant
	roles   map[string]*Role
}

func NewCauthCache() (*CauthCache, error) {
	c := &CauthCache{
		users:   make(map[string]*User),
		teams:   make(map[string]*Team),
		tenants: make(map[string]*Tenant),
		roles:   make(map[string]*Role),
	}
	return c, nil
}

func (c *CauthCache) Name() string {
	return CacheNameCauth
}

func (c *CauthCache) Refresh(client *http.Client, host string) error {
	const (
		mapNum = 4
	)
	var (
		users   map[string]*User
		teams   map[string]*Team
		tenants map[string]*Tenant
		roles   map[string]*Role
		wg      sync.WaitGroup
		ec      = make(chan error, mapNum)
	)
	wg.Add(mapNum)
	go func() {
		var e error
		users, e = GetUsersMap(client, host)
		if e != nil {
			log.Errorf("refresh get user map failed, %v", e)
			ec <- e
		}
		wg.Done()
	}()
	go func() {
		var e error
		teams, e = GetTeamsMap(client, host)
		if e != nil {
			log.Errorf("refresh get team map failed, %v", e)
			ec <- e
		}
		wg.Done()
	}()
	go func() {
		var e error
		tenants, e = GetTenantMap(client, host)
		if e != nil {
			log.Errorf("refresh get tenant map failed, %v", e)
			ec <- e
		}
		wg.Done()
	}()
	go func() {
		var e error
		roles, e = GetRolesMap(client, host)
		if e != nil {
			log.Errorf("refresh get role map failed, %v", e)
			ec <- e
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

	if len(ec) > 0 {
		errs := readAllErrorsFromChan(ec)
		return fmt.Errorf("failed %d/%d, %v", len(errs), mapNum, errs)
	}
	return nil
}

func (c *CauthCache) GetUsersMap() map[string]*User {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.users
}
func (c *CauthCache) GetTeamsMap() map[string]*Team {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.teams
}
func (c *CauthCache) GetTenantsMap() map[string]*Tenant {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.tenants
}
func (c *CauthCache) GetRolesMap() map[string]*Role {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.roles
}
