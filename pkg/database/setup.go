package database

import (
	"database/sql"
	"log"
	"task-one/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(cfg *config.Config) {
	var err error
	DB, err = sql.Open("postgres", cfg.DBDataSource())
	if err != nil {
		log.Fatalf("Failed to open a DB connection: %+v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Failed to ping DB: %+v", err)
	}

	log.Println("Database connection established")
}

func CloseDB() {
	if err := DB.Close(); err != nil {
		log.Fatalf("Failed to close DB connection: %+v", err)
	}
}
