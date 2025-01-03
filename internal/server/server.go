package server

import (
	"net/http"
	"time"

	"github.com/Arti9991/shortener/internal/app/cmpgzip"
	"github.com/Arti9991/shortener/internal/app/handlers"
	"github.com/Arti9991/shortener/internal/config"
	"github.com/Arti9991/shortener/internal/files"
	"github.com/Arti9991/shortener/internal/logger"
	"github.com/Arti9991/shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"
)

type Server struct {
	Storage *storage.Data
	Config  config.Config
	Files   *files.FileData
}

// инциализация всех необходимых струткур
func NewServer() (*Server, error) {
	// установка сида для случайных чисел
	rand.Seed(uint64(time.Now().UnixNano()))
	var Serv Server
	// инициализация конфигурации
	Serv.Config = config.InitConf()
	// инициализация логгера
	err := logger.Initialize(Serv.Config.LoggLevel)
	if err != nil {
		return nil, err
	}
	// инциализация хранилища в памяти
	Serv.Storage = storage.NewData()
	// инциализация структуры для файлов
	Serv.Files, err = files.NewFiles(Serv.Config.FilePath, Serv.Storage)
	if err != nil {
		logger.Log.Info("Error in creating or file! Setting in memory mode!", zap.Error(err))
	}

	return &Serv, nil
}

// создание роутера chi для хэндлеров
func (s *Server) MainRouter() chi.Router {
	hd := handlers.NewHandlersData(s.Storage, s.Config.BaseAdr, s.Files)

	rt := chi.NewRouter()
	rt.Post("/", logger.MiddlewareLogger(cmpgzip.MiddlewareGzip(handlers.PostAddr(hd))))
	rt.Post("/api/shorten", logger.MiddlewareLogger(cmpgzip.MiddlewareGzip(handlers.PostAddrJSON(hd))))
	rt.Get("/{id}", logger.MiddlewareLogger(cmpgzip.MiddlewareGzip(handlers.GetAddr(hd))))

	return rt
}

// запуск сервера со всеми полученными параметрами
func RunServer() error {
	serv, err := NewServer()
	if err != nil {
		return err
	}

	logger.Log.Info("New server initialyzed!",
		zap.String("Server addres:", serv.Config.HostAdr),
		zap.String("Base addres:", serv.Config.BaseAdr),
	)

	err = serv.Files.FileRead()
	if err != nil {
		logger.Log.Info("Error in reading file!", zap.Error(err))
	}

	err = http.ListenAndServe(serv.Config.HostAdr, serv.MainRouter())
	return err
}
