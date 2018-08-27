package rest

import (
	"github.com/caicloud/dashboard-admin/pkg/errors"
)

func listPrework(xTenant, xUser string, start, limit int) *errors.FormatError {
	if fe := ParamCheckTenantAndUser(xTenant, xUser); fe != nil {
		return fe
	}
	if fe := ParamCheckStartAndLimit(start, limit); fe != nil {
		return fe
	}
	return nil
}

func getClusterSubPrework(xTenant, xUser, cluster string) *errors.FormatError {
	if fe := ParamCheckTenantAndUser(xTenant, xUser); fe != nil {
		return fe
	}
	if len(cluster) == 0 {
		return errors.NewError().SetErrorEmptyCluster()
	}
	return nil
}

func getClusterAcrossPrework(xTenant, xUser string) *errors.FormatError {
	if fe := ParamCheckTenantAndUser(xTenant, xUser); fe != nil {
		return fe
	}
	return nil
}

func handleListClusterInfoPrework(xTenant, xUser string, start, limit int) *errors.FormatError {
	return listPrework(xTenant, xUser, start, limit)
}

func handleGetMachineSummaryPrework(xTenant, xUser, cluster string) *errors.FormatError {
	return getClusterSubPrework(xTenant, xUser, cluster)
}

func handleGetLoadBalancersSummaryPrework(xTenant, xUser, cluster string) *errors.FormatError {
	return getClusterSubPrework(xTenant, xUser, cluster)
}

func handleListStoragePrework(xTenant, xUser string, start, limit int) *errors.FormatError {
	return listPrework(xTenant, xUser, start, limit)
}

func handleGetContinuousIntegrationSummaryPrework(xTenant, xUser string) *errors.FormatError {
	return getClusterAcrossPrework(xTenant, xUser)
}

func handleGetCargoInfoPrework(xTenant, xUser string, start, limit int) *errors.FormatError {
	return listPrework(xTenant, xUser, start, limit)
}

func handleListEventPrework(xTenant, xUser string, start, limit int) *errors.FormatError {
	return listPrework(xTenant, xUser, start, limit)
}

func handleGetAddonHealthSummaryPrework(xTenant, xUser, cluster string) *errors.FormatError {
	return getClusterSubPrework(xTenant, xUser, cluster)
}

func handleGetKubeHealthSummaryPrework(xTenant, xUser, cluster string) *errors.FormatError {
	return getClusterSubPrework(xTenant, xUser, cluster)
}

func handleGetAlertSummaryPrework(xTenant, xUser string) *errors.FormatError {
	return getClusterAcrossPrework(xTenant, xUser)
}

func handleGetPlatformSummaryPrework(xTenant, xUser string) *errors.FormatError {
	return getClusterAcrossPrework(xTenant, xUser)
}

func handleGetAppSummaryPrework(xTenant, xUser, cluster string) *errors.FormatError {
	return getClusterSubPrework(xTenant, xUser, cluster)
}
