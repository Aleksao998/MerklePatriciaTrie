package mpt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMPTMemoryStorage_Has_NonExistentKey(t *testing.T) {
	storage := NewMPTMemoryStorage()
	key := []byte("nonexistent_key")

	has, err := storage.Has(key)
	assert.NoError(t, err)
	assert.False(t, has)
}

func TestMPTMemoryStorage_Put_Get_Valid(t *testing.T) {
	storage := NewMPTMemoryStorage()
	key := []byte("key")
	value := []byte("value")

	err := storage.Put(key, value)
	assert.NoError(t, err)

	retrievedValue, err := storage.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, value, retrievedValue)
}

func TestMPTMemoryStorage_Put_Overwrite(t *testing.T) {
	storage := NewMPTMemoryStorage()
	key := []byte("key")
	value1 := []byte("value1")
	value2 := []byte("value2")

	err := storage.Put(key, value1)
	assert.NoError(t, err)

	err = storage.Put(key, value2)
	assert.NoError(t, err)

	retrievedValue, err := storage.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, value2, retrievedValue)
}

func TestMPTMemoryStorage_Delete_NonExistentKey(t *testing.T) {
	storage := NewMPTMemoryStorage()
	key := []byte("nonexistent_key")

	err := storage.Delete(key)
	assert.NoError(t, err)
}

func TestMPTMemoryStorage_Delete_ExistentKey(t *testing.T) {
	storage := NewMPTMemoryStorage()
	key := []byte("key")
	value := []byte("value")

	err := storage.Put(key, value)
	assert.NoError(t, err)

	err = storage.Delete(key)
	assert.NoError(t, err)

	_, err = storage.Get(key)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "key not found")
}
