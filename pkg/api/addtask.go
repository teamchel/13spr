package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"13spr/pkg/db"
)

// writeJSON отправляет данные в формате JSON
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

// checkTask проверяет корректность задачи
func checkTask(task *db.Task) error {
	if task.Title == "" {
		return errors.New("не указан заголовок задачи")
	}

	now := time.Now().Format("20060102")

	// Проверяем дату
	var t time.Time
	var err error
	if task.Date == "" {
		task.Date = now
	} else {
		t, err = time.Parse("20060102", task.Date)
		if err != nil {
			return fmt.Errorf("некорректный формат даты: %w", err)
		}
		if afterNow(t, now) {
			// Если дата меньше текущей — обновляем на сегодня
			task.Date = now
		}
	}

	// Проверяем правило повторения
	if task.Repeat != "" {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("неподдерживаемое правило повторения: %w", err)
		}
		task.Date = next
	}

	return nil
}

// addTaskHandler обрабатывает POST /api/task
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка чтения запроса"})
		return
	}

	var task db.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка парсинга JSON"})
		return
	}

	err = checkTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка сохранения задачи"})
		return
	}

	writeJSON(w, map[string]string{"id": fmt.Sprintf("%d", id)})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	writeJSON(w, task)
}

func putTaskHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка чтения запроса"})
		return
	}

	var task db.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка парсинга JSON"})
		return
	}

	// Проверка обязательных полей
	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "не указан заголовок задачи"})
		return
	}

	// Проверка даты
	now := time.Now().Format("20060102")
	var t time.Time
	var errParse error
	if task.Date == "" {
		task.Date = now
	} else {
		t, errParse = time.Parse("20060102", task.Date)
		if errParse != nil {
			writeJSON(w, map[string]string{"error": "некорректный формат даты"})
			return
		}
		if afterNow(t, now) {
			task.Date = now
		}
	}

	// Проверка правила повторения
	if task.Repeat != "" {
		next, errNext := NextDate(now, task.Date, task.Repeat)
		if errNext != nil {
			writeJSON(w, map[string]string{"error": "неподдерживаемое правило повторения"})
			return
		}
		task.Date = next
	}

	// Обновляем задачу
	err = db.UpdateTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Успешное обновление — пустой ответ
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte("{}"))
}
