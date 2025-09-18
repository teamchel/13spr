package db

import (
	"database/sql"
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func GetTask(id string) (*Task, error) {
	query := `
		SELECT id, date, title, comment, repeat 
		FROM scheduler 
		WHERE id = ?
	`
	row := db.QueryRow(query, id)
	task := &Task{}
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("задача не найдена")
		}
		return nil, fmt.Errorf("ошибка чтения задачи: %w", err)
	}
	return task, nil
}

func UpdateTask(task *Task) error {
	query := `
		UPDATE scheduler 
		SET date = ?, title = ?, comment = ?, repeat = ? 
		WHERE id = ?
	`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("ошибка обновления: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка подсчёта строк: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("некорректный идентификатор для обновления")
	}

	return nil
}

func DeleteTask(id string) error {
	query := "DELETE FROM scheduler WHERE id = ?"
	res, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка подсчёта строк: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

func UpdateDate(next string, id string) error {
	query := "UPDATE scheduler SET date = ? WHERE id = ?"
	res, err := db.Exec(query, next, id)
	if err != nil {
		return fmt.Errorf("ошибка обновления даты: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка подсчёта строк: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("некорректный идентификатор для обновления")
	}

	return nil
}
