package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"main/back/db"
	"main/back/server"
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
	if err := server.Run(conn); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
