package config

import (
	"encoding/json"
)

// JSON implements [Loader].
// It deserializes an [ConfigVariable] as JSON as a value with the type of the generic type parameter `T`.
type JSON[T any] struct {
	val T
}

var _ Loader = (*JSON[string])(nil)

// Load implements [Loader].
// It loads and deserializes the value as JSON.
func (j *JSON[T]) Load(envVar ConfigVariable, _ LoadConfig) error {
	if j == nil {
		return ErrLoaderNotInitialized
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

// Get returns the value.
func (j *JSON[T]) Get() (T, error) {
	var zeroVal T
	if j == nil {
		return zeroVal, ErrLoaderNotInitialized
	}
	return j.val, nil
}
