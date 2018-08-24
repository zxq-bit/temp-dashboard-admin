package api

import (
	"net/http"
	"path"
)

const (
	devopUrlBase    = "api"
	devopApiVersion = "v1"

	workspacesListPath = "workspaces"
	pipelinesListPath  = "pipelines"

	workspacesListCode = 200
	pipelinesListCode  = 200
)

func ListWorkspaces(c *http.Client, devopHost string) (*WorkspaceList, error) {
	re := new(WorkspaceList)
	url := "http://" + path.Join(devopHost, devopUrlBase, devopApiVersion, workspacesListPath)
	e := doGet(c, url, workspacesListCode, re)
	if e != nil {
		return nil, e
	}
	return re, nil
}

func ListPipelines(c *http.Client, devopHost, workspace string) (*PipelineList, error) {
	re := new(PipelineList)
	url := "http://" + path.Join(devopHost, devopUrlBase, devopApiVersion, workspacesListPath, workspace, pipelinesListPath)
	e := doGet(c, url, pipelinesListCode, re)
	if e != nil {
		return nil, e
	}
	return re, nil
}
