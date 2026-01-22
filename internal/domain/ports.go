package domain

// ConfigProvider defines the interface for retrieving configuration.
type ConfigProvider interface {
	GetConfig() (*Config, error)
}
