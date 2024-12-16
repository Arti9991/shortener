package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Arti9991/shortener/internal/app/gzipcomp"
	"github.com/Arti9991/shortener/internal/app/handlers"
	"github.com/Arti9991/shortener/internal/config"
	"github.com/Arti9991/shortener/internal/logger"
	"github.com/Arti9991/shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"
)

type Server struct {
	Storage *storage.Data
	Config  config.Config
}

func NewServer() *Server {
	rand.Seed(uint64(time.Now().UnixNano()))
	var Serv Server
	stor := storage.NewData()
	Serv.Storage = &stor
	Serv.Config = config.InitConf()
	return &Serv
}

func (s *Server) MainRouter() chi.Router {
	rt := chi.NewRouter()
	rt.Post("/", logger.MiddlewareLogger(gzipcomp.MiddlewareGzip(handlers.PostAddr(s.Storage, s.Config.BaseAdr))))
	rt.Post("/api/shorten", logger.MiddlewareLogger(gzipcomp.MiddlewareGzip(handlers.PostAddrJSON(s.Storage, s.Config.BaseAdr))))
	rt.Get("/{id}", logger.MiddlewareLogger(gzipcomp.MiddlewareGzip(handlers.GetAddr(s.Storage))))

	return rt
}

func RunServer() error {
	serv := NewServer()

	if err := logger.Initialize(serv.Config.LoggLevel); err != nil {
		return err
	}
	logger.Log.Info("New server initialyzed!", zap.String("Server addres:", serv.Config.HostAdr))
	fmt.Printf("Host adr: %s\n", serv.Config.HostAdr)
	fmt.Printf("Base adr: %s\n", serv.Config.BaseAdr)

	err := http.ListenAndServe(serv.Config.HostAdr, serv.MainRouter())
	return err
}
