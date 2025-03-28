package data

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"github.com/Mclazy108/GridironGo/internals/data/sqlc"
	"github.com/mattn/go-sqlite3"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"
)

//go:embed queries/migrations/*.sql
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
