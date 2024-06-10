package config

import (
	"os"

	"github.com/naoina/toml"
)

type Config struct {
	DB struct {
		Database string
		URL      string
	}

	Kafka struct {
		URL     string
		GroupId string
	}

	Info struct {
		Port string
	}
}

func NewConfig(path string) *Config {
	c := new(Config)

	if f, err := os.Open(path); err != nil {
		panic(err)
	} else if err = toml.NewDecoder(f).Decode(c); err != nil {
		panic(err)
	} else {
		return c
	}
}
