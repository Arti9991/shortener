package config

import (
	"flag"
)

type Config struct {
	HostAdr string
	BaseAdr string
}

func InitConf() Config {
	var conf Config

	flag.StringVar(&conf.HostAdr, "a", "localhost:8080", "server host adress")
	flag.StringVar(&conf.BaseAdr, "b", "localhost:8080", "base return adress")
	flag.Parse()

	return conf
}
