package api

import (
	"main/back/db"
	"net/http"
)

func Init(storage db.TaskStorage) {
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/task/done", taskDoneHandler)
	http.HandleFunc("/api/tasks", getTasksHandler)
}
