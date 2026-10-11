package config

import (
	"context"
	"errors"
	"os"
	"reflect"
	"time"
)

const (
	HotReloadInterval = time.Second
)

// HotReload implements [Loader].
// It reloads the environment variable periodically.
// Any type compatible with [Load] can be used here.
type HotReload[T any] struct {
	val     T
	Context context.Context // optional context for managing reloading lifecycle, mostly just for testing.
}

var _ Loader = (*HotReload[string])(nil)

func (h *HotReload[T]) Load(envVar EnvironmentVariable) error {
	if h == nil {
		return errors.New("hot reload is nil")
	}

	parsedVar, err := parseEnvVariable(envVar, reflect.TypeFor[T]())
	if err != nil {
		return err
	}
	h.val = parsedVar.(T)

	go func() {
		ticker := time.NewTicker(HotReloadInterval)
		if h.Context == nil {
			h.Context = context.Background()
		}

		for {
			select {
			case <-ticker.C:
				updatedEnvVarValue := os.Getenv(envVar.Name)
				envVar.EnvironmentValue = updatedEnvVarValue

				parsedVar, err := parseEnvVariable(envVar, reflect.TypeFor[T]())
				if err != nil {
					// TODO: our options are to either silently fail, and keep the last valid value,
					// or, return an error from Get()
					continue
				}
				h.val = parsedVar.(T)
			case <-h.Context.Done():
				return
			}

		}
	}()

	return nil
}

func (h *HotReload[T]) Get() (T, error) {
	var zeroValue T
	if h == nil {
		return zeroValue, errors.New("hot reload is nil")
	}
	return h.val, nil
}
