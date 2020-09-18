package entities

// ConfigGetter defines methods to fetch config values.
type ConfigGetter interface {
	GetInt(key string) int
	GetString(key string) string
}
