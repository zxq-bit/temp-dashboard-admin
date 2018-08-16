package helper

import (
	"github.com/caicloud/dashboard-admin/pkg/errors"
	"github.com/caicloud/dashboard-admin/pkg/kubernetes"
)

type Content struct {
	kc kubernetes.Interface
}

func NewContent(kc kubernetes.Interface) (*Content, error) {
	if kc == nil {
		return nil, errors.ErrVarKubeClientNil
	}
	return &Content{kc: kc}, nil
}
