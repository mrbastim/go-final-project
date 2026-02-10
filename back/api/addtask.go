package api

import (
	"encoding/json"
	"errors"
	"main/back/db"
	"net/http"
	"strings"
	"time"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": "Invalid JSON: " + err.Error()})
		return
	}

	if strings.TrimSpace(task.Title) == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": "Title is required"})
		return
	}

	if err := checkDate(&task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(DB, task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, map[string]string{"error": "Failed to add task: " + err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	writeJson(w, map[string]int32{"id": id})
}

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(dateLayout)
	}

	t, err := time.ParseInLocation(dateLayout, task.Date, time.Local)
	if err != nil {
		return errors.New("invalid date format")
	}

	var next string

	if len(strings.TrimSpace(task.Repeat)) > 0 {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return errors.New("invalid repeat format")
		}
	}

	if afterNow(now, t) {
		if len(strings.TrimSpace(task.Repeat)) == 0 {
			task.Date = now.Format(dateLayout)
		} else {
			task.Date = next
		}
	}

	return nil
}
