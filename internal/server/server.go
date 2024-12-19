package server

import (
	"fmt"
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

func NewServer() *Server {
	rand.Seed(uint64(time.Now().UnixNano()))
	var Serv Server
	Serv.Config = config.InitConf()

	stor := storage.NewData()
	Serv.Storage = stor

	fl := files.NewFiles(Serv.Config.FilePath, Serv.Storage)
	Serv.Files = fl

	return &Serv
}

func (s *Server) MainRouter() chi.Router {
	hd := handlers.NewHandlersData(s.Storage, s.Config.BaseAdr, s.Files)

	rt := chi.NewRouter()
	rt.Post("/", logger.MiddlewareLogger(cmpgzip.MiddlewareGzip(handlers.PostAddr(hd))))
	rt.Post("/api/shorten", logger.MiddlewareLogger(cmpgzip.MiddlewareGzip(handlers.PostAddrJSON(hd))))
	rt.Get("/{id}", logger.MiddlewareLogger(cmpgzip.MiddlewareGzip(handlers.GetAddr(hd))))

	return rt
}

func RunServer() error {
	serv := NewServer()
	//defer serv.Files.FileSave(serv.Config.FilePath)

	if err := logger.Initialize(serv.Config.LoggLevel); err != nil {
		return err
	}
	logger.Log.Info("New server initialyzed!", zap.String("Server addres:", serv.Config.HostAdr))
	fmt.Printf("Host adr: %s\n", serv.Config.HostAdr)
	fmt.Printf("Base adr: %s\n", serv.Config.BaseAdr)
	serv.Files.FileRead()

	err := http.ListenAndServe(serv.Config.HostAdr, serv.MainRouter())
	return err
}
