package main

import (
	"fmt"
	"log"

	protoServer "github.com/Arti9991/shortener/internal/gRPC/server"
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

	RestServer := true
	// запуск сервера со всеми настройками
	if RestServer {
		err := server.RunRestServer()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := protoServer.RunGRPCServer()
		if err != nil {
			log.Fatal(err)
		}
	}
}
