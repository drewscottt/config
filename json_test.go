package config_test

import (
    "testing"

    "github.com/drewscottt/config"
)

func Test_JSON_Load(t *testing.T) {
    var j *config.JSON[string]
    err := j.Load(`"hello world"`)
    assertError(t, err)

    _, err = j.Get()
    assertError(t, err)
}
