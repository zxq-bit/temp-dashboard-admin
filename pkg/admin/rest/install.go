package rest

import (
	"fmt"
	"path"

	"github.com/caicloud/nirvana/definition"

	"github.com/caicloud/dashboard-admin/pkg/cache"
	"github.com/caicloud/dashboard-admin/pkg/constants"
)

var (
	QueryParamStart = definition.Parameter{
		Name:        constants.ParameterStart,
		Description: "page split start",
		Source:      definition.Query,
	}
	QueryParamLimit = definition.Parameter{
		Name:        constants.ParameterLimit,
		Description: "page split limit",
		Source:      definition.Query,
	}
	PathParamCluster = definition.Parameter{
		Name:        constants.ParameterCluster,
		Description: "cluster id",
		Source:      definition.Path,
	}
	BodyParamRequest = definition.Parameter{
		Name:        constants.ParameterRequestBody,
		Description: "request body",
		Source:      definition.Body,
	}
	HeaderParamXUser = definition.Parameter{
		Name:        constants.ParameterXUser,
		Description: "request operator",
		Source:      definition.Header,
	}
	HeaderParamXTenant = definition.Parameter{
		Name:        constants.ParameterXTenant,
		Description: "request tenant",
		Source:      definition.Header,
	}
)

func InitNirvanaDescriptors(c *cache.Cache) []definition.Descriptor {
	commonResults := definition.DataErrorResults("result")
	return []definition.Descriptor{
		{
			Path: path.Join(constants.RootPath, fmt.Sprintf("/clusters")),
			Definitions: []definition.Definition{
				{
					Description: "get cluster info",
					Method:      definition.List,
					Function:    HandleListClusterInfo(c),
					Consumes:    []string{definition.MIMEAll}, Produces: []string{definition.MIMEJSON},
					Parameters: []definition.Parameter{
						HeaderParamXTenant, HeaderParamXUser,
						QueryParamStart, QueryParamLimit,
					},
					Results: commonResults,
				},
			},
		},
		{
			Path: path.Join(constants.RootPath, fmt.Sprintf("/clusters/{%s}/machines", constants.ParameterCluster)),
			Definitions: []definition.Definition{
				{
					Description: "get cluster machines info",
					Method:      definition.Get,
					Function:    HandleGetMachineSummary(c),
					Consumes:    []string{definition.MIMEAll}, Produces: []string{definition.MIMEJSON},
					Parameters: []definition.Parameter{
						HeaderParamXTenant, HeaderParamXUser,
						PathParamCluster,
					},
					Results: commonResults,
				},
			},
		},
		{
			Path: path.Join(constants.RootPath, fmt.Sprintf("/clusters/{%s}/loadbalancers", constants.ParameterCluster)),
			Definitions: []definition.Definition{
				{
					Description: "get cluster loadbalancers summary",
					Method:      definition.Get,
					Function:    HandleGetLoadBalancersSummary(c),
					Consumes:    []string{definition.MIMEAll}, Produces: []string{definition.MIMEJSON},
					Parameters: []definition.Parameter{
						HeaderParamXTenant, HeaderParamXUser,
						PathParamCluster,
					},
					Results: commonResults,
				},
			},
		},
		{
			Path: path.Join(constants.RootPath, fmt.Sprintf("/clusters/{%s}/storage", constants.ParameterCluster)),
			Definitions: []definition.Definition{
				{
					Description: "list cluster storage info",
					Method:      definition.List,
					Function:    HandleListStorage(c),
					Consumes:    []string{definition.MIMEAll}, Produces: []string{definition.MIMEJSON},
					Parameters: []definition.Parameter{
						HeaderParamXTenant, HeaderParamXUser,
						PathParamCluster,
						QueryParamStart, QueryParamLimit,
					},
					Results: commonResults,
				},
			},
		},
		{
			Path: path.Join(constants.RootPath, fmt.Sprintf("/ci")),
			Definitions: []definition.Definition{
				{
					Description: "get ci summary",
					Method:      definition.Get,
					Function:    HandleGetContinuousIntegrationSummary(c),
					Consumes:    []string{definition.MIMEAll}, Produces: []string{definition.MIMEJSON},
					Parameters: []definition.Parameter{
						HeaderParamXTenant, HeaderParamXUser,
					},
					Results: commonResults,
				},
			},
		},
		{
			Path: path.Join(constants.RootPath, fmt.Sprintf("/cargo")),
			Definitions: []definition.Definition{
				{
					Description: "get cargo summary",
					Method:      definition.Get,
					Function:    HandleGetCargoInfo(c),
					Consumes:    []string{definition.MIMEAll}, Produces: []string{definition.MIMEJSON},
					Parameters: []definition.Parameter{
						HeaderParamXTenant, HeaderParamXUser,
						QueryParamStart, QueryParamLimit,
					},
					Results: commonResults,
				},
			},
		},
		{
			Path: path.Join(constants.RootPath, fmt.Sprintf("/events")),
			Definitions: []definition.Definition{
				{
					Description: "get events",
					Method:      definition.List,
					Function:    HandleListEvent(c),
					Consumes:    []string{definition.MIMEAll}, Produces: []string{definition.MIMEJSON},
					Parameters: []definition.Parameter{
						HeaderParamXTenant, HeaderParamXUser,
						QueryParamStart, QueryParamLimit,
					},
					Results: commonResults,
				},
			},
		},
		{
			Path: path.Join(constants.RootPath, fmt.Sprintf("/clusters/{%s}/addonhealth", constants.ParameterCluster)),
			Definitions: []definition.Definition{
				{
					Description: "get cluster addon health summary",
					Method:      definition.Get,
					Function:    HandleGetAddonHealthSummary(c),
					Consumes:    []string{definition.MIMEAll}, Produces: []string{definition.MIMEJSON},
					Parameters: []definition.Parameter{
						HeaderParamXTenant, HeaderParamXUser,
						PathParamCluster,
					},
					Results: commonResults,
				},
			},
		},
		{
			Path: path.Join(constants.RootPath, fmt.Sprintf("/clusters/{%s}/kubehealth", constants.ParameterCluster)),
			Definitions: []definition.Definition{
				{
					Description: "get cluster kube health summary",
					Method:      definition.Get,
					Function:    HandleGetKubeHealthSummary(c),
					Consumes:    []string{definition.MIMEAll}, Produces: []string{definition.MIMEJSON},
					Parameters: []definition.Parameter{
						HeaderParamXTenant, HeaderParamXUser,
						PathParamCluster,
					},
					Results: commonResults,
				},
			},
		},
		{
			Path: path.Join(constants.RootPath, fmt.Sprintf("/alerts")),
			Definitions: []definition.Definition{
				{
					Description: "get alert summary",
					Method:      definition.Get,
					Function:    HandleGetAlertSummary(c),
					Consumes:    []string{definition.MIMEAll}, Produces: []string{definition.MIMEJSON},
					Parameters: []definition.Parameter{
						HeaderParamXTenant, HeaderParamXUser,
					},
					Results: commonResults,
				},
			},
		},
		{
			Path: path.Join(constants.RootPath, fmt.Sprintf("/platform")),
			Definitions: []definition.Definition{
				{
					Description: "get platform summary",
					Method:      definition.Get,
					Function:    HandleGetPlatformSummary(c),
					Consumes:    []string{definition.MIMEAll}, Produces: []string{definition.MIMEJSON},
					Parameters: []definition.Parameter{
						HeaderParamXTenant, HeaderParamXUser,
					},
					Results: commonResults,
				},
			},
		},
		{
			Path: path.Join(constants.RootPath, fmt.Sprintf("/clusters/{%s}/apps", constants.ParameterCluster)),
			Definitions: []definition.Definition{
				{
					Description: "get platform summary",
					Method:      definition.Get,
					Function:    HandleGetAppSummary(c),
					Consumes:    []string{definition.MIMEAll}, Produces: []string{definition.MIMEJSON},
					Parameters: []definition.Parameter{
						HeaderParamXTenant, HeaderParamXUser,
						PathParamCluster,
					},
					Results: commonResults,
				},
			},
		},
	}
}
