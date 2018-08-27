package rest

import (
	"context"
	"fmt"
	"time"

	"github.com/caicloud/nirvana/log"

	"github.com/caicloud/dashboard-admin/pkg/admin/helper"
	apiv1a1 "github.com/caicloud/dashboard-admin/pkg/apis/v1alpha1"
	"github.com/caicloud/dashboard-admin/pkg/cache"
	"github.com/caicloud/dashboard-admin/pkg/errors"
	"github.com/caicloud/dashboard-admin/pkg/util"
)

func HandleListClusterInfo(c *cache.Cache) func(ctx context.Context,
	xTenant, xUser string, start, limit int) (*apiv1a1.ClusterInfoList, error) {
	return func(ctx context.Context, xTenant, xUser string, start, limit int) (*apiv1a1.ClusterInfoList, error) {
		logPrefix := fmt.Sprintf("HandleListClusterInfo[%v:%v][%v:%v]", xTenant, xUser, start, limit)
		startTime := time.Now()
		log.Infof("%s start", logPrefix)
		if fe := handleListClusterInfoPrework(xTenant, xUser, start, limit); fe != nil {
			log.Errorf("%s handleListClusterInfoPrework failed, %v", fe.Error())
			return nil, fe
		}

		cis, e := helper.ListClusterInfo(c, xTenant)
		if e != nil {
			log.Errorf("%s failed, %v", logPrefix, e)
			return nil, errors.NewError().SetErrorInternalServerError(e)
		}
		log.Infof("%s done in %v", logPrefix, time.Now().Sub(startTime))

		end := util.GetStartLimitEnd(start, limit, len(cis))
		return &apiv1a1.ClusterInfoList{
			MetaData: apiv1a1.ListMetaData{Total: len(cis)},
			Items:    cis[start:end],
		}, nil
	}
}

func HandleGetMachineSummary(c *cache.Cache) func(ctx context.Context,
	xTenant, xUser, cluster string) (*apiv1a1.MachineSummary, error) {
	return func(ctx context.Context, xTenant, xUser, cluster string) (*apiv1a1.MachineSummary, error) {
		return &apiv1a1.MachineSummary{}, nil
	}
}

func HandleGetLoadBalancersSummary(c *cache.Cache) func(ctx context.Context,
	xTenant, xUser, cluster string) (*apiv1a1.LoadBalancersSummary, error) {
	return func(ctx context.Context, xTenant, xUser, cluster string) (*apiv1a1.LoadBalancersSummary, error) {
		return &apiv1a1.LoadBalancersSummary{}, nil
	}
}

func HandleListStorage(c *cache.Cache) func(ctx context.Context,
	xTenant, xUser, cluster string, start, limit int) (*apiv1a1.StorageClassList, error) {
	return func(ctx context.Context, xTenant, xUser, cluster string, start, limit int) (*apiv1a1.StorageClassList, error) {
		return &apiv1a1.StorageClassList{}, nil
	}
}

func HandleGetContinuousIntegrationSummary(c *cache.Cache) func(ctx context.Context,
	xTenant, xUser string) (*apiv1a1.ContinuousIntegrationSummary, error) {
	return func(ctx context.Context, xTenant, xUser string) (*apiv1a1.ContinuousIntegrationSummary, error) {
		return &apiv1a1.ContinuousIntegrationSummary{}, nil
	}
}

func HandleGetCargoInfo(c *cache.Cache) func(ctx context.Context,
	xTenant, xUser string) (*apiv1a1.RegistryInfoList, error) {
	return func(ctx context.Context, xTenant, xUser string) (*apiv1a1.RegistryInfoList, error) {
		return &apiv1a1.RegistryInfoList{}, nil
	}
}

func HandleListEvent(c *cache.Cache) func(ctx context.Context,
	xTenant, xUser string, start, limit int) (*apiv1a1.EventList, error) {
	return func(ctx context.Context, xTenant, xUser string, start, limit int) (*apiv1a1.EventList, error) {
		return &apiv1a1.EventList{}, nil
	}
}

func HandleGetAddonHealthSummary(c *cache.Cache) func(ctx context.Context,
	xTenant, xUser, cluster string) (*apiv1a1.AddonHealthSummary, error) {
	return func(ctx context.Context, xTenant, xUser, cluster string) (*apiv1a1.AddonHealthSummary, error) {
		return &apiv1a1.AddonHealthSummary{}, nil
	}
}

func HandleGetKubeHealthSummary(c *cache.Cache) func(ctx context.Context,
	xTenant, xUser, cluster string) (*apiv1a1.KubeHealthSummary, error) {
	return func(ctx context.Context, xTenant, xUser, cluster string) (*apiv1a1.KubeHealthSummary, error) {
		return &apiv1a1.KubeHealthSummary{}, nil
	}
}

func HandleGetAlertSummary(c *cache.Cache) func(ctx context.Context,
	xTenant, xUser string) (*apiv1a1.AlertSummary, error) {
	return func(ctx context.Context, xTenant, xUser string) (*apiv1a1.AlertSummary, error) {
		return &apiv1a1.AlertSummary{}, nil
	}
}

func HandleGetPlatformSummary(c *cache.Cache) func(ctx context.Context,
	xTenant, xUser string) (*apiv1a1.PlatformSummary, error) {
	return func(ctx context.Context, xTenant, xUser string) (*apiv1a1.PlatformSummary, error) {
		return &apiv1a1.PlatformSummary{}, nil
	}
}

func HandleGetAppSummary(c *cache.Cache) func(ctx context.Context,
	xTenant, xUser, cluster string) (*apiv1a1.AppSummary, error) {
	return func(ctx context.Context, xTenant, xUser, cluster string) (*apiv1a1.AppSummary, error) {
		return &apiv1a1.AppSummary{}, nil
	}
}
