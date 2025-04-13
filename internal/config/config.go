package config

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/Arti9991/shortener/internal/logger"
	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
)

// Config структура со всемии конфигурируемыми параметрами.
type Config struct {
	HostAdr     string `env:"SERVER_ADDRESS" json:"server_address"`        // адрес сервера
	BaseAdr     string `env:"BASE_URL"  json:"base_url"`                   // базовый адрес возвращаемого URL
	LoggLevel   string `env:"LOG_LEVEL"  json:"logg_level"`                // уровень логгирования
	FilePath    string `env:"FILE_STORAGE_PATH"  json:"file_storage_path"` // путь к файлу хранения
	DBAddress   string `env:"DATABASE_DSN"  json:"database_dsn"`           // данные для подключения к базе
	EnableHTTPS bool   `env:"ENABLE_HTTPS"  json:"enable_https"`           // флаг работы через HTTPS или через HTTP
	ConfigAddr  string `env:"CONFIG" json:"config"`                        // флаг для конфигурации из файла
}

// InitConf инициализация конфигурации для чтения флагов и переменных окружения.
func InitConf() Config {
	var conf Config

	flag.StringVar(&conf.HostAdr, "a", "localhost:8080", "server host adress")
	flag.StringVar(&conf.BaseAdr, "b", "http://localhost:8080", "base return adress")
	flag.StringVar(&conf.LoggLevel, "l", "Info", "logging level")
	flag.StringVar(&conf.FilePath, "f", "", "storage file path")
	flag.StringVar(&conf.DBAddress, "d", "", "database address") //"host=localhost user=myuser password=123456 dbname=ShortURL sslmode=disable"
	flag.BoolVar(&conf.EnableHTTPS, "s", false, "Secure or not")
	flag.StringVar(&conf.ConfigAddr, "c", conf.ConfigAddr, "JSON config file")
	flag.Parse()

	err := env.Parse(&conf)
	if err != nil {
		fmt.Println(err)
	}

	JSONConfig := ReadConfig(conf.ConfigAddr)
	resConf := SaveConfig(&JSONConfig, &conf)

	return resConf
}

// InitConfTests инициализация тестовой конфигурации с заданными параметрами.
func InitConfTests() Config {
	var conf Config
	conf.HostAdr = "localhost:8080"
	conf.BaseAdr = "http://example.com"
	conf.LoggLevel = "Info"
	conf.FilePath = ""
	conf.DBAddress = ""
	conf.EnableHTTPS = false
	return conf
}

// ReadConfig функция чтения конфигурации из файла
func ReadConfig(cfgFilePath string) Config {
	var config Config
	file, err := os.OpenFile(cfgFilePath, os.O_RDONLY, 0644)
	if err != nil {
		logger.Log.Error("Wrong config file!", zap.Error(err))
		return config
	}
	defer file.Close()
	buff, err := io.ReadAll(file)
	if err != nil {
		logger.Log.Error("Bad read config file!", zap.Error(err))
		return config
	}
	err = json.Unmarshal(buff, &config)
	if err != nil {
		logger.Log.Error("Bad unmarshall config file!", zap.Error(err))
		return config
	}
	return config
}

// SaveConfig сохранение конфигурации в нужном приоритете
// (переменные окружения - флаги - файл конфигурации)
func SaveConfig(JSONConfig *Config, BaseConfig *Config) Config {

	if BaseConfig.HostAdr == "" {
		BaseConfig.HostAdr = JSONConfig.HostAdr
	}
	if BaseConfig.BaseAdr == "" {
		BaseConfig.BaseAdr = JSONConfig.BaseAdr
	}
	if BaseConfig.LoggLevel == "" {
		BaseConfig.LoggLevel = JSONConfig.LoggLevel
	}
	if BaseConfig.FilePath == "" {
		BaseConfig.FilePath = JSONConfig.FilePath
	}
	if BaseConfig.DBAddress == "" {
		BaseConfig.DBAddress = JSONConfig.DBAddress
	}
	if !BaseConfig.EnableHTTPS {
		BaseConfig.EnableHTTPS = JSONConfig.EnableHTTPS
	}
	return *BaseConfig
}

// CreateJSON создание файла конфигурации со стоковыми параметрами в папке проекта
func CreateJSON(cfgFilePath string) {
	var config Config

	config.HostAdr = "localhost:8080"
	config.BaseAdr = "http://localhost:8080"
	config.LoggLevel = "Info"
	config.FilePath = ""
	config.DBAddress = "host=localhost user=myuser password=123456 dbname=ShortURL sslmode=disable"
	config.EnableHTTPS = false
	config.ConfigAddr = cfgFilePath

	file, err := os.OpenFile(cfgFilePath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		logger.Log.Error("Wrong config file!", zap.Error(err))
	}
	defer file.Close()

	bt, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		logger.Log.Error("Bad marshall config file!", zap.Error(err))
	}
	buf := bytes.NewBuffer(bt)
	_, err = buf.WriteTo(file)
	if err != nil {
		logger.Log.Error("Bad write config file!", zap.Error(err))
	}
	logger.Log.Info("Config file created!", zap.Error(err))
}
