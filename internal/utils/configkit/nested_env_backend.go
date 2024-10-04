package configkit

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// NestedEnvBackend loads configuration from environment variables.
type NestedEnvBackend struct{}

// NewNestedEnvBackend creates a new environment variable backend which accounts nesting paths.
func NewNestedEnvBackend() *NestedEnvBackend {
	return &NestedEnvBackend{}
}

// Unmarshal loads values from environment variables into the provided struct.
func (b *NestedEnvBackend) Unmarshal(ctx context.Context, to any) error {
	val := reflect.ValueOf(to).Elem()
	t := val.Type()

	return b.processStruct("", val, t)
}

// processStruct processes each field in the struct and sets the value from environment variables.
func (b *NestedEnvBackend) processStruct(prefix string, val reflect.Value, t reflect.Type) error {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := val.Field(i)

		configTag := field.Tag.Get("config")
		if configTag == "" {
			configTag = strings.ToLower(field.Name)
		}

		envKey := b.buildEnvKey(prefix, configTag)

		switch fieldValue.Kind() { //nolint:exhaustive
		case reflect.Struct:
			err := b.processStruct(envKey, fieldValue, field.Type)
			if err != nil {
				return err
			}
		default:
			envValue, exists := os.LookupEnv(envKey)
			if exists {
				err := b.setFieldValue(fieldValue, envValue)
				if err != nil {
					return fmt.Errorf("failed to set field value for key \"%s\": %w", envKey, err)
				}
			}
		}
	}

	return nil
}

// buildEnvKey constructs the environment variable key from the prefix and config tag.
func (b *NestedEnvBackend) buildEnvKey(prefix, tag string) string {
	if prefix == "" {
		return strings.ToUpper(tag)
	}

	return strings.ToUpper(prefix + "__" + tag)
}

var (
	errInvalidDuration      = errors.New("invalid duration")
	errInvalidInt           = errors.New("invalid integer")
	errInvalidUint          = errors.New("invalid uint")
	errInvalidFloat         = errors.New("invalid float")
	errInvalidBool          = errors.New("invalid boolean")
	errUnsupportedFieldType = errors.New("unsupported field type")
)

// setFieldValue sets the value of the field based on its type.
func (b *NestedEnvBackend) setFieldValue(fieldValue reflect.Value, value string) error {
	switch fieldValue.Kind() { //nolint:exhaustive
	case reflect.String:
		fieldValue.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if fieldValue.Type() == reflect.TypeOf(time.Duration(0)) {
			duration, err := time.ParseDuration(value)
			if err != nil {
				return fmt.Errorf("%w: %w", errInvalidDuration, err)
			}

			fieldValue.Set(reflect.ValueOf(duration))
		} else {
			intValue, err := strconv.ParseInt(value, 10, fieldValue.Type().Bits())
			if err != nil {
				return fmt.Errorf("%w: %w", errInvalidInt, err)
			}

			fieldValue.SetInt(intValue)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintValue, err := strconv.ParseUint(value, 10, fieldValue.Type().Bits())
		if err != nil {
			return fmt.Errorf("%w: %w", errInvalidUint, err)
		}

		fieldValue.SetUint(uintValue)
	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, fieldValue.Type().Bits())
		if err != nil {
			return fmt.Errorf("%w: %w", errInvalidFloat, err)
		}

		fieldValue.SetFloat(floatValue)
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("%w: %w", errInvalidBool, err)
		}

		fieldValue.SetBool(boolValue)
	case reflect.Slice:
		elements := strings.Split(value, " ")
		sliceValue := reflect.MakeSlice(fieldValue.Type(), len(elements), len(elements))

		for i, elem := range elements {
			err := b.setFieldValue(sliceValue.Index(i), elem)
			if err != nil {
				return fmt.Errorf("failed to set slice element (index %d) value from value \"%s\": %w", i, value, err)
			}
		}

		fieldValue.Set(sliceValue)
	default:
		return fmt.Errorf("%w: %s", errUnsupportedFieldType, fieldValue.Type())
	}

	return nil
}

// Get is not implemented.
func (b *NestedEnvBackend) Get(ctx context.Context, _ string) ([]byte, error) {
	return nil, errors.New("not implemented")
}

// Name returns the name of the backend.
func (b *NestedEnvBackend) Name() string {
	return "nested-env"
}
