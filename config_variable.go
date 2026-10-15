package config

import (
	"fmt"
	"reflect"
	"strconv"
)

// ConfigVariable represents a config variable.
type ConfigVariable struct {
	Name          string
	ProviderValue string
	DefaultValue  string
}

// Value returns the value of the [ConfigVariable]. Either the value from the provider, or the default.
func (c ConfigVariable) Value() string {
	if c.ProviderValue != "" {
		return c.ProviderValue
	}
	return c.DefaultValue
}

// ParseIntoType parses the raw string value of [ConfigVariable.Value] into the provided type.
//
// It is returned as any.
// TODO: This could be revisited as a generic function, but there's lots of use of runtime reflection which makes generics hard.
func (c ConfigVariable) ParseIntoType(intoType reflect.Type, loadConfig LoadConfig) (any, error) {
	value := c.Value()

	intoKind := intoType.Kind()

	switch intoKind {
	case reflect.String:
		return reflect.ValueOf(value).Convert(intoType).Interface(), nil
	case reflect.Int:
		if value == "" {
			value = "0"
		}
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInvalidValueForVariableType, err)
		}
		return reflect.ValueOf(intValue).Convert(intoType).Interface(), nil
	case reflect.Float64:
		if value == "" {
			value = "0"
		}
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInvalidValueForVariableType, err)
		}
		return reflect.ValueOf(floatValue).Convert(intoType).Interface(), nil
	case reflect.Bool:
		if value == "" {
			value = "false"
		}
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInvalidValueForVariableType, err)
		}
		return reflect.ValueOf(boolValue).Convert(intoType).Interface(), nil
	case reflect.Pointer:
		newValue := reflect.New(intoType.Elem())

		iFace := newValue.Interface()
		loader, ok := iFace.(Loader)
		if !ok {
			return nil, fmt.Errorf("%w: pointer must implement %v", ErrInvalidVariableType, reflect.TypeFor[Loader]())
		}

		err := loader.Load(c, loadConfig)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInvalidValueForVariableType, err)
		}
		return loader, nil
	default:
		return nil, fmt.Errorf("%w: %v", ErrInvalidVariableType, intoKind)
	}
}
