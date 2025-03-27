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

/*
type DB struct {
	*sql.DB
	Queries *Queries
}
*/
