package server

import (
	"net/http"
	"time"

	"github.com/Arti9991/shortener/internal/app/cmpgzip"
	"github.com/Arti9991/shortener/internal/app/handlers"
	"github.com/Arti9991/shortener/internal/config"
	"github.com/Arti9991/shortener/internal/logger"
	"github.com/Arti9991/shortener/internal/storage/database"
	"github.com/Arti9991/shortener/internal/storage/files"
	"github.com/Arti9991/shortener/internal/storage/inmemory"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"
)

type Server struct {
	Inmemory *inmemory.Data
	Config   config.Config
	Files    *files.FileData
	DataBase *database.DBStor
	hd       *handlers.HandlersData
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
	Serv.StorInit()

	return &Serv, nil
}

// создание роутера chi для хэндлеров
func (s *Server) MainRouter() chi.Router {

	rt := chi.NewRouter()
	rt.Post("/", logger.MiddlewareLogger(cmpgzip.MiddlewareGzip(handlers.PostAddr(s.hd))))
	rt.Get("/{id}", logger.MiddlewareLogger(cmpgzip.MiddlewareGzip(handlers.GetAddr(s.hd))))
	rt.Get("/ping", logger.MiddlewareLogger(cmpgzip.MiddlewareGzip(handlers.Ping(s.hd))))
	rt.Post("/api/shorten", logger.MiddlewareLogger(cmpgzip.MiddlewareGzip(handlers.PostAddrJSON(s.hd))))
	rt.Post("/api/shorten/batch", logger.MiddlewareLogger(cmpgzip.MiddlewareGzip(handlers.PostBatch(s.hd))))

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

	//чтение всех данных из файла в память
	err = serv.FileRead(serv.hd.Files)
	if err != nil {
		logger.Log.Info("Error in reading file!", zap.Error(err))
	}

	err = http.ListenAndServe(serv.Config.HostAdr, serv.MainRouter())
	return err
}
