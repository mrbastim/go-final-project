package api

import (
	"main/back/db"
	"net/http"
)

var storage db.TaskStorage

func Init(taskStorage db.TaskStorage) {
	storage = taskStorage
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/task/done", taskDoneHandler)
	http.HandleFunc("/api/tasks", getTasksHandler)
}
