package config

// Loader is an interface for custom/specialized environment variable types, like [JSON] and [HotReload].
type Loader interface {
	Load(ConfigVariable, LoadConfig) error
}
