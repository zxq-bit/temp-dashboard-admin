package server

import (
	"fmt"

	"github.com/caicloud/nirvana"
	"github.com/caicloud/nirvana/config"
	"github.com/caicloud/nirvana/log"

	"github.com/caicloud/dashboard-admin/pkg/admin/helper"
	"github.com/caicloud/dashboard-admin/pkg/admin/rest"
	"github.com/caicloud/dashboard-admin/pkg/constants"
	"github.com/caicloud/dashboard-admin/pkg/kubernetes"
)

type Server struct {
	cfg Config
	cmd config.NirvanaCommand

	stopCh chan struct{}

	c *helper.Content
}

func NewServer() (*Server, error) {
	s := &Server{
		cfg: Config{
			KubeHost:   constants.DefaultKubeHost,
			KubeConfig: constants.DefaultKubeConfig,
		},
		cmd: config.NewNirvanaCommand(&config.Option{
			Port: uint16(constants.DefaultListenPort),
		}),
		stopCh: make(chan struct{}),
	}
	s.cmd.AddOption("", &s.cfg)
	s.cmd.SetHook(&config.NirvanaCommandHookFunc{
		PreConfigureFunc: s.init,
	})
	return s, nil
}

func (s *Server) init(config *nirvana.Config) error {
	// config
	kubeHost := s.cfg.KubeHost
	kubeConfig := s.cfg.KubeConfig
	e := s.cfg.Validate()
	if e != nil {
		log.Errorf("Validate config %v failed, %v", s.cfg.String(), e)
		return e
	}
	log.Info(s.cfg.String())

	// kube
	kc, e := kubernetes.NewClientFromFlags(kubeHost, kubeConfig)
	if e != nil {
		return fmt.Errorf("NewClientFromFlags failed, %v", e)
	}

	// helper
	c, e := helper.NewContent(kc)
	if e != nil {
		return fmt.Errorf("NewContent failed, %v", e)
	}

	// descriptor
	config.Configure(
		nirvana.Descriptor(rest.InitNirvanaDescriptors(c)...),
	)
	return nil
}

func (s *Server) Run() error {
	return s.cmd.Execute()
}
