package db

import "path"

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

func (c *Client) ListUsers(dexHost string) (*UserList, error) {
	re := new(UserList)
	url := "http://" + path.Join(dexHost, dexUrlBase, dexApiVersion, usersListPath)
	e := c.doGet(url, usersListCode, re)
	if e != nil {
		return nil, e
	}
	return re, nil
}

func (c *Client) ListTenant(dexHost string) (*TenantList, error) {
	re := new(TenantList)
	url := "http://" + path.Join(dexHost, dexUrlBase, dexApiVersion, tenantListPath)
	e := c.doGet(url, tenantListCode, re)
	if e != nil {
		return nil, e
	}
	return re, nil
}

func (c *Client) ListTeams(dexHost string) (*TeamList, error) {
	re := new(TeamList)
	url := "http://" + path.Join(dexHost, dexUrlBase, dexApiVersion, teamsListPath)
	e := c.doGet(url, teamsListCode, re)
	if e != nil {
		return nil, e
	}
	return re, nil
}

func (c *Client) ListRoles(dexHost string) (*RoleList, error) {
	re := new(RoleList)
	url := "http://" + path.Join(dexHost, dexUrlBase, dexApiVersion, rolesListPath)
	e := c.doGet(url, rolesListCode, re)
	if e != nil {
		return nil, e
	}
	return re, nil
}

func (c *Client) GetUsersMap(dexHost string) (map[string]*User, error) {
	usersList, e := c.ListUsers(dexHost)
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
func (c *Client) GetTenantMap(dexHost string) (map[string]*Tenant, error) {
	tenantList, e := c.ListTenant(dexHost)
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
func (c *Client) GetTeamsMap(dexHost string) (map[string]*Team, error) {
	teamsList, e := c.ListTeams(dexHost)
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
func (c *Client) GetRolesMap(dexHost string) (map[string]*Role, error) {
	rolesList, e := c.ListRoles(dexHost)
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
