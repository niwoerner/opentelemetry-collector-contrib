package gitlabreceiver

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
)

// Config that is exposed to this gitlab receiver through the OTEL config.yaml
type Config struct {
	WebHook WebHook `mapstructure:"webhook"`
}

type WebHook struct {
	confighttp.ServerConfig `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct
	Path                    string                   `mapstructure:"path"`            // path for data collection. Default is /events
	HealthPath              string                   `mapstructure:"health_path"`     // path for health check api. Default is /health_check
	RequiredHeader          RequiredHeader           `mapstructure:"required_header"` // optional setting to set a required header for all requests to have
	Secret                  string                   `mapstructure:"secret"`          // secret for webhook
}

type RequiredHeader struct {
	Key   string `mapstructure:"key"`
	Value string `mapstructure:"value"`
}

func createDefaultConfig() component.Config {
	return &Config{
		WebHook: WebHook{
			ServerConfig: confighttp.ServerConfig{
				Endpoint:     defaultEndpoint,
				ReadTimeout:  defaultReadTimeout,
				WriteTimeout: defaultWriteTimeout,
			},
			Path:       defaultPath,
			HealthPath: defaultHealthPath,
		},
	}
}
