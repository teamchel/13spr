package api

import (
	"net/http"

	"13spr/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	const limit = 50 // максимальное количество задач

	tasks, err := db.Tasks(limit)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка получения задач"})
		return
	}

	// Если нет задач — отправляем пустой слайс
	writeJSON(w, TasksResp{Tasks: tasks})
}
