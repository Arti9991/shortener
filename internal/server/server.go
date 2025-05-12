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
	Hd       *handlers.HandlersData
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
	// wait group для ожидания завершения горутин
	// у хэндлеров
	var wgWgHandler sync.WaitGroup
	// инциализация хранилища с нужным интерфейсом
	Serv.StorInit(ctx, &wgWgHandler, Serv.Config.TrustedNet)

	return &Serv, nil
}

// MainRouter создает роутер chi для хэндлеров.
func (s *Server) MainRouter() chi.Router {

	rt := chi.NewRouter()

	rt.Use(logger.MiddlewareLogger, cmpgzip.MiddlewareGzip, auth.MiddlewareAuth)

	rt.Mount("/debug", middleware.Profiler())

	rt.Post("/", handlers.PostAddr(s.Hd))
	rt.Get("/{id}", handlers.GetAddr(s.Hd))
	rt.Get("/ping", handlers.Ping(s.Hd))
	rt.Post("/api/shorten", handlers.PostAddrJSON(s.Hd))
	rt.Post("/api/shorten/batch", handlers.PostBatch(s.Hd))
	rt.Get("/api/user/urls", handlers.GetAddrUser(s.Hd))
	rt.Delete("/api/user/urls", handlers.DeleteAddr(s.Hd))
	rt.Get("/api/internal/stats", handlers.GetStats(s.Hd))

	return rt
}

// RunServer запускает сервер со всеми полученными параметрами.
func RunRestServer() error {
	// канал для сообщения о Shutdown
	shutCh := make(chan struct{})
	// Wait Group для ожидания завершения горутины удаления
	var WgStor sync.WaitGroup
	// контекст для ожидания системного сигнала на завершение работы
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serv, err := NewServer(ctx)
	if err != nil {
		return err
	}

	srv := http.Server{
		Handler: serv.MainRouter(),
		Addr:    serv.Config.HostAdr,
	}

	logger.Log.Info("New server initialyzed!",
		zap.String("Server addres:", serv.Config.HostAdr),
		zap.String("Base addres:", serv.Config.BaseAdr),
	)

	// чтение всех данных из файла в память
	err = serv.FileRead(serv.Hd.Files)
	if err != nil {
		logger.Log.Info("Error in reading file!", zap.Error(err))
	}

	WgStor.Add(1)
	// запуск горутины (описана в initStor.go).
	RunDeleteStor(serv.Hd, &WgStor)

	// запуск функции ожидающей сигнала о завершении
	// (описана в initStor.go).
	RunWaitShutDown(serv.Hd, &srv, shutCh)

	// запуск сервера.
	if serv.Config.EnableHTTPS {
		err = srv.ListenAndServeTLS("server.crt", "server.key")
		if err != nil && err != http.ErrServerClosed {
			logger.Log.Info("Error in ListenAndServeTLS", zap.Error(err))
			return err
		}
	} else {
		err = srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Log.Info("Error in ListenAndServe", zap.Error(err))
			return err
		}
	}
	// ожидание сообщения о Shutdown
	<-shutCh
	// ожидания закрытия горутины у хэндлера
	serv.Hd.Wg.Wait()
	// закрытие канала отправки URL под удаление
	close(serv.Hd.OutDelCh)
	// ожидание остановки горутины с функцией удаления
	WgStor.Wait()
	// закртытие соединения с базой данных
	err = serv.Hd.Dt.CloseDB()
	if err != nil {
		logger.Log.Info("Error in database close", zap.Error(err))
	}
	logger.Log.Info("Server shutted down!")
	return nil
}
