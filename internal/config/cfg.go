package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string `json:"databaseUrl"`
	Port        string `json:"port"`
}

func New() *Config {
	return &Config{}
}

func (c *Config) Load() error {
	_ = godotenv.Load()

	c.DatabaseURL = os.Getenv("DATABASE_URL")
	c.Port = os.Getenv("PORT")

	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	if c.Port == "" {
		c.Port = "8080"
	}

	return nil
}
