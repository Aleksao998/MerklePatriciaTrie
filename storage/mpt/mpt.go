package mpt

import (
	"errors"
	"sync"
)

type MPTMemoryStorage struct {
	data map[string][]byte
	mu   sync.RWMutex
}

func NewMPTMemoryStorage() *MPTMemoryStorage {
	return &MPTMemoryStorage{
		data: make(map[string][]byte),
	}
}

func (m *MPTMemoryStorage) Has(key []byte) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, ok := m.data[string(key)]

	return ok, nil
}

func (m *MPTMemoryStorage) Get(key []byte) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, ok := m.data[string(key)]
	if !ok {
		return nil, errors.New("key not found")
	}

	return value, nil
}

func (m *MPTMemoryStorage) Put(key []byte, value []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[string(key)] = value

	return nil
}

func (m *MPTMemoryStorage) Delete(key []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.data, string(key))

	return nil
}
