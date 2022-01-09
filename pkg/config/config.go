package config

import (
	"os"

	"github.com/joho/godotenv"
)

const (
	defaultAddr = ":4001"
)

type Config struct {
	Addr           string
	AllowedOrigins []string
	JWTSecret      string
	PostgresURL    string
}

func (c Config) setDefaults() Config {
	if c.Addr == "" {
		c.Addr = defaultAddr
	}

	return c
}

func Load() (Config, error) {
	// Optionally load env vars on .env file.
	_ = godotenv.Load()

	var c Config
	c.Addr = os.Getenv("ADDR")
	// Cleans input by removing spaces and split by comma.
	// c.AllowedOrigins = strings.Split(strings.ReplaceAll(os.Getenv("ALLOWED_ORIGINS"), " ", ""), ",")

	c.JWTSecret = os.Getenv("JWT_SECRET")
	c.PostgresURL = os.Getenv("POSTGRES_URL")
	c = c.setDefaults()
	return c, nil
}
