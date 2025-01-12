package db

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	ErrNotFound  = errors.New("not found")
	ErrDuplicate = errors.New("duplicate")
)

type DB struct {
	db *sqlx.DB
}

func Connect(connString string) DB {
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		panic(fmt.Errorf("failed to connect db: %w", err))
	}
	return DB{db}
}
