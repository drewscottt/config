package config_test

import (
	"context"
	"testing"
	"testing/synctest"

	"github.com/drewscottt/config"
)

func Test_HotReload_LoadGet_Succeeds(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		childCtx, cancelCtx := context.WithCancel(t.Context())

		h := &config.HotReload[int]{Context: childCtx}
		err := h.Load(config.EnvironmentVariable{
			Name:             "VAR_A",
			EnvironmentValue: "10",
			DefaultValue:     "5",
		})
		assertNoError(t, err)
		varA, err := h.Get()
		assertNoError(t, err)
		assertEqual(t, 10, varA)

		t.Setenv("VAR_A", "15")
		synctest.Sleep(config.HotReloadInterval * 2)
		varA, err = h.Get()
		assertNoError(t, err)
		assertEqual(t, 15, varA)

		t.Setenv("VAR_A", "")
		synctest.Sleep(config.HotReloadInterval * 2)
		varA, err = h.Get()
		assertNoError(t, err)
		assertEqual(t, 5, varA)

		t.Setenv("VAR_A", "not_an_int")
		synctest.Sleep(config.HotReloadInterval * 2)
		varA, err = h.Get()
		assertNoError(t, err)
		assertEqual(t, 5, varA)

		cancelCtx()
		synctest.Wait()
	})
}

func Test_HotReload_LoadGet_Fails_WhenNil(t *testing.T) {
	var h *config.HotReload[string]
	err := h.Load(config.EnvironmentVariable{EnvironmentValue: `"hello world"`})
	assertError(t, err)

	_, err = h.Get()
	assertError(t, err)
}

func Test_HotReload_Load_Fails_WhenInvalid(t *testing.T) {
	h := &config.HotReload[int]{}
	err := h.Load(config.EnvironmentVariable{EnvironmentValue: "not_an_int"})
	assertError(t, err)
}
