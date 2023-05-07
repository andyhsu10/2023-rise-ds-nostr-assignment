package configs

import "os"

const (
	prod = "production"
)

// Config object
type Config struct {
	Env        string `env:"ENV"`
	RelayUrl   string `env:"RELAY_URL"`
	PublicKey  string `env:"PUBLIC_KEY"`
	PrivateKey string `env:"PRIVATE_KEY"`
}

// IsProd Checks if env is production
func (c Config) IsProd() bool {
	return c.Env == prod
}

// GetConfig gets all config for the application
func GetConfig() Config {
	return Config{
		Env:        os.Getenv("ENV"),
		RelayUrl:   os.Getenv("RELAY_URL"),
		PublicKey:  os.Getenv("PUBLIC_KEY"),
		PrivateKey: os.Getenv("PRIVATE_KEY"),
	}
}
