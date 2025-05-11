package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"sync"

	"testing"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed init.sql
var initSQL embed.FS

var (
	db   *sql.DB
	once sync.Once
)

type Config struct {
	Path string
}

func NewConfig(path string) *Config {
	return &Config{Path: path}
}

func Connect(cfg *Config) (*sql.DB, error) {
	if testing.Testing() {
		return connectDB(cfg.Path)
	}

	var err error
	once.Do(func() {
		db, err = connectDB(cfg.Path)
	})

	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return db, nil
}

func connectDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	sqlBytes, err := initSQL.ReadFile("init.sql")
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("database error: %w", err)
	}

	if _, err = db.Exec(string(sqlBytes)); err != nil {
		db.Close()
		return nil, fmt.Errorf("database error: %w", err)
	}

	return db, nil
}

func GetDB() *sql.DB {
	return db
}

func WithTx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
