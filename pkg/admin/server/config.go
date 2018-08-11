package server

import "encoding/json"

type Config struct {
	// kube
	KubeHost   string `desc:"control cluster kubernetes host"`
	KubeConfig string `desc:"control cluster kubernetes config"`
}

func (c *Config) Validate() error {
	return nil
}
func (c *Config) String() string {
	b, _ := json.MarshalIndent(c, "", "  ")
	return string(b)
}
