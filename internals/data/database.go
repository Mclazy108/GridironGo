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
	// If no path is specified, use GridironGo.db in the current directory
	if dbPath == "" {
		dbPath = "./GridironGo.db"
		log.Printf("No database path specified, using default: %s", dbPath)
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
		if config.Path != "" {
			cfg.Path = config.Path
		} else {
			// Ensure path is never empty
			cfg.Path = "./GridironGo.db"
			log.Printf("Warning: Empty database path in config, using default: %s", cfg.Path)
		}

		// Copy other config settings if provided
		if config.ForeignKeys {
			cfg.ForeignKeys = config.ForeignKeys
		}
		if config.JournalMode != "" {
			cfg.JournalMode = config.JournalMode
		}
		if config.BusyTimeout > 0 {
			cfg.BusyTimeout = config.BusyTimeout
		}
		if config.MaxOpenConns > 0 {
			cfg.MaxOpenConns = config.MaxOpenConns
		}
		if config.MaxIdleConns > 0 {
			cfg.MaxIdleConns = config.MaxIdleConns
		}
		if config.ConnMaxLifetime > 0 {
			cfg.ConnMaxLifetime = config.ConnMaxLifetime
		}
	}

	// Debug: Print path being used
	log.Printf("Using database path: %s", cfg.Path)

	// Check if database file exists
	if _, err := os.Stat(cfg.Path); err == nil {
		log.Printf("Database file already exists at: %s", cfg.Path)
	} else {
		log.Printf("Database file does not exist, will be created at: %s", cfg.Path)
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

	log.Printf("SQLite connection string: %s", connStr)

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

	log.Println("Successfully connected to database")

	// Apply schema migrations from the embedded schema.sql file
	log.Println("Checking if tables already exist...")

	// Check if tables already exist
	var nflGameCount, nflTeamCount, nflPlayerCount int
	err = db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='nfl_games'").Scan(&nflGameCount)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("error checking if nfl_games table exists: %w", err)
	}

	err = db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='nfl_teams'").Scan(&nflTeamCount)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("error checking if nfl_teams table exists: %w", err)
	}

	err = db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='nfl_players'").Scan(&nflPlayerCount)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("error checking if nfl_players table exists: %w", err)
	}

	if nflGameCount > 0 && nflTeamCount > 0 && nflPlayerCount > 0 {
		log.Println("All required tables already exist, skipping schema creation")
	} else {
		// Read schema file
		schemaSQL, err := migrationFS.ReadFile("migrations/schema.sql")
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to read schema.sql: %w", err)
		}

		log.Println("Applying schema migrations from schema.sql...")
		_, err = db.Exec(string(schemaSQL))
		if err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to execute schema migrations: %w", err)
		}

		log.Println("Schema migrations applied successfully")
	}

	// Verify the nfl_games table was created
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='nfl_games'").Scan(&count)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("error checking if nfl_games table exists: %w", err)
	}

	if count == 0 {
		db.Close()
		return nil, fmt.Errorf("failed to create nfl_games table")
	}

	log.Println("Successfully verified nfl_games table exists")

	// Verify the nfl_teams table was created
	count = 0
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='nfl_teams'").Scan(&count)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("error checking if nfl_teams table exists: %w", err)
	}

	if count == 0 {
		db.Close()
		return nil, fmt.Errorf("failed to create nfl_teams table")
	}

	log.Println("Successfully verified nfl_teams table exists")

	// Verify the nfl_players table was created
	count = 0
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='nfl_players'").Scan(&count)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("error checking if nfl_players table exists: %w", err)
	}

	if count == 0 {
		db.Close()
		return nil, fmt.Errorf("failed to create nfl_players table")
	}

	log.Println("Successfully verified nfl_players table exists")

	// Create sqlc queries
	queries := sqlc.New(db)

	return &DB{
		DB:      db,
		Queries: queries,
		config:  cfg,
	}, nil
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
