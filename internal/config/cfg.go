package config

import (
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
	if err := godotenv.Load(); err != nil {
		return err
	}

	c.DatabaseURL = os.Getenv("DATABASE_URL")
	c.Port = os.Getenv("PORT")

	return nil
}
