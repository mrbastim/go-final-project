package api

import (
	"encoding/json"
	"main/back/db"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type TasksResponse struct {
	Tasks []db.Task `json:"tasks"`
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
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

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": "ID is required"})
		return
	}

	task, err := db.GetTaskByID(DB, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, map[string]string{"error": "Failed to retrieve task: " + err.Error()})
		return
	}

	if task == nil {
		w.WriteHeader(http.StatusNotFound)
		writeJson(w, map[string]string{"error": "Task not found"})
		return
	}

	writeJson(w, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.UseNumber()

	var payload map[string]interface{}
	if err := decoder.Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": "Invalid JSON: " + err.Error()})
		return
	}

	id, err := parseTaskID(payload["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	existingTask, err := db.GetTaskByID(DB, strconv.Itoa(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, map[string]string{"error": "Failed to check task: " + err.Error()})
		return
	}
	if existingTask == nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": "Task not found"})
		return
	}

	title, err := getStringField(payload, "title")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	if strings.TrimSpace(title) == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": "Title is required"})
		return
	}

	date, err := getStringField(payload, "date")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	comment, err := getStringField(payload, "comment")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	repeat, err := getStringField(payload, "repeat")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	task := db.Task{
		ID:      id,
		Date:    date,
		Title:   title,
		Comment: comment,
		Repeat:  repeat,
	}

	if err := checkDate(&task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	err = db.UpdateTask(DB, task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, map[string]string{"error": "Failed to update task: " + err.Error()})
		return
	}

	writeJson(w, map[string]string{})
}

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": "ID is required"})
		return
	}

	task, err := db.GetTaskByID(DB, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, map[string]string{"error": "Failed to retrieve task: " + err.Error()})
		return
	}
	if task == nil {
		w.WriteHeader(http.StatusNotFound)
		writeJson(w, map[string]string{"error": "Task not found"})
		return
	}

	if strings.TrimSpace(task.Repeat) == "" {
		err = db.DeleteTask(DB, id)
	} else {
		next, errNext := NextDate(time.Now(), task.Date, task.Repeat)
		if errNext != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeJson(w, map[string]string{"error": errNext.Error()})
			return
		}
		task.Date = next
		err = db.UpdateTask(DB, *task)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, map[string]string{"error": "Failed to update task: " + err.Error()})
		return
	}

	writeJson(w, map[string]string{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": "ID is required"})
		return
	}

	task, err := db.GetTaskByID(DB, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, map[string]string{"error": "Failed to retrieve task: " + err.Error()})
		return
	}
	if task == nil {
		w.WriteHeader(http.StatusNotFound)
		writeJson(w, map[string]string{"error": "Task not found"})
		return
	}

	err = db.DeleteTask(DB, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, map[string]string{"error": "Failed to delete task: " + err.Error()})
		return
	}

	writeJson(w, map[string]string{})
}
