package api

import (
	"encoding/json"
	"errors"
	"main/back/db"
	"net/http"
	"strconv"
	"strings"
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

	// Проверяем, существует ли задача
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

func parseTaskID(value interface{}) (int, error) {
	if value == nil {
		return 0, errors.New("ID is required")
	}

	switch v := value.(type) {
	case string:
		v = strings.TrimSpace(v)
		if v == "" {
			return 0, errors.New("ID is required")
		}
		id, err := strconv.Atoi(v)
		if err != nil || id <= 0 {
			return 0, errors.New("Invalid ID")
		}
		return id, nil
	case json.Number:
		id64, err := v.Int64()
		if err != nil || id64 <= 0 {
			return 0, errors.New("Invalid ID")
		}
		return int(id64), nil
	case float64:
		id64 := int64(v)
		if v != float64(id64) || id64 <= 0 {
			return 0, errors.New("Invalid ID")
		}
		return int(id64), nil
	default:
		return 0, errors.New("Invalid ID type")
	}
}

func getStringField(payload map[string]interface{}, key string) (string, error) {
	value, ok := payload[key]
	if !ok || value == nil {
		return "", nil
	}
	str, ok := value.(string)
	if !ok {
		return "", errors.New("Invalid " + key)
	}
	return str, nil
}
