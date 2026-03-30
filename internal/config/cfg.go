package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string `json:"databaseUrl"`
}

func New() *Config {
	return &Config{}
}

func (c *Config) Load() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	c.DatabaseURL = os.Getenv("DATABASE_URL")

	return nil
}
