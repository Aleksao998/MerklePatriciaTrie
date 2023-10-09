package pebble

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/stretchr/testify/assert"
)

func createPebbleStorage() (string, *Storage, error) {
	// Create a temporary directory for Pebble storage
	tempDir, err := os.MkdirTemp("", "pebble-test")
	if err != nil {
		return "", nil, fmt.Errorf("error creating temporary directory: %w", err)
	}

	// Initialize a new Pebble storage instance
	store, err := NewStorage(tempDir)
	if err != nil {
		os.RemoveAll(tempDir)

		return "", nil, fmt.Errorf("error creating new Pebble storage: %w", err)
	}

	return tempDir, store, nil
}

func TestPebbleStorage_GetNonExistentKey(t *testing.T) {
	t.Parallel()

	// Initialize PebbleStorage
	tempDir, store, err := createPebbleStorage()
	if err != nil {
		t.Fatalf("error creating pebble storage, %v", err)
	}

	defer os.RemoveAll(tempDir)
	defer store.Close()

	// Test Get on non-existent key
	nonExistentKey := []byte("non_existent_key")

	_, err = store.Get(nonExistentKey)
	if !assert.ErrorIs(t, err, pebble.ErrNotFound) {
		t.Errorf("Expected error not found when getting non-existent key")
	}
}

func TestPebbleStorage_WriteRead(t *testing.T) {
	t.Parallel()

	// Initialize PebbleStorage
	tempDir, store, err := createPebbleStorage()
	if err != nil {
		t.Fatalf("error creating pebble storage, %v", err)
	}

	defer os.RemoveAll(tempDir)
	defer store.Close()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Test Put and Get methods for multiple key-value pairs
	for i := 0; i < 10; i++ {
		key := []byte(fmt.Sprintf("key_%d", i))
		value := make([]byte, r.Intn(100000))

		_, err := r.Read(value)
		if err != nil {
			t.Fatalf("Error generating random number: %v", err)
		}

		if err := store.Put(key, value); err != nil {
			t.Fatalf("Error setting value for key '%s': %v", string(key), err)
		}

		retrievedValue, err := store.Get(key)
		if err != nil {
			t.Fatalf("Error getting value for key '%s': %v", string(key), err)
		}

		assert.Equal(t, value, retrievedValue)
	}
}

// Test for Has method
func TestPebbleStorage_Has(t *testing.T) {
	t.Parallel()

	// Initialize PebbleStorage
	tempDir, store, err := createPebbleStorage()
	if err != nil {
		t.Fatalf("error creating pebble storage, %v", err)
	}

	defer os.RemoveAll(tempDir)
	defer store.Close()

	key := []byte("test_key")
	value := []byte("test_value")

	// Check for non-existent key
	has, err := store.Has(key)
	assert.NoError(t, err)
	assert.False(t, has)

	// Put key-value pair
	err = store.Put(key, value)
	assert.NoError(t, err)

	// Check for existent key
	has, err = store.Has(key)
	assert.NoError(t, err)
	assert.True(t, has)
}

// Test for Delete method
func TestPebbleStorage_Delete(t *testing.T) {
	t.Parallel()

	// Initialize PebbleStorage
	tempDir, store, err := createPebbleStorage()
	if err != nil {
		t.Fatalf("error creating pebble storage, %v", err)
	}

	defer os.RemoveAll(tempDir)
	defer store.Close()

	key := []byte("test_key")
	value := []byte("test_value")

	// Put key-value pair
	err = store.Put(key, value)
	assert.NoError(t, err)

	// Delete key
	err = store.Delete(key)
	assert.NoError(t, err)

	// Check for non-existent key
	_, err = store.Get(key)
	assert.ErrorIs(t, err, pebble.ErrNotFound)
}
