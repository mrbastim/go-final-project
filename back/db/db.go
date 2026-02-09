package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:date`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

const schemaSQL = `
CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date TEXT NOT NULL,
	title TEXT NOT NULL,
	comment TEXT NOT NULL DEFAULT '',
	repeat TEXT NOT NULL DEFAULT '' CHECK (length(repeat) <= 128)
);

CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

func OpenAndInit(dbFile string) (*sql.DB, error) {
	_, err := os.Stat(dbFile)
	install := err != nil

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	installed, err := schemaInstalled(db)
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	if install || !installed {
		if _, err := db.Exec(schemaSQL); err != nil {
			_ = db.Close()
			return nil, fmt.Errorf("install schema: %w", err)
		}
	}

	return db, nil
}

func schemaInstalled(db *sql.DB) (bool, error) {
	var name string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='scheduler'`).Scan(&name)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return name == "scheduler", nil
}

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
