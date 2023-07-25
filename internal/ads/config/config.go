package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	Env              string `env:"ENV" env-default:"local"`
	HTTPAddress      string `env:"HTTP_ADDRESS" env-default:"18080"`
	GRPCAddress      string `env:"GRPC_ADDRESS" env-default:"18081"`
	PostgresHost     string `env:"POSTGRES_HOST" env-required:"true"`
	PostgresPort     uint16 `env:"POSTGRES_PORT" env-required:"true"`
	PostgresUser     string `env:"POSTGRES_USER" env-required:"true"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" env-required:"true"`
	PostgresDB       string `env:"POSTGRES_DB" env-required:"true"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	var cfg Config
	var err error
	if configPath != "" {
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			log.Fatalf("config file does not exist: %s", configPath)
		}
		err = cleanenv.ReadConfig(configPath, &cfg)
	} else {
		err = cleanenv.ReadEnv(&cfg)
	}
	if err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
