package config

import (
	"context"
	"reflect"
	"time"
)

const (
	defaultHotReloadInterval = time.Second
)

// HotReload implements [Loader].
// It reloads the [ConfigVariable] periodically.
// Any type compatible with [Load] can be used here.
//
// Due to the way that operating systems manage processes and environment variables, [EnvironmentProvider] is not sufficient for hot reloading updates.
// Using another [ConfigProvider], like [FileProvider], is necessary.
type HotReload[T any] struct {
	val     T
	context context.Context // optional context for managing reloading lifecycle, mostly just for testing.
}

var _ Loader = (*HotReload[string])(nil)

// Load implements [Loader].
// It loads the variable, then starts a background goroutine to pick up changes.
func (h *HotReload[T]) Load(configVariable ConfigVariable, loadConfig LoadConfig) error {
	if h == nil {
		return ErrLoaderNotInitialized
	}

	parsedVar, err := configVariable.ParseIntoType(reflect.TypeFor[T](), loadConfig)
	if err != nil {
		return err
	}
	h.val = parsedVar.(T)

	go func() {
		hotReloadInterval := defaultHotReloadInterval
		if loadConfig.HotReloadInterval != 0 {
			hotReloadInterval = loadConfig.HotReloadInterval
		}
		ticker := time.NewTicker(hotReloadInterval)
		if h.context == nil {
			if loadConfig.HotReloadContext == nil {
				h.context = context.Background()
			} else {
				h.context = loadConfig.HotReloadContext
			}
		}

		for {
			select {
			case <-ticker.C:
				err := loadVariableValue(&configVariable, loadConfig)
				if err != nil {
					// TODO: our options are to either silently fail and keep the last valid value,
					// or return an error from Get()
					continue
				}

				parsedVar, err := configVariable.ParseIntoType(reflect.TypeFor[T](), loadConfig)
				if err != nil {
					// TODO: our options are to either silently fail and keep the last valid value,
					// or return an error from Get()
					continue
				}

				typedVar, ok := parsedVar.(T)
				if !ok {
					// TODO: our options are to either silently fail and keep the last valid value,
					// or return an error from Get()
					continue
				}

				// TODO: need to re-validate Config here?

				h.val = typedVar
			case <-h.context.Done():
				return
			}

		}
	}()

	return nil
}

func (h *HotReload[T]) Get() (T, error) {
	var zeroValue T
	if h == nil {
		return zeroValue, ErrLoaderNotInitialized
	}
	return h.val, nil
}
