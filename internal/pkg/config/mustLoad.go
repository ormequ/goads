package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

func MustLoadENV[T any](path string) *T {
	var cfg T
	var err error
	if path != "" {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Fatalf("config file does not exist: %s", path)
		}
		err = cleanenv.ReadConfig(path, &cfg)
	} else {
		err = cleanenv.ReadEnv(&cfg)
	}
	if err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
