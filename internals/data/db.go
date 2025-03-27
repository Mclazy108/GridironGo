package data

import (
	"GridironGo/internals/data/sqlc"
	"context"
	"database/sql"
	"log"
	"path/filepath"

	"github.com/mattn/go-sqlite3"
)

func InitDB(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	log.Println("Database connection established")
	return db
}

func NewQueries(db *sql.DB) *sqlc.Queries {
	return sqlc.New(db)
}

/*
funct GetPlayers(db *sql.dB) *sqlc.Queries{
	queries := NewQueries(db)
	query_context := context.Background()
	players, err := queries.GetPlayers(query_context)
}
*/
