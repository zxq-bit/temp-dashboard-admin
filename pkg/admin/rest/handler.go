package rest

import (
	"context"
	"fmt"
	"time"

	"github.com/caicloud/nirvana/log"

	"github.com/caicloud/dashboard-admin/pkg/admin/fake"
	apiv1a1 "github.com/caicloud/dashboard-admin/pkg/apis/v1alpha1"
	"github.com/caicloud/dashboard-admin/pkg/cache"
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

		// cis, e := helper.ListClusterInfo(c, xTenant)
		// if e != nil {
		// 	log.Errorf("%s failed, %v", logPrefix, e)
		// 	return nil, errors.NewError().SetErrorInternalServerError(e)
		// }
		cis := fake.ListClusterInfo()
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
		logPrefix := fmt.Sprintf("HandleGetMachineSummary[%v:%v][cid:%v]", xTenant, xUser, cluster)
		startTime := time.Now()
		log.Infof("%s start", logPrefix)
		if fe := handleGetMachineSummaryPrework(xTenant, xUser, cluster); fe != nil {
			log.Errorf("%s handleGetMachineSummaryPrework failed, %v", fe.Error())
			return nil, fe
		}

		re := fake.GetMachineSummary()

		log.Infof("%s done in %v", logPrefix, time.Now().Sub(startTime))
		return re, nil
	}
}

func HandleGetLoadBalancersSummary(c *cache.Cache) func(ctx context.Context,
	xTenant, xUser, cluster string) (*apiv1a1.LoadBalancersSummary, error) {
	return func(ctx context.Context, xTenant, xUser, cluster string) (*apiv1a1.LoadBalancersSummary, error) {
		logPrefix := fmt.Sprintf("HandleGetLoadBalancersSummary[%v:%v][cid:%v]", xTenant, xUser, cluster)
		startTime := time.Now()
		log.Infof("%s start", logPrefix)
		if fe := handleGetLoadBalancersSummaryPrework(xTenant, xUser, cluster); fe != nil {
			log.Errorf("%s handleGetLoadBalancersSummaryPrework failed, %v", fe.Error())
			return nil, fe
		}

		re := fake.GetLoadBalancersSummary()

		log.Infof("%s done in %v", logPrefix, time.Now().Sub(startTime))
		return re, nil
	}
}

func HandleListStorage(c *cache.Cache) func(ctx context.Context,
	xTenant, xUser, cluster string, start, limit int) (*apiv1a1.StorageClassList, error) {
	return func(ctx context.Context, xTenant, xUser, cluster string, start, limit int) (*apiv1a1.StorageClassList, error) {
		logPrefix := fmt.Sprintf("HandleListStorage[%v:%v][%v:%v]", xTenant, xUser, start, limit)
		startTime := time.Now()
		log.Infof("%s start", logPrefix)
		if fe := handleListStoragePrework(xTenant, xUser, start, limit); fe != nil {
			log.Errorf("%s handleListStoragePrework failed, %v", fe.Error())
			return nil, fe
		}

		scs := fake.ListStorage()

		log.Infof("%s done in %v", logPrefix, time.Now().Sub(startTime))

		end := util.GetStartLimitEnd(start, limit, len(scs))
		return &apiv1a1.StorageClassList{
			MetaData: apiv1a1.ListMetaData{Total: len(scs)},
			Items:    scs[start:end],
		}, nil
	}
}

func HandleGetContinuousIntegrationSummary(c *cache.Cache) func(ctx context.Context,
	xTenant, xUser string) (*apiv1a1.ContinuousIntegrationSummary, error) {
	return func(ctx context.Context, xTenant, xUser string) (*apiv1a1.ContinuousIntegrationSummary, error) {
		logPrefix := fmt.Sprintf("HandleGetContinuousIntegrationSummary[%v:%v]", xTenant, xUser)
		startTime := time.Now()
		log.Infof("%s start", logPrefix)
		if fe := handleGetContinuousIntegrationSummaryPrework(xTenant, xUser); fe != nil {
			log.Errorf("%s handleGetContinuousIntegrationSummaryPrework failed, %v", fe.Error())
			return nil, fe
		}

		re := fake.GetContinuousIntegrationSummary()

		log.Infof("%s done in %v", logPrefix, time.Now().Sub(startTime))
		return re, nil
	}
}

func HandleGetCargoInfo(c *cache.Cache) func(ctx context.Context,
	xTenant, xUser string, start, limit int) (*apiv1a1.RegistryInfoList, error) {
	return func(ctx context.Context, xTenant, xUser string, start, limit int) (*apiv1a1.RegistryInfoList, error) {
		logPrefix := fmt.Sprintf("HandleGetCargoInfo[%v:%v]", xTenant, xUser)
		startTime := time.Now()
		log.Infof("%s start", logPrefix)
		if fe := handleGetCargoInfoPrework(xTenant, xUser, start, limit); fe != nil {
			log.Errorf("%s handleGetCargoInfoPrework failed, %v", fe.Error())
			return nil, fe
		}

		ris := fake.ListRegistryInfo()

		log.Infof("%s done in %v", logPrefix, time.Now().Sub(startTime))

		end := util.GetStartLimitEnd(start, limit, len(ris))
		return &apiv1a1.RegistryInfoList{
			MetaData: apiv1a1.ListMetaData{Total: len(ris)},
			Items:    ris[start:end],
		}, nil
	}
}

func HandleListEvent(c *cache.Cache) func(ctx context.Context,
	xTenant, xUser string, start, limit int) (*apiv1a1.EventList, error) {
	return func(ctx context.Context, xTenant, xUser string, start, limit int) (*apiv1a1.EventList, error) {
		logPrefix := fmt.Sprintf("HandleListEvent[%v:%v][%v:%v]", xTenant, xUser, start, limit)
		startTime := time.Now()
		log.Infof("%s start", logPrefix)
		if fe := handleListEventPrework(xTenant, xUser, start, limit); fe != nil {
			log.Errorf("%s handleListEventPrework failed, %v", fe.Error())
			return nil, fe
		}

		evs := fake.ListEvent()

		log.Infof("%s done in %v", logPrefix, time.Now().Sub(startTime))

		end := util.GetStartLimitEnd(start, limit, len(evs))
		return &apiv1a1.EventList{
			MetaData: apiv1a1.ListMetaData{Total: len(evs)},
			Items:    evs[start:end],
		}, nil
	}
}

func HandleGetAddonHealthSummary(c *cache.Cache) func(ctx context.Context,
	xTenant, xUser, cluster string) (*apiv1a1.AddonHealthSummary, error) {
	return func(ctx context.Context, xTenant, xUser, cluster string) (*apiv1a1.AddonHealthSummary, error) {
		logPrefix := fmt.Sprintf("HandleGetAddonHealthSummary[%v:%v][cid:%v]", xTenant, xUser, cluster)
		startTime := time.Now()
		log.Infof("%s start", logPrefix)
		if fe := handleGetAddonHealthSummaryPrework(xTenant, xUser, cluster); fe != nil {
			log.Errorf("%s handleGetAddonHealthSummaryPrework failed, %v", fe.Error())
			return nil, fe
		}

		re := fake.GetAddonHealthSummary()

		log.Infof("%s done in %v", logPrefix, time.Now().Sub(startTime))
		return re, nil
	}
}

func HandleGetKubeHealthSummary(c *cache.Cache) func(ctx context.Context,
	xTenant, xUser, cluster string) (*apiv1a1.KubeHealthSummary, error) {
	return func(ctx context.Context, xTenant, xUser, cluster string) (*apiv1a1.KubeHealthSummary, error) {
		logPrefix := fmt.Sprintf("HandleGetKubeHealthSummary[%v:%v][cid:%v]", xTenant, xUser, cluster)
		startTime := time.Now()
		log.Infof("%s start", logPrefix)
		if fe := handleGetKubeHealthSummaryPrework(xTenant, xUser, cluster); fe != nil {
			log.Errorf("%s handleGetKubeHealthSummaryPrework failed, %v", fe.Error())
			return nil, fe
		}

		re := fake.GetKubeHealthSummary()

		log.Infof("%s done in %v", logPrefix, time.Now().Sub(startTime))
		return re, nil
	}
}

func HandleGetAlertSummary(c *cache.Cache) func(ctx context.Context,
	xTenant, xUser string) (*apiv1a1.AlertSummary, error) {
	return func(ctx context.Context, xTenant, xUser string) (*apiv1a1.AlertSummary, error) {
		logPrefix := fmt.Sprintf("HandleGetAlertSummary[%v:%v]", xTenant, xUser)
		startTime := time.Now()
		log.Infof("%s start", logPrefix)
		if fe := handleGetAlertSummaryPrework(xTenant, xUser); fe != nil {
			log.Errorf("%s handleGetAlertSummaryPrework failed, %v", fe.Error())
			return nil, fe
		}

		re := fake.GetAlertSummary()

		log.Infof("%s done in %v", logPrefix, time.Now().Sub(startTime))
		return re, nil
	}
}

func HandleGetPlatformSummary(c *cache.Cache) func(ctx context.Context,
	xTenant, xUser string) (*apiv1a1.PlatformSummary, error) {
	return func(ctx context.Context, xTenant, xUser string) (*apiv1a1.PlatformSummary, error) {
		logPrefix := fmt.Sprintf("HandleGetPlatformSummary[%v:%v]", xTenant, xUser)
		startTime := time.Now()
		log.Infof("%s start", logPrefix)
		if fe := handleGetPlatformSummaryPrework(xTenant, xUser); fe != nil {
			log.Errorf("%s handleGetPlatformSummaryPrework failed, %v", fe.Error())
			return nil, fe
		}

		re := fake.GetPlatformSummary()

		log.Infof("%s done in %v", logPrefix, time.Now().Sub(startTime))
		return re, nil
	}
}

func HandleGetAppSummary(c *cache.Cache) func(ctx context.Context,
	xTenant, xUser, cluster string) (*apiv1a1.AppSummary, error) {
	return func(ctx context.Context, xTenant, xUser, cluster string) (*apiv1a1.AppSummary, error) {
		logPrefix := fmt.Sprintf("HandleGetAppSummary[%v:%v][cid:%v]", xTenant, xUser, cluster)
		startTime := time.Now()
		log.Infof("%s start", logPrefix)
		if fe := handleGetAppSummaryPrework(xTenant, xUser, cluster); fe != nil {
			log.Errorf("%s handleGetAppSummaryPrework failed, %v", fe.Error())
			return nil, fe
		}

		re := fake.GetAppSummary()

		log.Infof("%s done in %v", logPrefix, time.Now().Sub(startTime))
		return re, nil
	}
}
