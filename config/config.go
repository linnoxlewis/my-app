package config

type Config struct{}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) GetAddress() string {
	return ":8080"
}
