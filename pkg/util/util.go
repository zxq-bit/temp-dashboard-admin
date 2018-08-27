package util

import "github.com/caicloud/dashboard-admin/pkg/kubernetes"

func DoUtilNotConflict(f func() error) error {
	e := f()
	for kubernetes.IsConflict(e) {
		e = f()
	}
	return e
}

func GetStartLimitEnd(start, limit, arrayLen int) (end int) {
	if limit == 0 { // no limit
		end = arrayLen
	} else {
		end = start + limit
		if end > arrayLen {
			end = arrayLen
		}
	}
	return end
}
