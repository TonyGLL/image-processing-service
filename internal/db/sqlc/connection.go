package db

import (
	"database/sql"
)

type Store interface {
	UsersQuerier
	AuthQuerier
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}
