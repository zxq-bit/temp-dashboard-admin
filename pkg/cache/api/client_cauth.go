package api

import (
	"net/http"
	"path"
)

const (
	dexUrlBase    = "api"
	dexApiVersion = "v2"

	usersListPath  = "users"
	tenantListPath = "tenant"
	teamsListPath  = "teams"
	rolesListPath  = "roles"

	usersListCode  = 200
	tenantListCode = 200
	teamsListCode  = 200
	rolesListCode  = 200
)

func ListUsers(c *http.Client, dexHost string) (*UserList, error) {
	re := new(UserList)
	url := "http://" + path.Join(dexHost, dexUrlBase, dexApiVersion, usersListPath)
	e := doGet(c, url, usersListCode, re)
	if e != nil {
		return nil, e
	}
	return re, nil
}

func ListTenant(c *http.Client, dexHost string) (*TenantList, error) {
	re := new(TenantList)
	url := "http://" + path.Join(dexHost, dexUrlBase, dexApiVersion, tenantListPath)
	e := doGet(c, url, tenantListCode, re)
	if e != nil {
		return nil, e
	}
	return re, nil
}

func ListTeams(c *http.Client, dexHost string) (*TeamList, error) {
	re := new(TeamList)
	url := "http://" + path.Join(dexHost, dexUrlBase, dexApiVersion, teamsListPath)
	e := doGet(c, url, teamsListCode, re)
	if e != nil {
		return nil, e
	}
	return re, nil
}

func ListRoles(c *http.Client, dexHost string) (*RoleList, error) {
	re := new(RoleList)
	url := "http://" + path.Join(dexHost, dexUrlBase, dexApiVersion, rolesListPath)
	e := doGet(c, url, rolesListCode, re)
	if e != nil {
		return nil, e
	}
	return re, nil
}

func GetUsersMap(c *http.Client, dexHost string) (map[string]*User, error) {
	usersList, e := ListUsers(c, dexHost)
	if e != nil {
		return nil, e
	}
	m := make(map[string]*User, len(usersList.Items))
	for i := range usersList.Items {
		user := &usersList.Items[i]
		m[user.Username] = user
	}
	return m, nil
}
func GetTenantMap(c *http.Client, dexHost string) (map[string]*Tenant, error) {
	tenantList, e := ListTenant(c, dexHost)
	if e != nil {
		return nil, e
	}
	m := make(map[string]*Tenant, len(tenantList.Items))
	for i := range tenantList.Items {
		tenant := &tenantList.Items[i]
		m[tenant.ID] = tenant
	}
	return m, nil
}
func GetTeamsMap(c *http.Client, dexHost string) (map[string]*Team, error) {
	teamsList, e := ListTeams(c, dexHost)
	if e != nil {
		return nil, e
	}
	m := make(map[string]*Team, len(teamsList.Items))
	for i := range teamsList.Items {
		team := &teamsList.Items[i]
		m[team.ID] = team
	}
	return m, nil
}
func GetRolesMap(c *http.Client, dexHost string) (map[string]*Role, error) {
	rolesList, e := ListRoles(c, dexHost)
	if e != nil {
		return nil, e
	}
	m := make(map[string]*Role, len(rolesList.Items))
	for i := range rolesList.Items {
		role := &rolesList.Items[i]
		m[role.ID] = role
	}
	return m, nil
}
