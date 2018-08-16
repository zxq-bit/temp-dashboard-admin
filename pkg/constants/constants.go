package constants

import "fmt"

const (
	APIVersion = "v1alpha1"
)

var (
	RootPath = fmt.Sprintf("/api/%s", APIVersion)
)

const (
	ParameterStart = "start"
	ParameterLimit = "limit"

	ParameterCluster     = "cluster"
	ParameterMachine     = "machine"
	ParameterNode        = "node"
	ParameterCI          = "ci"
	ParameterRequestBody = "req"
	ParameterXUser       = "X-User"
	ParameterXTenant     = "X-Tenant"
)

const (
	DefaultKubeHost   = ""
	DefaultKubeConfig = ""
	DefaultListenPort = 2587
)
