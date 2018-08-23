package crd

import (
	resv1b1 "github.com/caicloud/clientset/pkg/apis/resource/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	"github.com/caicloud/dashboard-admin/pkg/errors"
	"github.com/caicloud/dashboard-admin/pkg/kubernetes"
)

const (
	CacheNameMachine = "Machine"
)

func (scc *subClusterCaches) GetMachineCache() (*MachinesCache, bool) {
	return scc.GetAsMachineCache(CacheNameMachine)
}
func (scc *subClusterCaches) GetAsMachineCache(name string) (*MachinesCache, bool) {
	c, ok := scc.m[name]
	if ok {
		return &MachinesCache{lwCache: c, kc: scc.kc}, true
	}
	return nil, false
}

type MachinesCache struct {
	lwCache *ListWatchCache
	kc      kubernetes.Interface
}

func NewMachinesCache(kc kubernetes.Interface) (*MachinesCache, error) {
	listWatcher, objType := GetMachineCacheConfig(kc)
	c, e := NewListWatchCache(listWatcher, objType)
	if e != nil {
		return nil, e
	}
	return &MachinesCache{
		lwCache: c,
		kc:      kc,
	}, nil
}

func (tc *MachinesCache) Run(stopCh chan struct{}) {
	tc.lwCache.Run(stopCh)
}

func (tc *MachinesCache) Get(key string) (*resv1b1.Machine, error) {
	return CacheGetMachine(key, tc.lwCache.indexer, tc.kc)
}
func (tc *MachinesCache) List() ([]resv1b1.Machine, error) {
	return CacheListMachines(tc.lwCache.indexer, tc.kc)
}
func (tc *MachinesCache) ListCachePointer() (re []*resv1b1.Machine) {
	return CacheListMachinesPointer(tc.lwCache.indexer, tc.kc)
}

func (tc *MachinesCache) Indexes() cache.Indexer {
	return tc.lwCache.indexer
}

func GetMachineCacheConfig(kc kubernetes.Interface) (cache.ListerWatcher, runtime.Object) {
	return &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			options.FieldSelector = fields.Everything().String()
			return kc.ResourceV1beta1().Machines().List(options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			options.FieldSelector = fields.Everything().String()
			options.Watch = true
			return kc.ResourceV1beta1().Machines().Watch(options)
		},
	}, &resv1b1.Machine{}
}

func CacheGetMachine(key string, indexer cache.Indexer, kc kubernetes.Interface) (*resv1b1.Machine, error) {
	if indexer != nil {
		if obj, exist, e := indexer.GetByKey(key); exist && obj != nil && e == nil {
			if machine, _ := obj.(*resv1b1.Machine); machine != nil && machine.Name == key {
				return machine, nil
			}
		}
	}
	if kc == nil {
		return nil, errors.ErrVarKubeClientNil
	}
	machine, e := kc.ResourceV1beta1().Machines().Get(key, metav1.GetOptions{})
	if e != nil {
		return nil, e
	}
	return machine, nil
}

func CacheListMachines(indexer cache.Indexer, kc kubernetes.Interface) ([]resv1b1.Machine, error) {
	if items := indexer.List(); len(items) > 0 {
		re := make([]resv1b1.Machine, 0, len(items))
		for _, obj := range items {
			machine, _ := obj.(*resv1b1.Machine)
			if machine != nil {
				re = append(re, *machine)
			}
		}
		if len(re) > 0 {
			return re, nil
		}
	}
	machineList, e := kc.ResourceV1beta1().Machines().List(metav1.ListOptions{})
	if e != nil {
		return nil, e
	}
	return machineList.Items, nil
}

func CacheListMachinesPointer(indexer cache.Indexer, kc kubernetes.Interface) (re []*resv1b1.Machine) {
	// from cache
	items := indexer.List()
	if len(items) > 0 {
		re = make([]*resv1b1.Machine, 0, len(items))
		for _, obj := range items {
			machine, _ := obj.(*resv1b1.Machine)
			if machine != nil {
				re = append(re, machine)
			}
		}
	}
	if len(re) > 0 {
		return re
	}
	// from source
	machineList, e := kc.ResourceV1beta1().Machines().List(metav1.ListOptions{})
	if e != nil || len(machineList.Items) == 0 {
		return nil
	}
	re = make([]*resv1b1.Machine, len(machineList.Items))
	for i := range machineList.Items {
		re[i] = &machineList.Items[i]
	}
	return re
}
