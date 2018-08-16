package v1alpha1

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// basic

type ObjectMetaData struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	CreationTime string `json:"creationTime"`

	Alias       string `json:"alias,omitempty"`
	Description string `json:"description,omitempty"`

	DeletionTime       string            `json:"deletionTime,omitempty"`
	LastUpdateTime     string            `json:"lastUpdateTime,omitempty"`
	LastTransitionTime string            `json:"lastTransitionTime,omitempty"`
	ResourceVersion    string            `json:"resourceVersion,omitempty"`
	Annotations        map[string]string `json:"annotations,omitempty"`
	Labels             map[string]string `json:"labels,omitempty"`
}
type ListMetaData struct {
	Total int `json:"total"`
}

// cluster

type ClusterInfo struct {
	// meta
	Metadata ObjectMetaData `json:"metadata"`
	UserID   string         `json:"userID"`
	// sys-admin
	Physical *Physical `json:"physical,omitempty"`
	// all
	Request Logical `json:"request"`
	Limit   Logical `json:"limit"`
	NodeNum int     `json:"nodeNum"`
	AppNum  int     `json:"appNum"`
	PodNum  int     `json:"podNum"`
}

type ClusterInfoList struct {
	MetaData ListMetaData  `json:"metadata"`
	Items    []ClusterInfo `json:"items"`
}

type Resource struct {
	// sys-admin
	Physical *Physical `json:"physical"`
	// all
	Request Logical `json:"request"`
	Limit   Logical `json:"limit"`
}
type Physical struct {
	Capacity corev1.ResourceList `json:"capacity"`
	Used     corev1.ResourceList `json:"used"`
}
type Logical struct {
	Capacity   corev1.ResourceList `json:"capacity"`
	SystemUsed corev1.ResourceList `json:"systemUsed"`
	UserUsed   corev1.ResourceList `json:"userUsed"`
}

// machine

type MachineSummary struct {
	Cluster     string        `json:"cluster"`
	UserID      string        `json:"userID"`
	NormalNum   int           `json:"normalNum"`
	AbnormalNum int           `json:"abnormalNum"`
	OfflineNum  int           `json:"offlineNum"`
	MaxLoads    []MachineLoad `json:"maxLoads"`
}

type MachineLoad struct {
	IP    string `json:"ip"`
	Score int    `json:"score"`
}

// load balancer

type LoadBalancersSummary struct {
	Cluster     string           `json:"cluster"`
	UserID      string           `json:"userID"`
	NormalNum   int              `json:"normalNum"`
	AbnormalNum int              `json:"abnormalNum"`
	TopIO       []LoadBalancerIO `json:"topIO"`
}

type LoadBalancerIO struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	In        uint64 `json:"in"`
	Out       uint64 `json:"out"`
}

// storage

type StorageSummary struct {
	Cluster        string                 `json:"cluster"`
	UserID         string                 `json:"userID"`
	StorageClasses []StorageClassesStatus `json:"storageClasses"`
}

type StorageClassesStatus struct {
	Name     string     `json:"name"`
	Capacity StorageSet `json:"capacity"`
	Used     StorageSet `json:"used"`
}

type StorageSet struct {
	Num  resource.Quantity `json:"num"`
	Size resource.Quantity `json:"size"`
}

// CI

type ContinuousIntegrationSummary struct {
	UserID       string `json:"userID"`
	WorkspaceNum int    `json:"workspaceNum"`
	PipelineNum  int    `json:"pipelineNum"`
}

// cargo

type CargoInfo struct {
	UserID     string         `json:"userID"`
	Registries []RegistryInfo `json:"registries"`
}
type RegistryInfo struct {
	User       string `json:"user"`
	Name       string `json:"name"`
	ProjectNum int    `json:"projectNum"`
	ImageNum   int    `json:"imageNum"`
	DiskUsage  string `json:"diskUsage"`
}

// event

type Event struct {
	Type     string    `json:"type"`
	Result   string    `json:"result"`
	Time     time.Time `json:"time"`
	Operator string    `json:"operator"`
	Message  string    `json:"message"`
}

type EventList struct {
	MetaData ListMetaData `json:"metadata"`
	Items    []Event      `json:"items"`
}

// addon health

type AddonHealthSummary struct {
	AbnormalNum int         `json:"abnormalNum"`
	TotalNum    int         `json:"totalNum"`
	Addons      []Component `json:"addons"`
}

type Component struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

// kube alerts

type KubeHealthSummary struct {
	AbnormalNum int         `json:"abnormalNum"`
	TotalNum    int         `json:"totalNum"`
	Components  []Component `json:"components"`
}

// alert rules

type AlertSummary struct {
	AlertingRulesNum int           `json:"alertingRulesNum"`
	RecentRecordsNum int           `json:"recentRecordsNum"`
	LatestRecords    []AlertRecord `json:"latestRecords"`
}

type AlertRecord struct {
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}

// platform

type PlatformSummary struct {
	Cluster        string `json:"cluster"`
	UserID         string `json:"userID"`
	TeamNum        int    `json:"teamNum"`
	UserNum        int    `json:"userNum"`
	FreeMachineNum *int   `json:"freeMachineNum,omitempty"`
}

// app

type AppSummary struct {
	Cluster     string `json:"cluster"`
	UserID      string `json:"userID"`
	NormalNum   int    `json:"normalNum"`
	UpdatingNum int    `json:"updatingNum"`
	AbnormalNum int    `json:"abnormalNum"`
}
