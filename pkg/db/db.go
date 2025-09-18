package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

const (
	dbFile = "scheduler.db"
	schema = `
CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date TEXT NOT NULL DEFAULT '',
	title TEXT NOT NULL DEFAULT '',
	comment TEXT DEFAULT '',
	repeat VARCHAR(128) DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
`
)

var db *sql.DB

func Init(dbPath string) error {
	var err error

	// Проверяем, существует ли файл базы данных
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Println("Файл базы данных не найден. Создаём новую.")
	} else if err != nil {
		return fmt.Errorf("ошибка при проверке файла: %w", err)
	}

	// Открываем базу данных
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("ошибка открытия базы данных: %w", err)
	}

	// Выполняем SQL-схему (создание таблицы и индекса)
	_, err = db.Exec(schema)
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы: %w", err)
	}

	log.Println("База данных успешно инициализирована")
	return nil
}

func AddTask(task *Task) (int64, error) {
	query := `
		INSERT INTO scheduler (date, title, comment, repeat) 
		VALUES (?, ?, ?, ?)
	`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
