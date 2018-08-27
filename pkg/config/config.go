package config

import (
	"encoding/json"
	"fmt"

	"github.com/caicloud/dashboard-admin/pkg/constants"
)

type Config struct {
	// kube
	KubeHost   string `desc:"control cluster kubernetes host"`
	KubeConfig string `desc:"control cluster kubernetes config"`

	// cache
	TimeoutSecond int
	RefreshSecond int

	// hosts
	CauthHost      string
	DevOpAdminHost string
	CargoAdminHost string
}

func NewDefaultConfig() *Config {
	return &Config{
		KubeHost:       constants.DefaultKubeHost,
		KubeConfig:     constants.DefaultKubeConfig,
		TimeoutSecond:  constants.DefaultTimeoutSecond,
		RefreshSecond:  constants.DefaultRefreshSecond,
		CauthHost:      constants.DefaultCauthHost,
		DevOpAdminHost: constants.DefaultDevOpAdminHost,
		CargoAdminHost: constants.DefaultCargoAdminHost,
	}
}

func (c *Config) Validate() error {
	if c.TimeoutSecond < 0 {
		return fmt.Errorf("illegal timeout seconds %d", c.TimeoutSecond)
	}
	if c.RefreshSecond < 1 {
		return fmt.Errorf("illegal refresh seconds %d", c.RefreshSecond)
	}
	if len(c.CauthHost) == 0 {
		return fmt.Errorf("empty cauth host")
	}
	if len(c.DevOpAdminHost) == 0 {
		return fmt.Errorf("empty devop admin host")
	}
	if len(c.CargoAdminHost) == 0 {
		return fmt.Errorf("empty cargo admin host")
	}
	return nil
}
func (c *Config) String() string {
	b, _ := json.MarshalIndent(c, "", "  ")
	return string(b)
}
