package main

import (
	"fmt"
	"log"

	"github.com/Arti9991/shortener/internal/server"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	// запуск сервера со всеми настройками
	err := server.RunServer()
	if err != nil {
		log.Fatal(err)
	}
}
