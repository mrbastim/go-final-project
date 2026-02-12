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
		log.Printf("db init failed: %v", err)
		return
	}
	defer conn.Close()
	storage := db.NewStorage(conn)

	log.Println("DB ready")
	if err := server.Run(storage); err != nil {
		log.Printf("server failed: %v", err)
		return
	}

}
