package db

import (
	"fmt"

	"github.com/dgraph-io/badger/v3"
)

type DB struct {
	db *badger.DB
}

func Open(path string) (*DB, error) {
	db, err := badger.Open(badger.DefaultOptions(path).WithLoggingLevel(badger.WARNING))
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return &DB{
		db: db,
	}, nil
}

func (d *DB) Close() error {
	return d.db.Close()
}
