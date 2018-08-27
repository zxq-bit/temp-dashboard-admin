package helper

import (
	tntv1al "github.com/caicloud/clientset/pkg/apis/tenant/v1alpha1"
)

func GetPartitionTenant(p *tntv1al.Partition) (tenant string, ok bool) {
	if p == nil || p.Labels == nil {
		return "", false
	}
	tenant, ok = p.Annotations[LabelKeyTenant]
	return
}
