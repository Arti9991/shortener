package main

import (
	"net/http"

	"github.com/Arti9991/shortener/internal/app/handlers"
	"github.com/Arti9991/shortener/internal/app/storage"
)

func main() {
	data := storage.NewData()
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, handlers.MainPage(&data))
	mux.HandleFunc(`/{id}`, handlers.GetAddr(&data))

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
