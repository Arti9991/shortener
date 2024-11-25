package main

import (
	"log"
	"net/http"

	"github.com/Arti9991/shortener/internal/app/handlers"
	"github.com/Arti9991/shortener/internal/storage"
	"github.com/go-chi/chi/v5"
)

func MainRouter(data storage.Data) chi.Router {
	rt := chi.NewRouter()

	rt.Post("/", handlers.MainPage(&data))
	rt.Get("/{id}", handlers.GetAddr(&data))
	rt.HandleFunc("/foo/{baz}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/bar/"+chi.URLParam(r, "baz"), http.StatusPermanentRedirect)
	})

	return rt
}
func main() {
	data := storage.NewData()

	log.Fatal(http.ListenAndServe(":8080", MainRouter(data)))
}
