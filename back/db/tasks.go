package db

import (
	"database/sql"
	"errors"
	"strings"
	"time"
)

type Storage struct {
	db *sql.DB
}

type TaskStorage interface {
	AddTask(task Task) (int32, error)
	GetTasks(search string, limit int) ([]Task, error)
	GetTaskByID(id string) (*Task, error)
	UpdateTask(task Task) error
	DeleteTask(id string) error
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

const (
	queryLayout = "20060102"
	dateLayout  = "02.01.2006"
)

func (s *Storage) AddTask(task Task) (int32, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	result, err := s.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int32(id), nil

}

func (s *Storage) GetTasks(search string, limit int) ([]Task, error) {
	if limit <= 0 {
		limit = 100
	}

	var rows *sql.Rows
	var err error

	search = strings.TrimSpace(search)
	parsedDate, dateErr := time.Parse(dateLayout, search)

	switch {
	case search == "":
		rows, err = s.db.Query(`SELECT 
		id, date, title, comment, repeat FROM scheduler 
		ORDER BY date LIMIT ?`, limit)
	case dateErr == nil:
		dateStr := parsedDate.Format(queryLayout)
		rows, err = s.db.Query(`SELECT 
		id, date, title, comment, repeat FROM scheduler 
		WHERE date = ? ORDER BY date LIMIT ?`,
			dateStr, limit)
	default:
		likePattern := "%" + search + "%"
		rows, err = s.db.Query(`SELECT 
		id, date, title, comment, repeat FROM scheduler 
		WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`,
			likePattern, likePattern, limit)
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

func (s *Storage) GetTaskByID(id string) (*Task, error) {
	t := &Task{}
	err := s.db.QueryRow(`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`, id).
		Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (s *Storage) UpdateTask(task Task) error {
	tx, err := s.db.Begin()
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

func (s *Storage) DeleteTask(id string) error {
	_, err := s.db.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	return err
}
