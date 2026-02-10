package api

import (
	"database/sql"
	"net/http"
)

var DB *sql.DB

func Init(db *sql.DB) {
	DB = db
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/task/done", taskDoneHandler)
	http.HandleFunc("/api/tasks", getTasksHandler)
}
