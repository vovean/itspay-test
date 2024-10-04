package config

import (
	"itspay/internal/config"
	"os"
	"testing"

	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func checkUnknownFields(configPath string, to any) error {
	f, err := os.Open(configPath)
	if err != nil {
		return err
	}

	decoder := yaml.NewDecoder(f)
	decoder.KnownFields(true)

	return decoder.Decode(to)
}

func TestConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		configPath string
		conf       v.Validatable
	}{
		{
			configPath: "rates_api.yml",
			conf:       &config.Config{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.configPath, func(t *testing.T) {
			t.Parallel()

			err := checkUnknownFields(tc.configPath, tc.conf)
			require.NoError(t, err)

			err = tc.conf.Validate()
			require.NoError(t, err)
		})
	}
}
