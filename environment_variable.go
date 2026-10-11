package config

// EnvironmentVariable represents an Environment Variable.
type EnvironmentVariable struct {
	Name             string
	EnvironmentValue string
	DefaultValue     string
}

// Value returns the value of the [EnvironmentVariable]. Either the value from the environment, or the default.
func (e EnvironmentVariable) Value() string {
	if e.EnvironmentValue != "" {
		return e.EnvironmentValue
	}
	return e.DefaultValue
}
