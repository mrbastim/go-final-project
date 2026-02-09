package server

import (
	"database/sql"
	"fmt"
	"main/back/api"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

const defaultPort = 7540

var DB *sql.DB

func Run(db *sql.DB) error {
	DB = db

	api.Init(DB)

	webDir, err := findWebDir()
	if err != nil {
		return err
	}

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	port := defaultPort
	if fromEnv := os.Getenv("TODO_PORT"); fromEnv != "" {
		if parsed, err := strconv.Atoi(fromEnv); err == nil {
			port = parsed
		}
	}

	addr := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(addr, nil)
}

func findWebDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	candidates := []string{
		filepath.Join(cwd, "web"),
		filepath.Join(cwd, "..", "web"),
	}

	for _, candidate := range candidates {
		info, err := os.Stat(candidate)
		if err == nil && info.IsDir() {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("web directory not found")
}
