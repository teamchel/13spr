package server

import (
	"log"
	"net/http"

	"13spr/pkg/api"
	"13spr/pkg/auth"
)

const webDir = "./web"

func Start(port string) error {
	api.Init()

	// Защита API-запросов
	http.HandleFunc("/api/task", auth.AuthMiddleware(api.taskHandler))
	http.HandleFunc("/api/tasks", auth.AuthMiddleware(api.tasksHandler))
	http.HandleFunc("/api/task/done", auth.AuthMiddleware(api.doneTaskHandler))
	http.HandleFunc("/api/task/delete", auth.AuthMiddleware(api.deleteTaskHandler))

	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	log.Printf("Сервер запущен на порту %s\n", port)
	return http.ListenAndServe(":"+port, nil)
}
