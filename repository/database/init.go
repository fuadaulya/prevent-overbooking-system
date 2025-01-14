package repository

import (
	"database/sql"
	"task-one/config"
	db "task-one/pkg/database"
)

type Postgres struct {
	DB *sql.DB
}

func InitDatabase() {
	cfg := config.GetConfig()
	db.InitDB(cfg)
}

// NewRepository menginisialisasi URLRepository dengan koneksi database.
func NewRepository() *Postgres {
	return &Postgres{
		DB: db.DB,
	}
}
