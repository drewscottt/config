package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Config interface {
	Validate() error
}

type Loader interface {
	Load(string) error
}

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
	val := reflect.ValueOf(c).Elem()
	typ := val.Type()

	for i := range val.NumField() {
		field := val.Field(i)
		fieldType := typ.Field(i)
		fieldKind := field.Kind()

		tag := fieldType.Tag.Get("config")
		if tag == "" {
			continue
		}

		if !field.CanInterface() {
			return ErrCannotPutConfigTagOnPrivateField
		}

		parts := strings.Split(tag, ",")

		envVarValue := os.Getenv(parts[0])
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

		if envVarValue == "" {
			envVarValue = defaultValue
		}

		switch fieldKind {
		case reflect.String:
			field.SetString(envVarValue)
		case reflect.Int:
			if envVarValue == "" {
				envVarValue = "0"
			}
			intValue, err := strconv.Atoi(envVarValue)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrInvalidValueForType, err)
			}
			field.SetInt(int64(intValue))
		case reflect.Float64:
			if envVarValue == "" {
				envVarValue = "0"
			}
			floatValue, err := strconv.ParseFloat(envVarValue, 64)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrInvalidValueForType, err)
			}
			field.SetFloat(floatValue)
		case reflect.Bool:
			if envVarValue == "" {
				envVarValue = "false"
			}
			boolValue, err := strconv.ParseBool(envVarValue)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrInvalidValueForType, err)
			}
			field.SetBool(boolValue)
		case reflect.Pointer:
			fieldTypeType := fieldType.Type
			if field.IsNil() {
				newValue := reflect.New(fieldTypeType.Elem())
				field.Set(newValue)
			}

			iFace := field.Interface()
			loader, ok := iFace.(Loader)
			if !ok {
				return fmt.Errorf("%w: pointer must implement %v", ErrInvalidType, reflect.TypeFor[Loader]())
			}

			err := loader.Load(envVarValue)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrInvalidValueForType, err)
			}
		default:
			return fmt.Errorf("%w: %v", ErrInvalidType, fieldKind)
		}
	}

	return nil
}
