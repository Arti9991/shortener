package server

import (
	"fmt"
	"net/http"

	"github.com/Arti9991/shortener/internal/app/handlers"
	"github.com/Arti9991/shortener/internal/config"
	"github.com/Arti9991/shortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	Storage *storage.Data
	Config  config.Config
}

func NewServer() *Server {
	var Serv Server
	stor := storage.NewData()
	Serv.Storage = &stor
	Serv.Config = config.InitConf()
	return &Serv
}

func (s *Server) MainRouter() chi.Router {
	rt := chi.NewRouter()

	rt.Post("/", handlers.MainPage(s.Storage, s.Config.BaseAdr))
	rt.Get("/{id}", handlers.GetAddr(s.Storage))

	return rt
}

func RunServer() error {
	serv := NewServer()

	fmt.Printf("Host adr: %s\n", serv.Config.HostAdr)
	fmt.Printf("Base adr: %s\n", serv.Config.BaseAdr)

	err := http.ListenAndServe(serv.Config.HostAdr, serv.MainRouter())
	return err
}
