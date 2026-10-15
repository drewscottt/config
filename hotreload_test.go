package config_test

import (
	"context"
	"testing"
	"testing/synctest"
	"time"

	"github.com/drewscottt/config"
)

func Test_HotReload_LoadGet_Succeeds(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		childCtx, cancelCtx := context.WithCancel(t.Context())

		hotReloadInterval := time.Second
		h := &config.HotReload[int]{}
		err := h.Load(
			config.ConfigVariable{
				Name:          "VAR_A",
				ProviderValue: "10",
				DefaultValue:  "5",
			},
			config.LoadConfig{
				ConfigProviders:   []config.ConfigProvider{config.EnvironmentProvider{}},
				HotReloadInterval: hotReloadInterval,
				HotReloadContext:  childCtx,
			},
		)
		assertNoError(t, err)
		varA, err := h.Get()
		assertNoError(t, err)
		assertEqual(t, 10, varA)

		t.Setenv("VAR_A", "15")
		synctest.Sleep(hotReloadInterval * 2)
		varA, err = h.Get()
		assertNoError(t, err)
		assertEqual(t, 15, varA)

		t.Setenv("VAR_A", "")
		synctest.Sleep(hotReloadInterval * 2)
		varA, err = h.Get()
		assertNoError(t, err)
		assertEqual(t, 5, varA)

		t.Setenv("VAR_A", "not_an_int")
		synctest.Sleep(hotReloadInterval * 2)
		varA, err = h.Get()
		assertNoError(t, err)
		assertEqual(t, 5, varA)

		cancelCtx()
		synctest.Wait()
	})
}

func Test_HotReload_LoadGet_Fails_WhenNil(t *testing.T) {
	var h *config.HotReload[string]
	err := h.Load(config.ConfigVariable{ProviderValue: `"hello world"`}, config.LoadConfig{})
	assertErrorIs(t, err, config.ErrLoaderNotInitialized)

	_, err = h.Get()
	assertErrorIs(t, err, config.ErrLoaderNotInitialized)
}

func Test_HotReload_Load_Fails_WhenInvalid(t *testing.T) {
	h := &config.HotReload[int]{}
	err := h.Load(config.ConfigVariable{ProviderValue: "not_an_int"}, config.LoadConfig{})
	assertError(t, err)
}
