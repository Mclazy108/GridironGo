package data

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

//go:embed queries/migrations/*.sql
var migrations embed.FS

// struct that has the database connection
type DB struct {
	*sql.DB
	Queries *Queries
}

type Matchup struct {
	HomeTeamID uint32
	AwayTeamID uint32
}

func initSchema(db *sql.DB) error {
	// Create schema from schema.sql
	schema, err := fs.ReadFile(migrations, "queries/migrations/schema.sql")
	if err != nil {
		return fmt.Errorf("Error in reading schema: %w", err)
	}

	// Execute schema.sql
	if _, err := db.Exec(string(schema)); err != nil {
		return fmt.Errorf("Error in executing schema : %w", err)
	}

	log.Println("Database schema initialized")
	return nil
}

// DBConnection create a new database connection and initializes the schema
func DBConnection(dbPath string) (*DB, error) {
	// Checks the directory exists
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create a directory: %w", err)
	}

	// Open the database
	sqlDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Check the Connection
	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Initialize the schema
	if err := initSchema(sqlDB); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("Error in initializing schema: %w", err)
	}

	// DB Struc
	db := &DB{
		DB:      sqlDB,
		Queries: DBConnection(sqlDB),
	}
	return db, nil
}

// Closes database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

func (db *DB) RunInTransaction(ctx context.Context, fn func(*Queries) error) error {
	transaction, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer transaction.Rollback()

	operation := DBConnection(transaction)
	if err := fn(operation); err != nil {
		return err
	}

	return transaction.Commit()
}
