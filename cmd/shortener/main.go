package main

import (
	"net/http"

	"github.com/Arti9991/shortener/internal/app/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, handlers.MainPage)
	mux.HandleFunc(`/api/`, handlers.ApiPage)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
