package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	*sqlx.DB
}

func Connect(connString string) DB {
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		panic(fmt.Errorf("failed to connect db: %w", err))
	}
	return DB{db}
}
