package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/caicloud/nirvana/log"

	"github.com/caicloud/dashboard-admin/pkg/config"
)

type Refresher interface {
	Name() string
	Refresh(client *http.Client, host string) error
}

type Cache struct {
	cfg config.Config
	clt *http.Client

	CauthCache *CauthCache
	DevopCache *DaCache
	CargoCache *CargoCache
}

func NewCache(cfg *config.Config) (*Cache, error) {
	if cfg == nil {
		return nil, fmt.Errorf("nil cache config")
	}
	if e := cfg.Validate(); e != nil {
		return nil, e
	}
	clt, e := NewHttpClient(time.Duration(cfg.TimeoutSecond) * time.Second)
	if e != nil {
		return nil, e
	}
	cc, e := NewCauthCache()
	if e != nil {
		return nil, e
	}
	dc, e := NewDaCache()
	if e != nil {
		return nil, e
	}
	cac, e := NewCargoCache()
	if e != nil {
		return nil, e
	}
	return &Cache{
		cfg:        *cfg,
		clt:        clt,
		CauthCache: cc,
		DevopCache: dc,
		CargoCache: cac,
	}, nil
}

func (c *Cache) Run(stopCh chan struct{}) {
	refreshTime := time.Duration(c.cfg.RefreshSecond) * time.Second

	go RunRefresher(c.clt, c.cfg.CauthHost, c.CauthCache, stopCh, refreshTime)
	go RunRefresher(c.clt, c.cfg.DevOpAdminHost, c.DevopCache, stopCh, refreshTime)
	go RunRefresher(c.clt, c.cfg.CargoAdminHost, c.CargoCache, stopCh, refreshTime)

	<-stopCh
}

func RunRefresher(client *http.Client, host string, r Refresher, stopCh chan struct{}, refreshTime time.Duration) {
	name := r.Name()
	log.Infof("%s cache start in refresh time: %v", name, refreshTime)
	r.Refresh(client, host)

	tk := time.NewTicker(refreshTime)
	for {
		select {
		case <-stopCh:
			log.Warningf("%s cache stopped", name)
			tk.Stop()
			return
		case <-tk.C:
			start := time.Now()
			log.Infof("%s cache refresh start", name)
			e := r.Refresh(client, host)
			cost := time.Now().Sub(start)
			if e != nil {
				log.Infof("%s cache refresh done in %v", name, cost)
			} else {
				log.Errorf("%s cache refresh failed in %v, %v", name, cost, e)
			}
		}
	}
}

func readAllErrorsFromChan(ec chan error) []error {
	if len(ec) == 0 {
		return nil
	}
	errs := make([]error, 0, len(ec))
	for len(ec) > 0 {
		e := <-ec
		errs = append(errs, e)
	}
	return errs
}
