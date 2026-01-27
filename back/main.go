package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"main/back/db"
)

func main() {
	_ = godotenv.Load()

	dbFile := "scheduler.db"
	if fromEnv := os.Getenv("TODO_DBFILE"); fromEnv != "" {
		dbFile = fromEnv
	}

	conn, err := db.OpenAndInit(dbFile)
	if err != nil {
		log.Fatalf("db init failed: %v", err)
	}
	defer conn.Close()

	log.Println("DB ready")
	// Дальше будет запуск HTTP-сервера (в следующих шагах проекта).
}
