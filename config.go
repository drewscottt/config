package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
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

		tagValue := os.Getenv(tag)

		switch fieldKind {
		case reflect.String:
			field.SetString(tagValue)
		case reflect.Int:
			if tagValue == "" {
				tagValue = "0"
			}
			intValue, err := strconv.Atoi(tagValue)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrInvalidValueForType, err)
			}
			field.SetInt(int64(intValue))
		case reflect.Float64:
			if tagValue == "" {
				tagValue = "0"
			}
			floatValue, err := strconv.ParseFloat(tagValue, 64)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrInvalidValueForType, err)
			}
			field.SetFloat(floatValue)
		case reflect.Bool:
			if tagValue == "" {
				tagValue = "false"
			}
			boolValue, err := strconv.ParseBool(tagValue)
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

			err := loader.Load(tagValue)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrInvalidValueForType, err)
			}
		default:
			return fmt.Errorf("%w: %v", ErrInvalidType, fieldKind)
		}
	}

	return nil
}
