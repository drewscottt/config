package config

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// ConfigProvider is an interface for a source for config variables.
type ConfigProvider interface {
	LoadConfigVariable(variableName string) (variableValue string, err error)
}

// EnvironmentProvider implements [ConfigProvider].
// It reads from the environment variables.
type EnvironmentProvider struct{}

// LoadConfigVariable implements [ConfigProvider].
func (e EnvironmentProvider) LoadConfigVariable(variableName string) (variableValue string, err error) {
	return os.Getenv(variableName), nil
}

var _ ConfigProvider = EnvironmentProvider{}

// FileProvider implements [ConfigProvider].
// It read froms from a source file, where each line is a key value pair, separated by `=`.
// For example:
// VAR_A=10
// VAR_B=string_value
type FileProvider struct {
	fileName string
}

// NewFileProvider creates a new [FileProvider].
func NewFileProvider(fileName string) FileProvider {
	return FileProvider{
		fileName: fileName,
	}
}

// LoadConfigVariable implements [ConfigProvider].
func (f FileProvider) LoadConfigVariable(variableName string) (variableValue string, err error) {
	file, err := os.Open(f.fileName)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	fileLines := strings.SplitSeq(string(fileContent), "\n")
	for fileLine := range fileLines {
		keyVal := strings.SplitN(fileLine, "=", 2)
		if keyVal[0] != variableName {
			continue
		}
		if len(keyVal) != 2 {
			return "", nil
		}
		return keyVal[1], nil
	}

	return "", nil
}

var _ ConfigProvider = FileProvider{}
