package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

func writeJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
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
