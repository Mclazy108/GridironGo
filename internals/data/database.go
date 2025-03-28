package data

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Mclazy108/GridironGo/internals/data/sqlc"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

type DBConfig struct {
	Path            string
	ForeignKeys     bool
	JournalMode     string
	BusyTimeout     int
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

func DefaultDBConfig(dbPath string) DBConfig {
	// If no path is specified, use GridironGo.db in the executable's directory
	if dbPath == "" {
		// Get the executable's directory
		exePath, err := os.Executable()
		if err == nil {
			// Use the directory of the executable
			exeDir := filepath.Dir(exePath)
			// Go to the parent directory (assuming executables is a subdirectory)
			parentDir := filepath.Dir(exeDir)
			// If the parent directory contains build.sh, use that (root directory)
			if _, err := os.Stat(filepath.Join(parentDir, "build.sh")); err == nil {
				dbPath = filepath.Join(parentDir, "GridironGo.db")
			} else {
				// Otherwise just use the current directory
				dbPath = "./GridironGo.db"
			}
		} else {
			// Fallback to the current directory if we can't determine the executable path
			dbPath = "./GridironGo.db"
		}
	}

	return DBConfig{
		Path:            dbPath,
		ForeignKeys:     true,
		JournalMode:     "WAL",
		BusyTimeout:     5000,
		MaxOpenConns:    1,
		MaxIdleConns:    1,
		ConnMaxLifetime: 3600,
	}
}

type DB struct {
	*sql.DB
	*sqlc.Queries
	config DBConfig
}

// NewDB creates a new database connection and prepares the database.
// If config is nil, it uses a default configuration with a database at "./GridironGo.db"
func NewDB(config *DBConfig) (*DB, error) {
	// Use default config if none provided
	cfg := DefaultDBConfig("")
	if config != nil {
		cfg = *config
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(cfg.Path)
	if dir != "." && dir != "/" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory for database: %w", err)
		}
	}

	// Connect to database with correct parameters
	connStr := fmt.Sprintf("file:%s?_foreign_keys=%t&_journal_mode=%s&_busy_timeout=%d",
		cfg.Path, cfg.ForeignKeys, cfg.JournalMode, cfg.BusyTimeout)

	db, err := sql.Open("sqlite3", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	// Check connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Apply migrations
	if err := applyMigrations(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	// Create sqlc queries
	queries := sqlc.New(db)

	return &DB{
		DB:      db,
		Queries: queries,
		config:  cfg,
	}, nil
}

// tablesExist checks if the database schema has already been applied
func tablesExist(db *sql.DB) (bool, error) {
	// Check for the existence of one of your tables
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='seasons'").Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// applyMigrations applies all SQL migrations in the migrations directory
func applyMigrations(db *sql.DB) error {
	// Check if tables already exist
	exists, err := tablesExist(db)
	if err != nil {
		return fmt.Errorf("failed to check if tables exist: %w", err)
	}

	// Skip migrations if tables already exist
	if exists {
		log.Println("Database schema already exists, skipping migrations")
		return nil
	}

	// Get migrations from embedded filesystem
	migrations, err := migrationFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Start transaction for migrations
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Apply migrations in order
	for _, migration := range migrations {
		if !migration.IsDir() && filepath.Ext(migration.Name()) == ".sql" {
			migrationPath := filepath.Join("migrations", migration.Name())
			migrationSQL, err := migrationFS.ReadFile(migrationPath)
			if err != nil {
				return fmt.Errorf("failed to read migration %s: %w", migration.Name(), err)
			}

			if _, err := tx.Exec(string(migrationSQL)); err != nil {
				return fmt.Errorf("failed to execute migration %s: %w", migration.Name(), err)
			}

			log.Printf("Applied migration: %s", migration.Name())
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migrations: %w", err)
	}

	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

// ExecTx executes a function within a database transaction
func (db *DB) ExecTx(ctx context.Context, fn func(*sqlc.Queries) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := db.Queries.WithTx(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx failed: %v, rollback failed: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
