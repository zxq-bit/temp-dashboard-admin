package api

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/caicloud/nirvana/log"
)

const (
	CacheNameDevopAdmin = "devop admin"
)

type DaCache struct {
	lock  sync.RWMutex
	wsMap map[string]*WorkspaceDetail
}

type WorkspaceDetail struct {
	Workspace *Workspace
	Pipelines []Pipeline
}

func NewDaCache() (*DaCache, error) {
	c := &DaCache{
		wsMap: make(map[string]*WorkspaceDetail),
	}
	return c, nil
}

func (c *DaCache) Name() string {
	return CacheNameDevopAdmin
}

func (c *DaCache) Refresh(client *http.Client, host string) error {
	workspaces, e := ListWorkspaces(client, host)
	if e != nil {
		log.Errorf("refresh list workspace failed, %v", e)
		return e
	}
	wsMap := make(map[string]*WorkspaceDetail)
	wds := make([]WorkspaceDetail, len(workspaces.Items))
	wg := sync.WaitGroup{}
	ec := make(chan error, len(workspaces.Items))
	for i := range workspaces.Items {
		wds[i].Workspace = &workspaces.Items[i]
		wsMap[wds[i].Workspace.Name] = &wds[i]
		wg.Add(1)
		go func(wd *WorkspaceDetail) {
			pipelineList, e := ListPipelines(client, host, wd.Workspace.Name)
			if e != nil {
				ec <- e
				log.Errorf("refresh list pipeline in workspace %v failed, %v", wd.Workspace.Name, e)
			} else {
				wd.Pipelines = pipelineList.Items
			}
			wg.Done()
		}(&wds[i])
	}
	wg.Wait()

	if len(ec) > 0 {
		errs := readAllErrorsFromChan(ec)
		return fmt.Errorf("failed %d/%d, %v", len(errs), len(wds), errs)
	}

	c.lock.Lock()
	c.wsMap = wsMap
	c.lock.Unlock()
	return nil
}

func (c *DaCache) GetWorkspaceMap() map[string]*WorkspaceDetail {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.wsMap
}
