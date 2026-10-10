package config_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/drewscottt/config"
)

var _ config.Config = TestConfigAllTypes{}

type TestConfigAllTypes struct {
	VarA string                   `config:"VAR_A"`
	VarB int                      `config:"VAR_B"`
	VarC float64                  `config:"VAR_C"`
	VarD *config.JSON[JSONConfig] `config:"VAR_D"`
	VarE bool                     `config:"VAR_E"`
}

func (t TestConfigAllTypes) Validate() error {
	return nil
}

type JSONConfig struct {
	Field1 string `json:"field_1"`
}

func Test_Load_Succeeds(t *testing.T) {
	t.Setenv("VAR_A", "val1")
	t.Setenv("VAR_B", "10")
	t.Setenv("VAR_C", "10.5")
	t.Setenv("VAR_D", `{"field_1": "value1"}`)
	t.Setenv("VAR_E", "true")
	c := TestConfigAllTypes{}
	err := config.Load(&c)
	assertNoError(t, err)

	assertEqual(t, "val1", c.VarA)
	assertEqual(t, 10, c.VarB)
	assertEqual(t, 10.5, c.VarC)
	varD, err := c.VarD.Get()
	assertNoError(t, err)
	assertEqual(t, JSONConfig{Field1: "value1"}, varD)
	assertEqual(t, true, c.VarE)
}

func Test_Load_Succeeds_ZeroValues(t *testing.T) {
	c := TestConfigAllTypes{}
	err := config.Load(&c)
	assertNoError(t, err)

	assertEqual(t, "", c.VarA)
	assertEqual(t, 0, c.VarB)
	assertEqual(t, 0, c.VarC)
	varD, err := c.VarD.Get()
	assertNoError(t, err)
	assertEqual(t, JSONConfig{}, varD)
	assertEqual(t, false, c.VarE)
}

func Test_Load_Succeeds_NegativeNumericValues(t *testing.T) {
	t.Setenv("VAR_B", "-10")
	t.Setenv("VAR_C", "-10.5")
	c := TestConfigAllTypes{}
	err := config.Load(&c)
	assertNoError(t, err)

	assertEqual(t, -10, c.VarB)
	assertEqual(t, -10.5, c.VarC)
}

type TestConfigInvalidTypes struct {
	VarA time.Time `config:"VAR_A"`
}

func (t TestConfigInvalidTypes) Validate() error {
	return nil
}

func Test_Load_Fails_InvalidTypes(t *testing.T) {
	c := TestConfigInvalidTypes{}
	err := config.Load(&c)
	assertErrorIs(t, err, config.ErrInvalidType)
}

type TestConfigInvalidPointer struct {
	VarA *JSONConfig `config:"VAR_A"`
}

func (t TestConfigInvalidPointer) Validate() error {
	return nil
}

func Test_Load_Fails_InvalidPointer(t *testing.T) {
	c := TestConfigInvalidPointer{}
	err := config.Load(&c)
	assertErrorIs(t, err, config.ErrInvalidType)
}

type TestConfigValidateError struct {
	VarA string `config:"VAR_A"`
}

var errValidation = errors.New("an error")

func (t TestConfigValidateError) Validate() error {
	return errValidation
}

func Test_Load_Fails_ValidationError(t *testing.T) {
	c := TestConfigValidateError{}
	err := config.Load(&c)
	assertErrorIs(t, err, config.ErrValidation)
	assertErrorIs(t, err, errValidation)
}

type TestConfigInt struct {
	VarA int `config:"VAR_A"`
}

func (t TestConfigInt) Validate() error {
	return nil
}

type TestConfigBool struct {
	VarA bool `config:"VAR_A"`
}

func (t TestConfigBool) Validate() error {
	return nil
}

type TestConfigJSON struct {
	VarA *config.JSON[JSONConfig] `config:"VAR_A"`
}

func (t TestConfigJSON) Validate() error {
	return nil
}

type TestConfigFloat64 struct {
	VarA float64 `config:"VAR_A"`
}

func (t TestConfigFloat64) Validate() error {
	return nil
}

func Test_Load_Fails_InvalidValuesForTypes(t *testing.T) {
	t.Run("invalid int", func(t *testing.T) {
		t.Setenv("VAR_A", "doesnt_parse_as_int")
		c := TestConfigInt{}
		err := config.Load(&c)
		assertErrorIs(t, err, config.ErrInvalidValueForType)
	})

	t.Run("invalid bool", func(t *testing.T) {
		t.Setenv("VAR_A", "doesnt_parse_as_bool")
		c := TestConfigBool{}
		err := config.Load(&c)
		assertErrorIs(t, err, config.ErrInvalidValueForType)
	})

	t.Run("invalid json", func(t *testing.T) {
		t.Setenv("VAR_A", `{"invalid_json":}`)
		c := TestConfigJSON{}
		err := config.Load(&c)
		assertErrorIs(t, err, config.ErrInvalidValueForType)
	})

	t.Run("invalid float64", func(t *testing.T) {
		t.Setenv("VAR_A", "doesnt_parse_as_float64")
		c := TestConfigFloat64{}
		err := config.Load(&c)
		assertErrorIs(t, err, config.ErrInvalidValueForType)
	})
}

type TestConfigIntPointer struct {
	VarA *int `config:"VAR_A"`
}

func (t TestConfigIntPointer) Validate() error {
	return nil
}

func Test_Load_Fails_InvalidPointerType(t *testing.T) {
	c := TestConfigIntPointer{}
	err := config.Load(&c)
	assertError(t, err)
}

type TestConfigNoConfigTags struct {
	VarA int
}

func (t TestConfigNoConfigTags) Validate() error {
	return nil
}

func Test_Load_Succeeds_NoConfigStructTags(t *testing.T) {
	c := TestConfigNoConfigTags{}
	err := config.Load(&c)
	assertNoError(t, err)
}

type TestConfigPrivateField struct {
	varA int `config:"VAR_A"`
}

func (t TestConfigPrivateField) Validate() error {
	return nil
}

func Test_Load_Succeeds_PrivateField(t *testing.T) {
	c := TestConfigPrivateField{}
	err := config.Load(&c)
	assertErrorIs(t, err, config.ErrCannotPutConfigTagOnPrivateField)
}

// Testing helpers... Thank you, Go standard library, for doing the absolute bare minimum. It is so "simple" to have to repeat basic stuff in every library ever written.

func assertEqual[T comparable](t *testing.T, expected, actual T, msg ...string) {
	if expected == actual {
		return
	}

	t.Helper()
	t.Errorf("expected %v, got %v; %s", expected, actual, strings.Join(msg, ""))
}

func assertError(t *testing.T, err error, msg ...string) {
	if err != nil {
		return
	}

	t.Helper()
	t.Errorf("expected error, but got nil; %s", strings.Join(msg, ""))
}

func assertNoError(t *testing.T, err error, msg ...string) {
	if err == nil {
		return
	}

	t.Helper()
	t.Errorf("received unexpected error: %v; %s", err, strings.Join(msg, ""))
}

func assertErrorIs(t *testing.T, err, targetErr error, msg ...string) {
	if errors.Is(err, targetErr) {
		return
	}

	t.Helper()
	t.Errorf("received unexpected error: [%v]; expected [%v]; %s", err, targetErr, strings.Join(msg, ""))
}
