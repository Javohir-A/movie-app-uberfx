package config

import (
	"os"

	"sync"

	"github.com/joho/godotenv"
	"github.com/spf13/cast"
)

var (
	conf *Config
	once sync.Once
)

// atomic
func Configs() *Config {
	once.Do(func() {
		conf = load()
	})
	return conf
}

func load() *Config {
	c := Config{}
	err := godotenv.Load()
	if err != nil {
		return nil
	}

	c.DBHost = cast.ToString(getOrReturnDefault("DB_HOST", "localhost"))
	c.DBName = cast.ToString(getOrReturnDefault("DB_NAME", "name"))
	c.DBPassword = cast.ToString(getOrReturnDefault("DB_PASSWORD", "password"))
	c.DBPort = cast.ToString(getOrReturnDefault("DB_PORT", "5432"))
	c.DBUser = cast.ToString(getOrReturnDefault("DB_USER", "postgres"))

	c.Port = cast.ToString(getOrReturnDefault("PORT", "7777"))
	c.JWTSecret = cast.ToString(getOrReturnDefault("JWT_SECRET", "2343rfe"))
	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "info"))

	return &c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	if os.Getenv(key) == "" {
		return defaultValue
	}
	return os.Getenv(key)
}

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	JWTSecret  string
	Port       string
	LogLevel   string
}
