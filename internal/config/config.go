package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	Env         string `env:"ENV" env-default:"local"`
	HTTPAddress string `env:"HTTP_ADDRESS" env-default:"18080"`
	GRPCAddress string `env:"GRPC_ADDRESS" env-default:"18081"`
	DBHost      string `env:"DB_HOST" env-required:"true"`
	DBPort      uint16 `env:"DB_PORT" env-required:"true"`
	DBUser      string `env:"DB_USER" env-required:"true"`
	DBPassword  string `env:"DB_PASSWORD" env-required:"true"`
	DBName      string `env:"DB_NAME" env-required:"true"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
