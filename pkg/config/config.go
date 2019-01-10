package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Specification contains application parameters
// configured from environment variables.
type Specification struct {
	Debug      bool   `json:"debug" envconfig:"debug" default:"true" desc:"Application Debug Mode"`
	Port       int    `json:"port" envconfig:"port" required:"true" desc:"HTTP Server Port"`
	SiteURL    string `json:"site_url" envconfig:"site_url" required:"true" desc:"Site URL"`
	SessionKey string `json:"session_key" envconfig:"session_key" required:"true" desc:"Session key"`
	SecretKey  string `json:"secret_key" envconfig:"secret_key" required:"true" desc:"Secret key used in CSRF"`
	Yandex     Yandex `json:"yandex" envconfig:"yandex"`
}

// Yandex contains yandex specific parameters.
type Yandex struct {
	CliendID     string `json:"client_id" envconfig:"client_id" required:"true" desc:"Yandex Cliend ID"`
	ClientSecret string `json:"client_secret" envconfig:"client_secret" required:"true" desc:"Yandex Cliend Secret"`
}

// Addr returns http address to listen.
func (s Specification) Addr() string {
	return fmt.Sprintf(":%d", s.Port)
}

// NewSpec loads environment variables
// and parse them into struct.
func NewSpec(prefix string) (Specification, error) {
	var spec Specification
	err := envconfig.Process(prefix, &spec)
	if err != nil {
		return spec, err
	}
	return spec, nil
}

// Usage prints usage information.
func Usage(prefix string, spec Specification) error {
	return envconfig.Usage(prefix, &spec)
}
