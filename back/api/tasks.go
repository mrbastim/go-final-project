package api

import (
	"main/back/db"
	"net/http"
)

type TasksResponse struct {
	Tasks []db.Task `json:"tasks"`
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {

	search := r.URL.Query().Get("search")

	tasks, err := db.GetTasks(DB, search, 50)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, map[string]string{"error": "Failed to retrieve tasks: " + err.Error()})
		return
	}

	if tasks == nil {
		tasks = []db.Task{}
	}

	writeJson(w, TasksResponse{Tasks: tasks})
}
