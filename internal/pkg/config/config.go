package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type DbConfig struct {
	Name     string `env:"DB_NAME"     env-default:"mydb"     env-description:"Database name"`
	User     string `env:"DB_USER"     env-default:"postgres" env-description:"Database user"`
	Password string `env:"DB_PASSWORD" env-default:"password" env-description:"Database password"`
	Host     string `env:"DB_HOST"     env-default:"localhost" env-description:"Database host"`
	Port     string `env:"DB_PORT"     env-default:"5432"     env-description:"Database port"`
	SSLMode  string `env:"DB_SSLMODE"  env-default:"disable"  env-description:"Database SSL mode"`
}

func (db DbConfig) URL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.Name,
		db.SSLMode,
	)
}

type Config struct {
	DB         DbConfig `env-prefix:""`
	ServerPort string   `env:"SERVER_PORT" env-default:"8080" env-description:"HTTP server port"`
}

func FromEnv() (Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return Config{}, fmt.Errorf("read env: %w", err)
	}
	return cfg, nil
}
