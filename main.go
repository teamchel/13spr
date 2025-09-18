package main

import (
	"log"
	"os"

	"13spr/pkg/server"
)

func main() {
	port := getPort()
	log.Printf("Запуск сервера на порту %s\n", port)

	if err := server.Start(port); err != nil {
		log.Fatal(err)
	}
}

func getPort() string {
	if port := os.Getenv("TODO_PORT"); port != "" {
		return port
	}
	return "7540"
}
