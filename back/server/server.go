package server

import (
	"fmt"
	"main/back/api"
	"main/back/db"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

const defaultPort = 7540

func Run(storage db.TaskStorage) error {

	api.Init(storage)

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
