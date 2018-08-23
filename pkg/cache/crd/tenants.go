package crd

import (
	tntv1al "github.com/caicloud/clientset/pkg/apis/tenant/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"

	"github.com/caicloud/dashboard-admin/pkg/errors"
	"github.com/caicloud/dashboard-admin/pkg/kubernetes"
)

const (
	CacheNameTenant = "Tenant"
)

func (scc *subClusterCaches) GetTenantCache() (*TenantsCache, bool) {
	return scc.GetAsTenantCache(CacheNameTenant)
}
func (scc *subClusterCaches) GetAsTenantCache(name string) (*TenantsCache, bool) {
	c, ok := scc.m[name]
	if ok {
		return &TenantsCache{lwCache: c, kc: scc.kc}, true
	}
	return nil, false
}

type TenantsCache struct {
	lwCache *ListWatchCache
	kc      kubernetes.Interface
}

func NewTenantsCache(kc kubernetes.Interface) (*TenantsCache, error) {
	listWatcher, objType := GetTenantCacheConfig(kc)
	c, e := NewListWatchCache(listWatcher, objType)
	if e != nil {
		return nil, e
	}
	return &TenantsCache{
		lwCache: c,
		kc:      kc,
	}, nil
}

func (tc *TenantsCache) Run(stopCh chan struct{}) {
	tc.lwCache.Run(stopCh)
}

func (tc *TenantsCache) Get(key string) (*tntv1al.Tenant, error) {
	return CacheGetTenant(key, tc.lwCache.indexer, tc.kc)
}
func (tc *TenantsCache) List() ([]tntv1al.Tenant, error) {
	return CacheListTenants(tc.lwCache.indexer, tc.kc)
}
func (tc *TenantsCache) ListCachePointer() (re []*tntv1al.Tenant) {
	return CacheListTenantsPointer(tc.lwCache.indexer, tc.kc)
}

func (tc *TenantsCache) Indexes() cache.Indexer {
	return tc.lwCache.indexer
}

func GetTenantCacheConfig(kc kubernetes.Interface) (cache.ListerWatcher, runtime.Object) {
	return &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			options.FieldSelector = fields.Everything().String()
			return kc.TenantV1alpha1().Tenants().List(options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			options.FieldSelector = fields.Everything().String()
			options.Watch = true
			return kc.TenantV1alpha1().Tenants().Watch(options)
		},
	}, &tntv1al.Tenant{}
}

func CacheGetTenant(key string, indexer cache.Indexer, kc kubernetes.Interface) (*tntv1al.Tenant, error) {
	if indexer != nil {
		if obj, exist, e := indexer.GetByKey(key); exist && obj != nil && e == nil {
			if tenant, _ := obj.(*tntv1al.Tenant); tenant != nil && tenant.Name == key {
				return tenant, nil
			}
		}
	}
	if kc == nil {
		return nil, errors.ErrVarKubeClientNil
	}
	tenant, e := kc.TenantV1alpha1().Tenants().Get(key, metav1.GetOptions{})
	if e != nil {
		return nil, e
	}
	return tenant, nil
}

func CacheListTenants(indexer cache.Indexer, kc kubernetes.Interface) ([]tntv1al.Tenant, error) {
	if items := indexer.List(); len(items) > 0 {
		re := make([]tntv1al.Tenant, 0, len(items))
		for _, obj := range items {
			tenant, _ := obj.(*tntv1al.Tenant)
			if tenant != nil {
				re = append(re, *tenant)
			}
		}
		if len(re) > 0 {
			return re, nil
		}
	}
	tenantList, e := kc.TenantV1alpha1().Tenants().List(metav1.ListOptions{})
	if e != nil {
		return nil, e
	}
	return tenantList.Items, nil
}

func CacheListTenantsPointer(indexer cache.Indexer, kc kubernetes.Interface) (re []*tntv1al.Tenant) {
	// from cache
	items := indexer.List()
	if len(items) > 0 {
		re = make([]*tntv1al.Tenant, 0, len(items))
		for _, obj := range items {
			tenant, _ := obj.(*tntv1al.Tenant)
			if tenant != nil {
				re = append(re, tenant)
			}
		}
	}
	if len(re) > 0 {
		return re
	}
	// from source
	tenantList, e := kc.TenantV1alpha1().Tenants().List(metav1.ListOptions{})
	if e != nil || len(tenantList.Items) == 0 {
		return nil
	}
	re = make([]*tntv1al.Tenant, len(tenantList.Items))
	for i := range tenantList.Items {
		re[i] = &tenantList.Items[i]
	}
	return re
}
