package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() error {
	var err error
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error al abrir conexión: %v", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("error al verificar conexión: %v", err)
	}

	fmt.Println("✓ Conectado a PostgreSQL")
	return nil
}