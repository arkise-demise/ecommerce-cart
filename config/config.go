package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUsername string
	DBPassword string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load("../.env")
	if err != nil {
		return nil, fmt.Errorf("failed to load .env file: %v", err)
	}

	config := &Config{
		DBHost:     os.Getenv("DBHost"),
		DBPort:     os.Getenv("DBPort"),
		DBName:     os.Getenv("DBName"),
		DBUsername: os.Getenv("DBUsername"),
		DBPassword: os.Getenv("DBPassword"),
	}

	if config.DBHost == "" || config.DBPort == "" || config.DBName == "" || config.DBUsername == "" || config.DBPassword == "" {
		return nil, fmt.Errorf("one or more environment variables are missing in the .env file")
	}

	return config, nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", c.DBUsername, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}
