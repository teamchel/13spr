// pkg/api/tasks.go
package api

import (
	"13spr/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	const limit = 50 // максимальное количество задач

	tasks, err := db.Tasks(limit)
	if err != nil {
		// Используем writeJSON для правильной кодировки
		writeJSON(w, map[string]interface{}{"error": "ошибка получения задач"})
		return
	}

	// Отправляем ответ с задачами
	resp := TasksResp{Tasks: tasks}
	writeJSON(w, resp)
}
