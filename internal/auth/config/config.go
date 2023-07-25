package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	Env          string `env:"ENV" env-default:"local"`
	PrivateKey   string `env:"AUTH_PRIVATE_KEY" env-required:"true"`
	PublicKey    string `env:"AUTH_PUBLIC_KEY" env-required:"true"`
	Expires      int    `env:"AUTH_EXPIRES_HOURS" env-default:"24"`
	PasswordCost int    `env:"PASSWORD_COST" env-default:"10"`
	HTTPAddress  string `env:"HTTP_ADDRESS" env-default:":8080"`
	GRPCAddress  string `env:"GRPC_ADDRESS" env-default:":8081"`
	DBHost       string `env:"POSTGRES_HOST" env-required:"true"`
	DBPort       uint16 `env:"POSTGRES_PORT" env-required:"true"`
	DBUser       string `env:"POSTGRES_USER" env-required:"true"`
	DBPassword   string `env:"POSTGRES_PASSWORD" env-required:"true"`
	DBName       string `env:"POSTGRES_DB" env-required:"true"`
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
