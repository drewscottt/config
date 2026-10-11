package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Config is an interface to provide runtime configuration based on environment.
type Config interface {
	Validate() error
}

// Load loads the [Config] according to the environment and struct tags.
//
// Config must be a pointer to a struct.
//
// Each field with a `config` tag in the implementation struct will be loaded with a value according to the
// corresponding environment variable, or default value.
// The `config` tag is structured like: `config:"ENV_VAR_NAME,default=<optional default value>"`
// The fields with `config` tags must have types with kind: string, int, float64, bool, or be a pointer to [Loader].
func Load(c Config) error {
	err := parseConfigStructTags(c)
	if err != nil {
		return err
	}

	err = c.Validate()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrValidation, err)
	}
	return nil
}

func parseConfigStructTags(c Config) error {
	if reflect.TypeOf(c).Kind() != reflect.Pointer {
		return ErrConfigMustBePointerToStruct
	}

	configStruct := reflect.ValueOf(c).Elem()
	configStructType := configStruct.Type()

	if configStructType.Kind() != reflect.Struct {
		return ErrConfigMustBePointerToStruct
	}

	for i := range configStruct.NumField() {
		structField := configStructType.Field(i)

		tag := structField.Tag.Get("config")
		if tag == "" {
			continue
		}

		structFieldValue := configStruct.Field(i)

		if !structFieldValue.CanInterface() {
			return ErrCannotPutConfigTagOnPrivateField
		}

		envVar := configStructTagToEnvVar(tag)

		parsedEnvVar, err := parseEnvVariable(envVar, structFieldValue.Type())
		if err != nil {
			return err
		}
		structFieldValue.Set(reflect.ValueOf(parsedEnvVar))
	}

	return nil
}

func parseEnvVariable(envVar EnvironmentVariable, intoType reflect.Type) (any, error) {
	envVarValue := envVar.Value()

	intoKind := intoType.Kind()

	switch intoKind {
	case reflect.String:
		return envVarValue, nil
	case reflect.Int:
		if envVarValue == "" {
			envVarValue = "0"
		}
		intValue, err := strconv.Atoi(envVarValue)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInvalidValueForType, err)
		}
		return intValue, nil
	case reflect.Float64:
		if envVarValue == "" {
			envVarValue = "0"
		}
		floatValue, err := strconv.ParseFloat(envVarValue, 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInvalidValueForType, err)
		}
		return floatValue, nil
	case reflect.Bool:
		if envVarValue == "" {
			envVarValue = "false"
		}
		boolValue, err := strconv.ParseBool(envVarValue)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInvalidValueForType, err)
		}
		return boolValue, nil
	case reflect.Pointer:
		newValue := reflect.New(intoType.Elem())

		iFace := newValue.Interface()
		loader, ok := iFace.(Loader)
		if !ok {
			return nil, fmt.Errorf("%w: pointer must implement %v", ErrInvalidType, reflect.TypeFor[Loader]())
		}

		err := loader.Load(envVar)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInvalidValueForType, err)
		}
		return loader, nil
	default:
		return nil, fmt.Errorf("%w: %v", ErrInvalidType, intoKind)
	}
}

func configStructTagToEnvVar(tagValue string) EnvironmentVariable {
	parts := strings.Split(tagValue, ",")

	envVarName := parts[0]
	defaultValue := ""
	for i, part := range parts {
		if i == 0 {
			continue
		}

		partParts := strings.SplitN(part, "=", 2)
		switch partParts[0] {
		case "default":
			defaultValue = partParts[1]
		default:
			continue
		}
	}

	return EnvironmentVariable{
		Name:             envVarName,
		EnvironmentValue: os.Getenv(envVarName),
		DefaultValue:     defaultValue,
	}
}
