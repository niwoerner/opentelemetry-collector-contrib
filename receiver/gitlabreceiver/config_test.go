// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package gitlabreceiver

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/configopaque"
	"go.opentelemetry.io/collector/otelcol/otelcoltest"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/gitlabreceiver/internal/metadata"
)

func TestCreateDefaultConfig(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()

	assert.Equal(t,
		&Config{
			WebHook: WebHook{
				ServerConfig: confighttp.ServerConfig{
					Endpoint:     defaultEndpoint,
					ReadTimeout:  defaultReadTimeout,
					WriteTimeout: defaultWriteTimeout,
				},
				Path:       defaultPath,
				HealthPath: defaultHealthPath,
			},
		},
		cfg, "failed to create default config")

	assert.NoError(t, componenttest.CheckConfigStruct(cfg))
}

func TestLoadConfig(t *testing.T) {
	factories, err := otelcoltest.NopFactories()
	require.NoError(t, err)

	factory := NewFactory()
	factories.Receivers[metadata.Type] = factory

	cfg, err := otelcoltest.LoadConfigAndValidate(filepath.Join("testdata", "config.yaml"), factories)

	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Len(t, cfg.Receivers, 2)

	expectedConfig := &Config{
		WebHook: WebHook{
			ServerConfig: confighttp.ServerConfig{
				Endpoint:     "localhost:8080",
				ReadTimeout:  500 * time.Millisecond,
				WriteTimeout: 500 * time.Millisecond,
			},
			Path:       "some/path",
			HealthPath: "health/path",
			RequiredHeaders: map[string]configopaque.String{
				"key-present": "value-present",
			},
		},
	}

	r0 := cfg.Receivers[component.NewID(metadata.Type)]

	assert.Equal(t, expectedConfig, r0)

	expectedConfig.WebHook.RequiredHeaders = map[string]configopaque.String{
		"key-present":  "value-present",
		"key2-present": "value2-present",
	}

	r1 := cfg.Receivers[component.NewIDWithName(metadata.Type, "customname")].(*Config)

	assert.Equal(t, expectedConfig, r1)
}
