package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

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
