package config_test

import (
	"testing"

	"github.com/drewscottt/config"
)

func Test_JSON_LoadGet_Succeeds(t *testing.T) {
	j := &config.JSON[string]{}
	err := j.Load(config.ConfigVariable{ProviderValue: `"hello world"`}, config.LoadConfig{})
	assertNoError(t, err)

	val, err := j.Get()
	assertNoError(t, err)
	assertEqual(t, "hello world", val)
}

func Test_JSON_LoadGet_Fails_WhenNil(t *testing.T) {
	var j *config.JSON[string]
	err := j.Load(config.ConfigVariable{ProviderValue: `"hello world"`}, config.LoadConfig{})
	assertErrorIs(t, err, config.ErrLoaderNotInitialized)

	_, err = j.Get()
	assertErrorIs(t, err, config.ErrLoaderNotInitialized)
}
