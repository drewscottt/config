package config

import (
	"encoding/json"
	"errors"
)

type JSON[T any] struct {
	val T
}

var _ Loader = (*JSON[string])(nil)

func (j *JSON[T]) Get() (T, error) {
	var zeroVal T
	if j == nil {
		return zeroVal, errors.New("value never set")
	}
	return j.val, nil
}

func (j *JSON[T]) Load(val string) error {
	if j == nil {
	    return errors.New("cannot load into nil")	
	}

    if val == "" {
        return nil
    }

	var parsedVal T
	err := json.Unmarshal([]byte(val), &parsedVal)
	if err != nil {
		return err
	}
	j.val = parsedVal
	return nil
}
