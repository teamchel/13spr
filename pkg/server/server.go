// pkg/server/server.go
package server

import (
	"log"
	"net/http"

	"13spr/pkg/api"
)

const webDir = "./web"

func Start(port string) error {
	api.Init()

	// Защита API-запросов
	http.HandleFunc("/api/task", api.TaskHandler)
	http.HandleFunc("/api/tasks", api.TasksHandler)
	http.HandleFunc("/api/task/done", api.DoneTaskHandler)
	http.HandleFunc("/api/task/delete", api.DeleteTaskHandler)

	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	log.Printf("Сервер запущен на порту %s\n", port)
	return http.ListenAndServe(":"+port, nil)
}
