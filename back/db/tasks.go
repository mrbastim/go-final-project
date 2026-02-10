package db

import (
	"database/sql"
	"errors"
	"strings"
	"time"
)

func AddTask(db *sql.DB, task Task) (int32, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int32(id), nil

}

func GetTasks(db *sql.DB, search string, limit int) ([]Task, error) {
	if limit <= 0 {
		limit = 100
	}

	var rows *sql.Rows
	var err error

	search = strings.TrimSpace(search)

	if t, errDate := time.Parse("02.01.2006", search); errDate == nil {
		dateStr := t.Format("20060102")
		rows, err = db.Query(`SELECT 
		id, date, title, comment, repeat FROM scheduler 
		WHERE date = ? ORDER BY date LIMIT ?`,
			dateStr, limit)
	} else if search != "" {
		likePattern := "%" + search + "%"
		rows, err = db.Query(`SELECT 
		id, date, title, comment, repeat FROM scheduler 
		WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`,
			likePattern, likePattern, limit)
	} else {
		rows, err = db.Query(`SELECT 
		id, date, title, comment, repeat FROM scheduler 
		ORDER BY date LIMIT ?`, limit)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []Task{}
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func GetTaskByID(db *sql.DB, id string) (*Task, error) {
	var t Task
	err := db.QueryRow(`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`, id).
		Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func UpdateTask(db *sql.DB, task Task) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`,
		task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func DeleteTask(db *sql.DB, id string) error {
	_, err := db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	return err
}
