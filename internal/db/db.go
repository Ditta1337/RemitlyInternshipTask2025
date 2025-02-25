package db

import (
	"context"
	"database/sql"
	"github.com/Ditta1337/RemitlyInternshipTask2025/internal/store"
	"time"
)

func New(addr string, maxOpenConns, maxIdleCons int, maxIdleTime string) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleCons)

	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func SeedDBIfEmpty(db *sql.DB, store store.Storage) error {
	ctx := context.Background()
	var count int

	query := `
		SELECT COUNT(*) 
		FROM banks
	`
	err := db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return Seed(store)
	}

	return nil
}
