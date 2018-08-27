package fake

import (
	"fmt"
	"strconv"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	apiv1a1 "github.com/caicloud/dashboard-admin/pkg/apis/v1alpha1"
)

func ListClusterInfo() []apiv1a1.ClusterInfo {
	count := 5
	ctrlID := 0

	isCtrl := func(index int) bool { return index == ctrlID }
	getUUID := func(index int) string { return "uuid-" + strconv.Itoa(index) }
	getName := func(index int) string {
		if isCtrl(index) {
			return "ctrl-" + strconv.Itoa(index)
		} else {
			return "user-" + strconv.Itoa(index)
		}
	}
	getAlias := func(index int) string { return "alias-" + getName(index) }

	re := make([]apiv1a1.ClusterInfo, count)
	for i := range re {
		ci := &re[i]

		ci.Metadata = apiv1a1.ObjectMetaData{
			ID:    getUUID(i),
			Name:  getName(i),
			Alias: getAlias(i),
		}
		ci.Physical = &apiv1a1.Physical{
			Capacity: corev1.ResourceList{
				"cpu":    *resource.NewQuantity(int64(i+1)*8, resource.DecimalSI),
				"memory": *resource.NewQuantity(int64(i+1)*16*1024*1024*1024, resource.BinarySI),
			},
			Used: corev1.ResourceList{
				"cpu":    *resource.NewQuantity(int64(i+1)*3, resource.DecimalSI),
				"memory": *resource.NewQuantity(int64(i+1)*6*1024*1024*1024, resource.BinarySI),
			},
		}
		ci.Request = apiv1a1.Logical{
			Capacity: corev1.ResourceList{
				"cpu":    *resource.NewQuantity(int64(i+1)*24, resource.DecimalSI),
				"memory": *resource.NewQuantity(int64(i+1)*15*1024*1024*1024, resource.BinarySI),
			},
			SystemUsed: corev1.ResourceList{
				"cpu":    *resource.NewQuantity(int64(i+1)*1, resource.DecimalSI),
				"memory": *resource.NewQuantity(int64(i+1)*1*1024*1024*1024, resource.BinarySI),
			},
			UserUsed: corev1.ResourceList{
				"cpu":    *resource.NewQuantity(int64(i+1)*2, resource.DecimalSI),
				"memory": *resource.NewQuantity(int64(i+1)*5*1024*1024*1024, resource.BinarySI),
			},
		}
		ci.Limit = apiv1a1.Logical{
			Capacity: corev1.ResourceList{
				"cpu":    *resource.NewQuantity(int64(i+1)*7, resource.DecimalSI),
				"memory": *resource.NewQuantity(int64(i+1)*15*1024*1024*1024, resource.BinarySI),
			},
			SystemUsed: corev1.ResourceList{
				"cpu":    *resource.NewQuantity(int64(i+1)*1, resource.DecimalSI),
				"memory": *resource.NewQuantity(int64(i+1)*1*1024*1024*1024, resource.BinarySI),
			},
			UserUsed: corev1.ResourceList{
				"cpu":    *resource.NewQuantity(int64(i+1)*2, resource.DecimalSI),
				"memory": *resource.NewQuantity(int64(i+1)*5*1024*1024*1024, resource.BinarySI),
			},
		}
		ci.NodeNum = (i + 1) * 3
		ci.AppNum = (i + 1) * 4
		ci.PodNum = (i + 1) * 5
		ci.IsControl = isCtrl(i)
	}
	return re
}

func GetMachineSummary() *apiv1a1.MachineSummary {
	count := 5
	masterID := 3
	isMaster := func(index int) bool { return index == masterID }
	re := &apiv1a1.MachineSummary{
		NormalNum:   9,
		AbnormalNum: 3,
		OfflineNum:  1,
		MaxLoads:    make([]apiv1a1.MachineLoad, count),
	}
	for i := range re.MaxLoads {
		re.MaxLoads[i].IP = fmt.Sprintf("192.168.1.%d", i+1)
		re.MaxLoads[i].Score = count - i
		re.MaxLoads[i].IsMaster = isMaster(i)
	}
	return re
}

func GetLoadBalancersSummary() *apiv1a1.LoadBalancersSummary {
	count := 5
	isSystemID := 3
	isSystem := func(index int) bool { return index == isSystemID }
	re := &apiv1a1.LoadBalancersSummary{
		NormalNum:   7,
		AbnormalNum: 2,
		TopIO:       make([]apiv1a1.LoadBalancerIO, count),
	}
	for i := range re.TopIO {
		re.TopIO[i].Name = fmt.Sprintf("lb-%d", i+1)
		re.TopIO[i].Namespace = fmt.Sprintf("ns-%d", count-i)
		re.TopIO[i].IsSystem = isSystem(i)
		re.TopIO[i].In = uint64(count-i) * 65536
		re.TopIO[i].Out = uint64(count-i) * 32768
	}
	return re
}

func ListStorage() []apiv1a1.StorageClassStatus {
	count := 5
	IsSystemID := 2

	isSystem := func(index int) bool { return index == IsSystemID }
	getName := func(index int) string { return "sc-" + strconv.Itoa(index) }
	getAlias := func(index int) string { return "alias-" + getName(index) }

	re := make([]apiv1a1.StorageClassStatus, count)
	for i := range re {
		ss := &re[i]

		ss.Name = getName(i)
		ss.Alias = getAlias(i)
		ss.IsSystem = isSystem(i)
		ss.Capacity.Num = *resource.NewQuantity(int64(count-i), resource.DecimalSI)
		ss.Capacity.Size = *resource.NewQuantity(int64(count-i)*10*1024*1024, resource.BinarySI)
		ss.Used.Num = *resource.NewQuantity(int64(count-i)/2, resource.DecimalSI)
		ss.Used.Size = *resource.NewQuantity(int64(count-i)*10*1024*1024/2, resource.BinarySI)
	}
	return re
}

func GetContinuousIntegrationSummary() *apiv1a1.ContinuousIntegrationSummary {
	return &apiv1a1.ContinuousIntegrationSummary{
		PipelineNum:  66,
		WorkspaceNum: 8,
	}
}

func ListRegistryInfo() []apiv1a1.RegistryInfo {
	count := 5

	re := make([]apiv1a1.RegistryInfo, count)
	for i := range re {
		ss := &re[i]

		ss.Name = fmt.Sprintf("registry-%d", count-i)
		ss.ProjectNum = (i + 1) * 5
		ss.ImageNum = (i + 1) * (i + 1) * (i + 1)
		ss.DiskUsage = fmt.Sprintf("%d%d%%", count-i, i)
	}
	return re
}

func ListEvent() []apiv1a1.Event {
	count := 5

	getType := func(index int) string {
		switch index % 3 {
		case 0:
			return "create"
		case 1:
			return "delete"
		default:
			return "update"
		}
	}
	getResult := func(index int) string {
		switch index % 2 {
		case 0:
			return "success"
		default:
			return "failed"
		}
	}
	re := make([]apiv1a1.Event, count)
	for i := range re {
		ss := &re[i]

		ss.Type = getType(i)
		ss.Result = getResult(i)
		ss.Time = time.Now().Add(time.Duration(0-i) * time.Minute)
		ss.User = fmt.Sprintf("uuu-%d", i)
		ss.Tenant = fmt.Sprintf("ttt-%d", i)
		ss.Message = fmt.Sprintf("mmm-%d", i)
	}
	return re
}

func GetAddonHealthSummary() *apiv1a1.AddonHealthSummary {
	return &apiv1a1.AddonHealthSummary{
		AbnormalNum: 3,
		TotalNum:    23,
		Addons: []apiv1a1.Component{
			{Name: "storage", Status: "normal"},
			{Name: "tenant", Status: "normal"},
			{Name: "load-balancer", Status: "normal"},
			{Name: "monitoring", Status: "abnormal"},
			{Name: "logging", Status: "abnormal"},
		},
	}
}

func GetKubeHealthSummary() *apiv1a1.KubeHealthSummary {
	return &apiv1a1.KubeHealthSummary{
		AbnormalNum: 3,
		TotalNum:    23,
		Components: []apiv1a1.Component{
			{Name: "apiserver-provider", Status: "normal"},
			{Name: "kube-dns", Status: "normal"},
			{Name: "canal", Status: "normal"},
			{Name: "hybrid-controller-manager", Status: "abnormal"},
			{Name: "kube-scheduler", Status: "abnormal"},
		},
	}
}

func GetAlertSummary() *apiv1a1.AlertSummary {
	count := 5

	re := &apiv1a1.AlertSummary{
		AlertingRulesNum: 23,
		RecentRecordsNum: 666,
		LatestRecords:    make([]apiv1a1.AlertRecord, count),
	}
	for i := range re.LatestRecords {
		rec := &re.LatestRecords[i]

		rec.Time = time.Now().Add(time.Duration(0-i) * time.Second)
		rec.Message = fmt.Sprintf("message %v", count-i)
	}
	return re
}

func GetPlatformSummary() *apiv1a1.PlatformSummary {
	re := &apiv1a1.PlatformSummary{
		TeamNum:        22,
		UserNum:        33,
		FreeMachineNum: new(int),
	}
	*re.FreeMachineNum = 66
	return re
}

func GetAppSummary() *apiv1a1.AppSummary {
	return &apiv1a1.AppSummary{
		NormalNum:   666,
		UpdatingNum: 33,
		AbnormalNum: 22,
	}
}
