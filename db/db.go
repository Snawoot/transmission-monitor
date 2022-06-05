package db

import (
	"errors"
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

func (d *DB) Clear() error {
	return d.db.DropAll()
}

func (d *DB) Get(key string) (b []byte, err error) {
	if err := d.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		if val, err := item.ValueCopy(nil); err != nil {
			return err
		} else {
			b = val
		}
		return nil
	}); err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get key %q: %w", key, err)
	}
	return
}

func (d *DB) Set(key string, value []byte) error {
	err := d.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), value)
	})
	if err != nil {
		return fmt.Errorf("failed to set key %q: %w", key, err)
	}
	return nil
}

func (d *DB) Delete(key string) error {
	err := d.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
	if err != nil {
		return fmt.Errorf("failed to delete key %q: %w", key, err)
	}
	return nil
}
