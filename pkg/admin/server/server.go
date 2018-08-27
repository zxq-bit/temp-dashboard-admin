package server

import (
	"fmt"

	"github.com/caicloud/dashboard-admin/pkg/cache"
	"github.com/caicloud/nirvana"
	"github.com/caicloud/nirvana/config"
	"github.com/caicloud/nirvana/log"

	"github.com/caicloud/dashboard-admin/pkg/admin/rest"
	cfg "github.com/caicloud/dashboard-admin/pkg/config"
	"github.com/caicloud/dashboard-admin/pkg/constants"
)

type Server struct {
	cfg cfg.Config
	cmd config.NirvanaCommand

	stopCh chan struct{}

	c *cache.Cache
}

func NewServer() (*Server, error) {
	s := &Server{
		cfg: *cfg.NewDefaultConfig(),
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
	e := s.cfg.Validate()
	if e != nil {
		log.Errorf("Validate config %v failed, %v", s.cfg.String(), e)
		return e
	}
	log.Info(s.cfg.String())

	// helper
	s.c, e = cache.NewCache(&s.cfg)
	if e != nil {
		return fmt.Errorf("NewCache failed, %v", e)
	}
	go s.c.Run(s.stopCh)

	// descriptor
	config.Configure(
		nirvana.Descriptor(rest.InitNirvanaDescriptors(s.c)...),
	)
	return nil
}

func (s *Server) Run() error {
	defer close(s.stopCh)
	return s.cmd.Execute()
}
