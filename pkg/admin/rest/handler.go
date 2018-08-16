package rest

import (
	"context"

	"github.com/caicloud/dashboard-admin/pkg/admin/helper"
	apiv1a1 "github.com/caicloud/dashboard-admin/pkg/apis/v1alpha1"
)

func HandleListClusterInfo(c *helper.Content) func(ctx context.Context,
	xTenant, xUser string, start, limit int) (*apiv1a1.ClusterInfoList, error) {
	return func(ctx context.Context, xTenant, xUser string, start, limit int) (*apiv1a1.ClusterInfoList, error) {
		return &apiv1a1.ClusterInfoList{}, nil
	}
}

func HandleGetMachineSummary(c *helper.Content) func(ctx context.Context,
	xTenant, xUser, cluster string) (*apiv1a1.MachineSummary, error) {
	return func(ctx context.Context, xTenant, xUser, cluster string) (*apiv1a1.MachineSummary, error) {
		return &apiv1a1.MachineSummary{}, nil
	}
}

func HandleGetLoadBalancersSummary(c *helper.Content) func(ctx context.Context,
	xTenant, xUser, cluster string) (*apiv1a1.LoadBalancersSummary, error) {
	return func(ctx context.Context, xTenant, xUser, cluster string) (*apiv1a1.LoadBalancersSummary, error) {
		return &apiv1a1.LoadBalancersSummary{}, nil
	}
}

func HandleGetContinuousIntegrationSummary(c *helper.Content) func(ctx context.Context,
	xTenant, xUser string) (*apiv1a1.ContinuousIntegrationSummary, error) {
	return func(ctx context.Context, xTenant, xUser string) (*apiv1a1.ContinuousIntegrationSummary, error) {
		return &apiv1a1.ContinuousIntegrationSummary{}, nil
	}
}

func HandleGetCargoInfo(c *helper.Content) func(ctx context.Context,
	xTenant, xUser string) (*apiv1a1.CargoInfo, error) {
	return func(ctx context.Context, xTenant, xUser string) (*apiv1a1.CargoInfo, error) {
		return &apiv1a1.CargoInfo{}, nil
	}
}

func HandleListEvent(c *helper.Content) func(ctx context.Context,
	xTenant, xUser string, start, limit int) (*apiv1a1.EventList, error) {
	return func(ctx context.Context, xTenant, xUser string, start, limit int) (*apiv1a1.EventList, error) {
		return &apiv1a1.EventList{}, nil
	}
}

func HandleGetAddonHealthSummary(c *helper.Content) func(ctx context.Context,
	xTenant, xUser, cluster string) (*apiv1a1.AddonHealthSummary, error) {
	return func(ctx context.Context, xTenant, xUser, cluster string) (*apiv1a1.AddonHealthSummary, error) {
		return &apiv1a1.AddonHealthSummary{}, nil
	}
}

func HandleGetKubeHealthSummary(c *helper.Content) func(ctx context.Context,
	xTenant, xUser, cluster string) (*apiv1a1.KubeHealthSummary, error) {
	return func(ctx context.Context, xTenant, xUser, cluster string) (*apiv1a1.KubeHealthSummary, error) {
		return &apiv1a1.KubeHealthSummary{}, nil
	}
}

func HandleGetAlertSummary(c *helper.Content) func(ctx context.Context,
	xTenant, xUser string) (*apiv1a1.AlertSummary, error) {
	return func(ctx context.Context, xTenant, xUser string) (*apiv1a1.AlertSummary, error) {
		return &apiv1a1.AlertSummary{}, nil
	}
}

func HandleGetPlatformSummary(c *helper.Content) func(ctx context.Context,
	xTenant, xUser string) (*apiv1a1.PlatformSummary, error) {
	return func(ctx context.Context, xTenant, xUser string) (*apiv1a1.PlatformSummary, error) {
		return &apiv1a1.PlatformSummary{}, nil
	}
}

func HandleGetAppSummary(c *helper.Content) func(ctx context.Context,
	xTenant, xUser, cluster string) (*apiv1a1.AppSummary, error) {
	return func(ctx context.Context, xTenant, xUser, cluster string) (*apiv1a1.AppSummary, error) {
		return &apiv1a1.AppSummary{}, nil
	}
}
