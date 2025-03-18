package main

import (
	"log"
	"net/http"
	_ "net/http/pprof" // подключаем пакет pprof

	"github.com/Arti9991/shortener/internal/server"
)

func main() {

	router, addr, err := server.InitServer()
	if err != nil {
		log.Fatal(err)
	}

	err = http.ListenAndServe(addr, router)
	if err != nil {
		log.Fatal(err)
	}
}
