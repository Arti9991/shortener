package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Arti9991/shortener/internal/app/handlers"
	"github.com/Arti9991/shortener/internal/config"
	"github.com/Arti9991/shortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

func MainRouter(data *storage.Data, BaseAdr string) chi.Router {
	rt := chi.NewRouter()

	rt.Post("/", handlers.MainPage(data, BaseAdr))
	rt.Get("/{id}", handlers.GetAddr(data))
	rt.HandleFunc("/foo/{baz}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/bar/"+chi.URLParam(r, "baz"), http.StatusPermanentRedirect)
	})

	return rt
}
func main() {
	data := storage.NewData()
	cfg := config.InitConf()

	fmt.Printf("Host adr: %s\n", cfg.HostAdr)
	fmt.Printf("Base adr: %s\n", cfg.BaseAdr)

	log.Fatal(http.ListenAndServe(cfg.HostAdr, MainRouter(&data, cfg.BaseAdr)))
}
