package main

import (
	"log"

	"github.com/Arti9991/shortener/internal/server"
)

func main() {
	// запуск сервера со всеми настройками
	err := server.RunServer()
	if err != nil {
		log.Fatal(err)
	}
}
