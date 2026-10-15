package config

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Config is an interface to provide runtime configuration based on environment.
type Config interface {
	Validate() error
}

// Load loads the [Config] from the config provider(s) and according to the struct tag definition.
//
// Config must be a pointer to a struct.
//
// Each field with a `config` tag in the implementation struct will be loaded with a value according to the
// corresponding environment variable, environment provider, or default value.
//
// The `config` tag is structured like: `config:"VAR_NAME,default=<optional default value>"`
// The fields with `config` tags must have types with kind: string, int, float64, bool, or be a pointer to [Loader].
func Load(c Config, opts ...LoadOption) error {
	lc := LoadConfig{
		ConfigProviders: []ConfigProvider{},
	}
	for _, opt := range opts {
		opt(&lc)
	}
	lc.ConfigProviders = append(lc.ConfigProviders, EnvironmentProvider{})

	err := loadConfigVariables(c, lc)
	if err != nil {
		return err
	}

	err = c.Validate()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrValidation, err)
	}

	return nil
}

// LoadConfig is the configuration for [Load] and for [Loader.Load].
type LoadConfig struct {
	ConfigProviders   []ConfigProvider
	HotReloadInterval time.Duration
	HotReloadContext  context.Context
}

// LoadOption is a functional option for customizing [Load].
type LoadOption func(*LoadConfig)

// AddConfigProvider adds a new [ConfigProvider] to [Load].
// By default, [EnvironmentProvider] is used.
// Order matters; the first added [ConfigProvider] takes first priority, then the [EnvironmentProvider] takes last priority.
func AddConfigProvider(cp ConfigProvider) LoadOption {
	return func(lc *LoadConfig) {
		if lc.ConfigProviders == nil {
			lc.ConfigProviders = make([]ConfigProvider, 0, 1)
		}
		lc.ConfigProviders = append(lc.ConfigProviders, cp)
	}
}

// SetHotReloadInterval sets the interval for hot reloading.
func SetHotReloadInterval(d time.Duration) LoadOption {
	return func(lc *LoadConfig) {
		lc.HotReloadInterval = d
	}
}

// SetHotReloadInterval sets the context for managing the hot reloading lifecycle.
func SetHotReloadContext(c context.Context) LoadOption {
	return func(lc *LoadConfig) {
		lc.HotReloadContext = c
	}
}

func loadConfigVariables(c Config, loadConfig LoadConfig) error {
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

		configVariable := configStructTagToConfigVariable(tag)
		err := loadVariableValue(&configVariable, loadConfig)
		if err != nil {
			return err // TODO
		}

		parsedConfigVariable, err := configVariable.ParseIntoType(structFieldValue.Type(), loadConfig)
		if err != nil {
			return err
		}
		structFieldValue.Set(reflect.ValueOf(parsedConfigVariable))
	}

	return nil
}

func configStructTagToConfigVariable(tagValue string) ConfigVariable {
	parts := strings.Split(tagValue, ",")

	configVarName := parts[0]
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

	return ConfigVariable{
		Name:         configVarName,
		DefaultValue: defaultValue,
	}
}

func loadVariableValue(configVariable *ConfigVariable, lc LoadConfig) error {
	value := ""
	defer func() {
		configVariable.ProviderValue = value
	}()

	for _, cp := range lc.ConfigProviders {
		providerValue, err := cp.LoadConfigVariable(configVariable.Name)
		if err != nil {
			return err
		}
		if providerValue != "" {
			value = providerValue
			return nil
		}
	}
	value = configVariable.DefaultValue
	return nil
}
