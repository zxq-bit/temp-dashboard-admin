package api

import (
	"net/http"
	"path"
)

const (
	cargoUrlBase    = "api"
	cargoApiVersion = "v2"

	registriesListPath = "registries"
	registriesListCode = 200
)

func ListRegistries(c *http.Client, caHost string) (*RegistryList, error) {
	re := new(RegistryList)
	url := "http://" + path.Join(caHost, cargoUrlBase, cargoApiVersion, registriesListPath)
	e := doGet(c, url, registriesListCode, re)
	if e != nil {
		return nil, e
	}
	return re, nil
}

func GetRegistriesMap(c *http.Client, dexHost string) (map[string]*Registry, error) {
	registryList, e := ListRegistries(c, dexHost)
	if e != nil {
		return nil, e
	}
	m := make(map[string]*Registry, len(registryList.Items))
	for i := range registryList.Items {
		registry := &registryList.Items[i]
		m[registry.Metadata.Name] = registry
	}
	return m, nil
}
