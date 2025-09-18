package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"13spr/pkg/auth"
	"13spr/pkg/db"
)

type LoginReq struct {
	Password string `json:"password"`
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка парсинга JSON"})
		return
	}

	// Получаем сохранённый пароль из переменной окружения
	expectedPass := os.Getenv("TODO_PASSWORD")
	if expectedPass == "" {
		writeJSON(w, map[string]string{"error": "сервер не настроен"})
		return
	}

	if req.Password != expectedPass {
		writeJSON(w, map[string]string{"error": "неверный пароль"})
		return
	}

	// Генерируем JWT-токен
	token, err := auth.GenerateToken(req.Password)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка генерации токена"})
		return
	}

	// Устанавливаем куку
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		MaxAge:   int(8 * 3600), // 8 часов
		SameSite: http.SameSiteLaxMode,
	})

	// Ответ с токеном
	writeJSON(w, map[string]string{"token": token})
}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	now := time.Now().Format("20060102")

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

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
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
