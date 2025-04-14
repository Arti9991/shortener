// Здесь производится запуск и настройка сервера.
// В Example содержится пример работы с эндпоинтами.
package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "net/http/pprof" // подключаем пакет pprof

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"

	"github.com/Arti9991/shortener/internal/app/auth"
	"github.com/Arti9991/shortener/internal/app/cmpgzip"
	"github.com/Arti9991/shortener/internal/app/handlers"
	"github.com/Arti9991/shortener/internal/config"
	"github.com/Arti9991/shortener/internal/logger"
	"github.com/Arti9991/shortener/internal/storage/database"
	"github.com/Arti9991/shortener/internal/storage/files"
	"github.com/Arti9991/shortener/internal/storage/inmemory"
)

// Server хранит всю информацию для работы сервера.
type Server struct {
	Inmemory *inmemory.Data
	Files    *files.FileData
	DataBase *database.DBStor
	hd       *handlers.HandlersData
	Config   config.Config
}

// NewServer инциализирует все необходимые струткуры.
func NewServer(ctx context.Context) (*Server, error) {
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
	var wg sync.WaitGroup
	// инциализация хранилища с нужным интерфейсом
	Serv.StorInit(ctx, &wg)

	return &Serv, nil
}

// MainRouter создает роутер chi для хэндлеров.
func (s *Server) MainRouter() chi.Router {

	rt := chi.NewRouter()

	rt.Use(logger.MiddlewareLogger, cmpgzip.MiddlewareGzip, auth.MiddlewareAuth)

	rt.Mount("/debug", middleware.Profiler())

	rt.Post("/", handlers.PostAddr(s.hd))
	rt.Get("/{id}", handlers.GetAddr(s.hd))
	rt.Get("/ping", handlers.Ping(s.hd))
	rt.Post("/api/shorten", handlers.PostAddrJSON(s.hd))
	rt.Post("/api/shorten/batch", handlers.PostBatch(s.hd))
	rt.Get("/api/user/urls", handlers.GetAddrUser(s.hd))
	rt.Delete("/api/user/urls", handlers.DeleteAddr(s.hd))

	return rt
}

// RunServer запускает сервер со всеми полученными параметрами.
func RunServer() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serv, err := NewServer(ctx)
	if err != nil {
		return err
	}
	defer close(serv.hd.OutDelCh)

	srv := http.Server{
		Handler: serv.MainRouter(),
		Addr:    serv.Config.HostAdr,
	}

	logger.Log.Info("New server initialyzed!",
		zap.String("Server addres:", serv.Config.HostAdr),
		zap.String("Base addres:", serv.Config.BaseAdr),
	)

	// чтение всех данных из файла в память
	err = serv.FileRead(serv.hd.Files)
	if err != nil {
		logger.Log.Info("Error in reading file!", zap.Error(err))
	}

	RunWaitShutDown(serv.hd, &srv)

	// запуск горутины (описана в initStor.go).
	RunDeleteStor(serv.hd)

	// запуск сервера.
	if serv.Config.EnableHTTPS {
		err = srv.ListenAndServeTLS("server.crt", "server.key")
	} else {
		err = srv.ListenAndServe()
	}
	logger.Log.Info("Server shutted down!")
	return err
}
