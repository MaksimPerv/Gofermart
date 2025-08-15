package database

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"path/filepath"
)

func RunMigrations(pool *pgxpool.Pool) error {
	migrationsDir := filepath.Join("internal", "migrations")
	goose.SetBaseFS(nil)
	conn := pool.Config().ConnConfig
	db := stdlib.OpenDB(*conn)
	return goose.Up(db, migrationsDir)
}
