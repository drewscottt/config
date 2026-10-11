package config

import (
	"encoding/json"
	"errors"
)

// JSON implements [Loader].
// It deserializes an [EnvironmentVariable] as JSON as a value with the type of the generic type parameter `T`.
type JSON[T any] struct {
	val T
}

var _ Loader = (*JSON[string])(nil)

// Get returns the value.
func (j *JSON[T]) Get() (T, error) {
	var zeroVal T
	if j == nil {
		return zeroVal, errors.New("value never set")
	}
	return j.val, nil
}

// Load loads and deserializes the value.
func (j *JSON[T]) Load(envVar EnvironmentVariable) error {
	if j == nil {
		return errors.New("cannot load into nil")
	}

	if envVar.Value() == "" {
		// Empty string is not valid JSON, but we'll treat it is the zero value rather than fail deserialization.
		return nil
	}

	var parsedVal T
	err := json.Unmarshal([]byte(envVar.Value()), &parsedVal)
	if err != nil {
		return err
	}
	j.val = parsedVal
	return nil
}
