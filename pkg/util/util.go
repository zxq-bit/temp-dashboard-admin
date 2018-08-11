package util

import "github.com/caicloud/dashboard-admin/pkg/kubernetes"

func DoUtilNotConflict(f func() error) error {
	e := f()
	for kubernetes.IsConflict(e) {
		e = f()
	}
	return e
}
