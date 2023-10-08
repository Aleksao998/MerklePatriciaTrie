package trie

import (
	"github.com/Aleksao998/Merkle-Patricia-Trie/core/storage/mockStorage"
	"github.com/stretchr/testify/require"
	"strconv"
	"sync"
	"testing"
)

// TestGetPutDelBasic tests basic trie instructions
func TestGetPutDelBasic(t *testing.T) {
	t.Parallel()

	t.Run("should get an error if key does not exists", func(t *testing.T) {
		t.Parallel()

		db := &mockStorage.MockStorage{}
		trie := NewTrie(db)

		_, err := trie.Get([]byte("notexist"))
		require.Error(t, errKeyNotFound, err)

		err = trie.Del([]byte("notexist"))
		require.Error(t, errKeyNotFound, err)
	})

	t.Run("should get an value if key exists", func(t *testing.T) {
		t.Parallel()

		db := &mockStorage.MockStorage{}
		trie := NewTrie(db)

		trie.Put([]byte("key"), []byte("value"))

		val, err := trie.Get([]byte("key"))
		require.Nil(t, err)
		require.Equal(t, []byte("value"), val)
	})

	t.Run("should get an error if we try to get deleted item", func(t *testing.T) {
		t.Parallel()

		db := &mockStorage.MockStorage{}
		trie := NewTrie(db)

		trie.Put([]byte("key"), []byte("value"))
		trie.Del([]byte("key"))

		_, err := trie.Get([]byte("notexist"))
		require.Error(t, errKeyNotFound, err)
	})

	t.Run("should get latest value on updated items", func(t *testing.T) {
		t.Parallel()

		db := &mockStorage.MockStorage{}
		trie := NewTrie(db)

		trie.Put([]byte("key"), []byte("value"))
		trie.Put([]byte("key1"), []byte("value1"))
		trie.Put([]byte("key"), []byte("new value"))

		val, err := trie.Get([]byte("key"))
		require.Nil(t, err)
		require.Equal(t, []byte("new value"), val)
	})
}

// TestGetPutMultipleItems tests adding and retrieving multiple items
func TestGetPutMultipleItems(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		key   []byte
		value []byte
	}{
		{key: []byte("dog"), value: []byte("dog")},
		{key: []byte("doge"), value: []byte("doge")},
		{key: []byte("dogger"), value: []byte("dogger")},
		{key: []byte("cat"), value: []byte("cat")},
		{key: []byte("car"), value: []byte("car")},
	}

	db := &mockStorage.MockStorage{}
	trie := NewTrie(db)

	// put items
	for _, tt := range testCases {
		trie.Put(tt.key, tt.value)
	}

	// get items
	for _, tt := range testCases {
		val, err := trie.Get(tt.key)
		require.Nil(t, err)
		require.Equal(t, tt.value, val)
	}
}

// TestTrieConcurrencyBasic tests adding, retrieving and deleting item in parallel
func TestTrieConcurrencyBasic(t *testing.T) {
	db := &mockStorage.MockStorage{}
	trie := NewTrie(db)

	// use WaitGroup to synchronize our go routines
	var wg sync.WaitGroup

	// launch multiple goroutines to write data to the trie
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "key" + strconv.Itoa(i)
			value := "value" + strconv.Itoa(i)
			trie.Put([]byte(key), []byte(value))
		}(i)
	}

	// wait for all inserts to finish
	wg.Wait()

	// launch multiple goroutines to read data from the trie
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "key" + strconv.Itoa(i)
			value, err := trie.Get([]byte(key))
			require.NoError(t, err, "Error reading key %s", key)
			expectedValue := "value" + strconv.Itoa(i)
			require.Equal(t, expectedValue, string(value), "Expected value for key %s", key)
		}(i)
	}

	// wait for all goroutines to finish
	wg.Wait()

	// launch multiple goroutines to read data from the trie
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "key" + strconv.Itoa(i)
			err := trie.Del([]byte(key))
			require.NoError(t, err, "Error deleting key %s", key)
		}(i)
	}

	// wait for all goroutines to finish
	wg.Wait()
}

// TestTrieConcurrencyAdvance tests getting and putting items in parallel
func TestTrieConcurrencyAdvance(t *testing.T) {
	db := &mockStorage.MockStorage{}
	trie := NewTrie(db)

	// use WaitGroup to synchronize our go routines
	var wg sync.WaitGroup

	// launch multiple goroutines to write initial data to the trie
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "key" + strconv.Itoa(i)
			value := "value" + strconv.Itoa(i)
			trie.Put([]byte(key), []byte(value))
		}(i)
	}

	// wait for all initial inserts to finish
	wg.Wait()

	// launch multiple goroutines to concurrently read old items and add new items with an offset
	offset := 1000
	for i := 0; i < 1000; i++ {
		// retrieve old items
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "key" + strconv.Itoa(i)
			value, err := trie.Get([]byte(key))
			require.NoError(t, err, "Error reading key %s", key)
			expectedValue := "value" + strconv.Itoa(i)
			require.Equal(t, expectedValue, string(value), "Expected value for key %s", key)
		}(i)

		// insert new items with offset
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "key" + strconv.Itoa(i+offset) // Using offset for new keys
			value := "value" + strconv.Itoa(i+offset)
			trie.Put([]byte(key), []byte(value))
		}(i)
	}

	// wait for all goroutines to finish
	wg.Wait()

	// launch multiple goroutines to read the new items added with offset
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "key" + strconv.Itoa(i+offset)
			value, err := trie.Get([]byte(key))
			require.NoError(t, err, "Error reading new key %s", key)
			expectedValue := "value" + strconv.Itoa(i+offset)
			require.Equal(t, expectedValue, string(value), "Expected value for new key %s", key)
		}(i)
	}

	// wait for all goroutines to finish
	wg.Wait()
}
