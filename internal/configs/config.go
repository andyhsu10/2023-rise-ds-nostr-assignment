package configs

import "os"

const (
	prod = "production"
)

// Config object
type Config struct {
	Env       string `env:"ENV"`
	PublicKey string `env:"PUBLIC_KEY"`
	SecretKey string `env:"SECRET_KEY"`
}

// IsProd Checks if env is production
func (c Config) IsProd() bool {
	return c.Env == prod
}

// GetConfig gets all config for the application
func GetConfig() Config {
	return Config{
		Env:       os.Getenv("ENV"),
		PublicKey: os.Getenv("PUBLIC_KEY"),
		SecretKey: os.Getenv("SECRET_KEY"),
	}
}
