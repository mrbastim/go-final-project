package main

import (
	"log"

	"github.com/mrbastim/go-final-project/back/db"
)

func main() {
	conn, err := db.OpenAndInit("scheduler.db")
	if err != nil {
		log.Fatalf("db init failed: %v", err)
	}
	defer conn.Close()

	log.Println("DB ready")
	// Дальше будет запуск HTTP-сервера (в следующих шагах проекта).
}
