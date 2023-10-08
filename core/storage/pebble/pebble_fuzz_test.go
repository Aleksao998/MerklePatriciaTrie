package pebble

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func FuzzPebbleStorageWriteRead(f *testing.F) {
	tempDir, store, err := createPebbleStorage()
	if err != nil {
		f.Fatalf("error creating pebble storage, %v", err)
	}

	defer os.RemoveAll(tempDir)
	defer store.Close()

	f.Fuzz(func(t *testing.T, key string, value string) {
		t.Parallel()

		// Convert string to []byte as PebbleStorage expects []byte type for key and value
		keyBytes := []byte(key)
		valueBytes := []byte(value)

		// Test Put method
		err := store.Put(keyBytes, valueBytes)
		if err != nil {
			t.Fatalf("Error setting value for key '%s': %v", key, err)
		}

		// Test Get method
		retrievedValue, err := store.Get(keyBytes)
		if err != nil {
			t.Fatalf("Error getting value for key '%s': %v", key, err)
		}

		assert.Equal(t, valueBytes, retrievedValue)
	})
}
