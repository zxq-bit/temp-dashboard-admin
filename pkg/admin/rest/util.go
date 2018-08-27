package rest

import (
	"strconv"

	"github.com/caicloud/dashboard-admin/pkg/errors"
)

func ParamCheckTenantAndUser(xTenant, xUser string) (fe *errors.FormatError) {
	// if len(xTenant) == 0 { // Comment out for fake
	// 	fe = errors.NewError().SetErrorBadTenantOrUser(xTenant, xUser)
	// }
	return
}

func ParamCheckStartAndLimit(start, limit int) (fe *errors.FormatError) {
	if start < 0 || limit < 0 {
		fe = errors.NewError().SetErrorBadPageStartOrLimit(strconv.Itoa(start), strconv.Itoa(limit))
	}
	return
}
