package kubernetes

import (
	"github.com/caicloud/clientset/kubernetes"
	kubeerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/caicloud/dashboard-admin/pkg/errors"
)

type Interface = kubernetes.Interface
type Config = rest.Config

var BuildConfigFromFlags = clientcmd.BuildConfigFromFlags

func NewClientFromFlags(masterUrl, kubeconfigPath string) (Interface, error) {
	cfg, e := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
	if e != nil {
		return nil, e
	}

	return NewClientFromRestConfig(cfg)
}

func NewClientFromRestConfig(restConf *rest.Config) (Interface, error) {
	cs, e := kubernetes.NewForConfig(restConf)
	if e != nil {
		return nil, e
	}
	return cs, nil
}

func NewClientFromUser(host, user, pwd string) (Interface, error) {
	return NewClientFromRestConfig(&rest.Config{
		Host:     host,
		Username: user,
		Password: pwd,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	})
}

// kube errors

var (
	// IsNotFound      = kubeerr.IsNotFound
	IsConflict      = kubeerr.IsConflict
	IsAlreadyExists = kubeerr.IsAlreadyExists
)

func IsNotFound(err error) bool {
	if !kubeerr.IsNotFound(err) {
		return false
	}
	switch err.(type) {
	case kubeerr.APIStatus:
		ke := err.(kubeerr.APIStatus)
		if ke == nil || ke.Status().Details == nil || len(ke.Status().Details.Causes) == 0 {
			return true
		}
		for _, sd := range ke.Status().Details.Causes {
			if sd.Type == metav1.CauseTypeUnexpectedServerResponse {
				return false
			}
		}
	}
	return true
}

func SwitchKubeUpdateError(name string, e error) *errors.FormatError {
	if IsNotFound(e) {
		return errors.NewError().SetErrorObjectNotFound(name, e)
	}
	if IsConflict(e) {
		return errors.NewError().SetErrorObjectConflict(name, e)
	}
	return errors.NewError().SetErrorInternalServerError(e)
}
func SwitchKubeGetError(name string, e error) *errors.FormatError {
	if IsNotFound(e) {
		return errors.NewError().SetErrorObjectNotFound(name, e)
	}
	return errors.NewError().SetErrorInternalServerError(e)
}
func SwitchKubeCreateError(name string, e error) *errors.FormatError {
	if IsAlreadyExists(e) {
		return errors.NewError().SetErrorObjectAlreadyExist(name, e)
	}
	return errors.NewError().SetErrorInternalServerError(e)
}
