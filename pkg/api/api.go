package api

import (
	"13spr/pkg/db"
	"net/http"
	"time"
)

func Init() {
	http.HandleFunc("/api/nextdate", NextDayHandler)
	http.HandleFunc("/api/task", TaskHandler)
	http.HandleFunc("/api/tasks", TasksHandler)
	http.HandleFunc("/api/task/done", DoneTaskHandler)
	http.HandleFunc("/api/task/delete", DeleteTaskHandler)
	http.HandleFunc("/api/signin", signinHandler)
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodPut:
		putTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func NextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse("20060102", nowStr)
		if err != nil {
			http.Error(w, "Invalid 'now' date", http.StatusBadRequest)
			return
		}
	}

	if dateStr == "" {
		http.Error(w, "Missing 'date'", http.StatusBadRequest)
		return
	}

	result, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(result))
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	const limit = 50 // максимальное количество задач

	tasks, err := db.Tasks(limit)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка получения задач"})
		return
	}

	// Если нет задач — отправляем пустой слайс
	writeJSON(w, TasksResp{Tasks: tasks})
}

func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	now := time.Now()

	if task.Repeat == "" {
		// Одноразовая задача — удаляем
		err = db.DeleteTask(id)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
	} else {
		// Периодическая задача — вычисляем следующую дату
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}

		err = db.UpdateDate(next, id)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
	}

	// Успешно — пустой ответ
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte("{}"))
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Успешное удаление — пустой ответ
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte("{}"))
}
