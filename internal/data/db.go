package data

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func NewDB(dataSourceName string) *sql.DB { //пример dataSourceName: postgres://user:password@localhost:5432/mydb?sslmode=disable
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Connected to PostgreSQL")
	return db
}
