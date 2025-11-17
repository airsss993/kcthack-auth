package database

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kcthack-auth/internal/config"
)

func ConnDB(cfg *config.Config) *sql.DB {
	db, err := sql.Open("pgx", cfg.Database.DSN)
	if err != nil {
		log.Fatalf("failed to connect to PostgreSQL: %s", err.Error())
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %s", err.Error())
	}

	log.Println("Successfully connected to DB!")

	return db
}
