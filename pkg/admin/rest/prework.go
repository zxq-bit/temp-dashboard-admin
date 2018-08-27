package rest

import (
	"github.com/caicloud/dashboard-admin/pkg/errors"
)

func handleListClusterInfoPrework(xTenant, xUser string, start, limit int) *errors.FormatError {
	if fe := ParamCheckTenantAndUser(xTenant, xUser); fe != nil {
		return fe
	}
	if fe := ParamCheckStartAndLimit(start, limit); fe != nil {
		return fe
	}
	return nil
}
