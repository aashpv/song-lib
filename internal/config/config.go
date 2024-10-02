package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"time"
)

type Config struct {
	Database `yaml:"database"`
	LogLevel string `env:"LOG_LEVEL"`
	Server   `yaml:"server"`
}

type Database struct {
	Host     string `env:"DB_HOST"`
	Port     int    `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Name     string `env:"DB_NAME"`
}

type Server struct {
	Host        string        `env:"SERVER_HOST"`
	Port        int           `env:"SERVER_PORT"`
	Timeout     time.Duration `env:"SERVER_TIMEOUT"`
	IdleTimeout time.Duration `env:"SERVER_IDLE_TIMEOUT"`
}

// MustLoad загружает конфигурацию из файла и переменных окружения
func MustLoad() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("Error reading environment variables: %v", err)
	}

	return &cfg
}
