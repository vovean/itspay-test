package configkit

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/heetch/confita"
	"github.com/stretchr/testify/require"
)

type NestedConfig struct {
	FieldString   string        `config:"field_string"`
	FieldInt      int           `config:"field_int"`
	FieldBool     bool          `config:"field_bool"`
	FieldFloat    float64       `config:"field_float"`
	FieldDuration time.Duration `config:"field_duration"`
	FieldStrings  []string      `config:"field_strings"`
	FieldInts     []int         `config:"field_ints"`
	FieldBools    []bool        `config:"field_bools"`
	FieldFloats   []float64     `config:"field_floats"`
}

type TestConfig struct {
	Nested        NestedConfig `config:"nested"`
	SimpleString  string       `config:"simple_string"`
	SimpleInt     int          `config:"simple_int"`
	SimpleUint    uint         `config:"simple_uint"`
	SimpleBool    bool         `config:"simple_bool"`
	SimpleFloat   float64      `config:"simple_float"`
	SimpleStrings []string     `config:"simple_strings"`
	SimpleInts    []int        `config:"simple_ints"`
	SimpleBools   []bool       `config:"simple_bools"`
	SimpleFloats  []float64    `config:"simple_floats"`
}

func setEnvVariables(t *testing.T, envVars map[string]string) {
	t.Helper()

	for key, value := range envVars {
		require.NoError(t, os.Setenv(key, value))
	}
}

func unsetEnvVariables(t *testing.T, envVars map[string]string) {
	t.Helper()

	for key := range envVars {
		require.NoError(t, os.Unsetenv(key))
	}
}

func TestNestedEnvBackend_Success(t *testing.T) { //nolint:paralleltest
	envVars := map[string]string{
		"NESTED__FIELD_STRING":   "nested string value",
		"NESTED__FIELD_INT":      "123",
		"NESTED__FIELD_BOOL":     "true",
		"NESTED__FIELD_FLOAT":    "45.67",
		"NESTED__FIELD_DURATION": "1h30m",
		"NESTED__FIELD_STRINGS":  "one two three",
		"NESTED__FIELD_INTS":     "1 2 3",
		"NESTED__FIELD_BOOLS":    "true false true",
		"NESTED__FIELD_FLOATS":   "1.1 2.2 3.3",
		"SIMPLE_STRING":          "simple string value",
		"SIMPLE_INT":             "456",
		"SIMPLE_UINT":            "789",
		"SIMPLE_BOOL":            "false",
		"SIMPLE_FLOAT":           "78.90",
		"SIMPLE_STRINGS":         "alpha beta gamma",
		"SIMPLE_INTS":            "10 20 30",
		"SIMPLE_BOOLS":           "false true false",
		"SIMPLE_FLOATS":          "4.4 5.5 6.6",
	}

	// Setting environment variables for the test
	setEnvVariables(t, envVars)
	defer unsetEnvVariables(t, envVars)

	// Creating a confita loader with the custom env backend
	loader := confita.NewLoader(NewNestedEnvBackend())

	// Creating an empty config struct to unmarshal into
	var cfg TestConfig

	// Loading the configuration using confita
	err := loader.Load(context.Background(), &cfg)
	require.NoError(t, err)

	// requiring the values
	require.Equal(t, "nested string value", cfg.Nested.FieldString)
	require.Equal(t, 123, cfg.Nested.FieldInt)
	require.True(t, cfg.Nested.FieldBool)
	require.InEpsilon(t, 45.67, cfg.Nested.FieldFloat, 0)
	require.Equal(t, 1*time.Hour+30*time.Minute, cfg.Nested.FieldDuration)
	require.Equal(t, []string{"one", "two", "three"}, cfg.Nested.FieldStrings)
	require.Equal(t, []int{1, 2, 3}, cfg.Nested.FieldInts)
	require.Equal(t, []bool{true, false, true}, cfg.Nested.FieldBools)
	require.InEpsilonSlice(t, []float64{1.1, 2.2, 3.3}, cfg.Nested.FieldFloats, 0)

	require.Equal(t, "simple string value", cfg.SimpleString)
	require.Equal(t, 456, cfg.SimpleInt)
	require.Equal(t, uint(789), cfg.SimpleUint)
	require.False(t, cfg.SimpleBool)
	require.InEpsilon(t, 78.90, cfg.SimpleFloat, 0)
	require.Equal(t, []string{"alpha", "beta", "gamma"}, cfg.SimpleStrings)
	require.Equal(t, []int{10, 20, 30}, cfg.SimpleInts)
	require.Equal(t, []bool{false, true, false}, cfg.SimpleBools)
	require.InEpsilonSlice(t, []float64{4.4, 5.5, 6.6}, cfg.SimpleFloats, 0)
}

func TestNestedEnvBackend_Error(t *testing.T) { //nolint:paralleltest
	// Creating a confita loader with the custom env backend
	loader := confita.NewLoader(NewNestedEnvBackend())

	// Creating an empty config struct to unmarshal into
	var cfg TestConfig

	tests := []struct {
		varName     string
		varValue    string
		expectedErr error
	}{
		{
			varName:     "NESTED__FIELD_INT",
			varValue:    "not_an_int",
			expectedErr: errInvalidInt,
		},
		{
			varName:     "NESTED__FIELD_BOOL",
			varValue:    "not_a_bool",
			expectedErr: errInvalidBool,
		},
		{
			varName:     "NESTED__FIELD_FLOAT",
			varValue:    "not_a_float",
			expectedErr: errInvalidFloat,
		},
		{
			varName:     "NESTED__FIELD_DURATION",
			varValue:    "not_a_duration",
			expectedErr: errInvalidDuration,
		},
		{
			varName:     "SIMPLE_UINT",
			varValue:    "not_a_uint",
			expectedErr: errInvalidUint,
		},
		{
			varName:     "NESTED__FIELD_INTS",
			varValue:    "1 2 three",
			expectedErr: errInvalidInt,
		},
		{
			varName:     "NESTED__FIELD_BOOLS",
			varValue:    "true false maybe",
			expectedErr: errInvalidBool,
		},
		{
			varName:     "NESTED__FIELD_FLOATS",
			varValue:    "1.1 2.2 three.point.three",
			expectedErr: errInvalidFloat,
		},
	}

	for _, tt := range tests { //nolint: paralleltest
		t.Run(tt.varName, func(t *testing.T) {
			setEnvVariables(t, map[string]string{tt.varName: tt.varValue})
			defer unsetEnvVariables(t, map[string]string{tt.varName: tt.varValue})

			err := loader.Load(context.Background(), &cfg)
			require.Error(t, err)
			require.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
