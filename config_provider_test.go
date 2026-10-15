package config_test

import (
	_ "embed"
	"testing"

	"github.com/drewscottt/config"
)

func Test_FileProvider_BasicFile(t *testing.T) {
	fileProvider := config.NewFileProvider("test_data/file_provider/test1.env")
	varA, err := fileProvider.LoadConfigVariable("VAR_A")
	assertNoError(t, err)
	assertEqual(t, "VALUE", varA)
	varB, err := fileProvider.LoadConfigVariable("VAR_B")
	assertNoError(t, err)
	assertEqual(t, "VALUE2", varB)
}

func Test_FileProvider_FileNoExist(t *testing.T) {
	fileProvider := config.NewFileProvider("test_data/file_provider/doesnt_exist.env")
	_, err := fileProvider.LoadConfigVariable("VAR_A")
	assertError(t, err)
}

func Test_FileProvider_FileMalformed(t *testing.T) {
	fileProvider := config.NewFileProvider("test_data/file_provider/test2.env")
	varA, err := fileProvider.LoadConfigVariable("VAR_A")
	assertNoError(t, err)
	assertEqual(t, "", varA)
}

func Test_FileProvider_VariableMissing(t *testing.T) {
	fileProvider := config.NewFileProvider("test_data/file_provider/test2.env")
	varA, err := fileProvider.LoadConfigVariable("VAR_B")
	assertNoError(t, err)
	assertEqual(t, "", varA)
}
