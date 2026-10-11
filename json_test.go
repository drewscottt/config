package config_test

import (
	"testing"

	"github.com/drewscottt/config"
)

func Test_JSON_LoadGet_Succeeds(t *testing.T) {
	j := &config.JSON[string]{}
	err := j.Load(config.EnvironmentVariable{EnvironmentValue: `"hello world"`})
	assertNoError(t, err)

	val, err := j.Get()
	assertNoError(t, err)
	assertEqual(t, "hello world", val)
}

func Test_JSON_LoadGet_Fails_WhenNil(t *testing.T) {
	var j *config.JSON[string]
	err := j.Load(config.EnvironmentVariable{EnvironmentValue: `"hello world"`})
	assertError(t, err)

	_, err = j.Get()
	assertError(t, err)
}
