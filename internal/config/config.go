package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	HostAdr   string `env:"SERVER_ADDRESS"`
	BaseAdr   string `env:"BASE_URL"`
	LoggLevel string `env:"LOG_LEVEL"`
}

func InitConf() Config {
	var conf Config

	flag.StringVar(&conf.HostAdr, "a", "localhost:8080", "server host adress")
	flag.StringVar(&conf.BaseAdr, "b", "http://localhost:8080", "base return adress")
	flag.StringVar(&conf.LoggLevel, "l", "Info", "logging level")
	flag.Parse()

	err := env.Parse(&conf)
	if err != nil {
		fmt.Println(err)
	}
	return conf
}
