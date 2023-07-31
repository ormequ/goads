package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

func MustLoadENV[T any](path string) *T {
	var cfg T
	var err error
	if path != "" {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			panic(fmt.Sprintf("config file does not exist: %s", path))
		}
		err = cleanenv.ReadConfig(path, &cfg)
	} else {
		err = cleanenv.ReadEnv(&cfg)
	}
	if err != nil {
		panic(fmt.Sprintf("cannot read config: %s", err))
	}

	return &cfg
}
