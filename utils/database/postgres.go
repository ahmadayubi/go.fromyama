package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/stdlib"
)

var DB *sql.DB

func ConnectToDatabase() error {
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}
	log.Print("Connected to Postgres Database")
	DB = db
	return nil
}