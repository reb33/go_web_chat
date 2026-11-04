package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	DSN string
	CommandTimeout int 
}

type Config struct {
	DB DatabaseConfig
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load(".env")

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		return nil, errors.New("dsn не установлен")
	}
	dbCommandTimeout := 10
	timeout := os.Getenv("DB_TIMEOUT")
	if timeout != "" {
		parsedTimeout, err := strconv.Atoi(timeout)
		if err == nil {
			dbCommandTimeout = parsedTimeout
		}
	}

	return &Config{
		DB: DatabaseConfig{
			DSN: dsn,
			CommandTimeout: dbCommandTimeout,
		},
	}, nil
}
