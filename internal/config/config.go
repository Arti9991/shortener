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
	FilePath  string `env:"FILE_STORAGE_PATH"`
	DBAddress string `env:"DATABASE_DSN"`
}

// инициализация конфигурации для чтения флагов и переменных окружения
func InitConf() Config {
	var conf Config

	flag.StringVar(&conf.HostAdr, "a", "localhost:8080", "server host adress")
	flag.StringVar(&conf.BaseAdr, "b", "http://localhost:8080", "base return adress")
	flag.StringVar(&conf.LoggLevel, "l", "Info", "logging level")
	flag.StringVar(&conf.FilePath, "f", "", "storage file path")
	flag.StringVar(&conf.DBAddress, "d", "", "database address") //"host=localhost user=myuser password=123456 dbname=ShortURL sslmode=disable"
	flag.Parse()

	err := env.Parse(&conf)
	if err != nil {
		fmt.Println(err)
	}
	return conf
}

// инициализация тестовой конфигурации с заданными параметрами
func InitConfTests() Config {
	var conf Config
	conf.HostAdr = "localhost:8080"
	conf.BaseAdr = "http://example.com"
	conf.LoggLevel = "Info"
	conf.FilePath = ""
	conf.DBAddress = ""
	return conf
}
