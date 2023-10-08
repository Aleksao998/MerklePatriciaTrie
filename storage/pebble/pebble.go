package pebble

import (
	"github.com/cockroachdb/pebble"
)

type Storage struct {
	db *pebble.DB
}

// NewStorage initializes a new Storage instance with a database at the given path
func NewStorage(path string) (*Storage, error) {
	db, err := pebble.Open(path, &pebble.Options{})
	if err != nil {
		return nil, err
	}

	return &Storage{
		db: db,
	}, nil
}

// Has retrieves if a key is present in the key-value data store.
func (p *Storage) Has(key []byte) (bool, error) {
	_, closer, err := p.db.Get(key)
	if err == pebble.ErrNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	closer.Close()
	return true, nil
}

// Get retrieves the value for a given key and returns an error if any issue occurs during the operation
func (p *Storage) Get(key []byte) ([]byte, error) {
	value, closer, err := p.db.Get(key)
	if err != nil {
		return nil, err
	}
	defer closer.Close()

	return value, nil
}

// Put inserts the given value into the key-value data store.
func (p *Storage) Put(key []byte, value []byte) error {
	return p.db.Set(key, value, pebble.Sync)
}

// Delete removes the key from the key-value data store.
func (p *Storage) Delete(key []byte) error {
	return p.db.Delete(key, pebble.Sync)
}

// Close closes the database connection and returns an error if any issue occurs during the operation
func (p *Storage) Close() error {
	return p.db.Close()
}
