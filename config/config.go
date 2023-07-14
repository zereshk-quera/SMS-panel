package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config represents the application configuration.
	Config struct {
		App           App           `yaml:"app"`
		HTTP          HTTP          `yaml:"http"`
		PG            PG            `yaml:"pg"`
		AdminRegister AdminRegister `yaml:"adminregister"`
	}

	// App represents the application-specific configuration.
	App struct {
		Name    string `yaml:"name" env:"APP_NAME" env-required:"true"`
		Version string `yaml:"version" env:"APP_VERSION" env-required:"true"`
	}

	// HTTP represents the HTTP server configuration.
	HTTP struct {
		Port string `yaml:"port" env:"HTTP_PORT" env-required:"true"`
	}

	// PG represents the PostgreSQL configuration.
	PG struct {
		HOST     string `env:"POSTGRES_HOST"`
		PORT     string `env:"POSTGRES_PORT"`
		DB       string `env:"POSTGRES_DB"`
		USER     string `env:"POSTGRES_USER"`
		PASSWORD string `env:"POSTGRES_PASSWORD"`
		SSLMODE  string `env:"POSTGRES_SSLMODE"`
		TIMEZONE string `env:"POSTGRES_TIMEZONE"`
	}

	// AdminRegister represents the configuration for the admin registration.
	AdminRegister struct {
		ADMIN_CODE string `yaml:"admin_code" env:"ADMIN_CODE"`
	}
)

// NewConfig returns the application configuration based on the provided configuration files and environment variables.
func NewConfig() (*Config, error) {
	var cfg Config

	// Read the configuration from the YAML file
	if err := cleanenv.ReadConfig("config.yaml", &cfg); err != nil {
		log.Printf("Error reading configuration: %v", err)
		return nil, err
	}

	// Read the environment variables
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Printf("Error reading environment variables: %v", err)
		return nil, err
	}

	return &cfg, nil
}
