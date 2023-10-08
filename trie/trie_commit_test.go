package trie

import (
	"github.com/Aleksao998/Merkle-Patricia-Trie/storage/mpt"
	"github.com/ethereum/go-ethereum/common"
	ethereumTrie "github.com/ethereum/go-ethereum/trie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestGetPutDelBasicAfterCommit tests basic trie instructions after commit
func TestGetPutDelBasicAfterCommit(t *testing.T) {
	t.Parallel()

	t.Run("should get an value if key exists", func(t *testing.T) {
		t.Parallel()

		db := mpt.NewMPTMemoryStorage()
		trie := NewTrie(db)

		trie.Put([]byte("dog"), []byte("dog"))
		trie.Put([]byte("dor"), []byte("dor"))

		_ = trie.Commit()

		val, err := trie.Get([]byte("dog"))
		require.Nil(t, err)
		require.Equal(t, []byte("dog"), val)

		val, err = trie.Get([]byte("missing_key"))
		require.Nil(t, val)
		require.Error(t, err)
	})

	t.Run("root hash should match has after commit", func(t *testing.T) {
		t.Parallel()

		db := mpt.NewMPTMemoryStorage()
		trie := NewTrie(db)

		trie.Put([]byte("dog"), []byte("dog"))
		trie.Put([]byte("dor"), []byte("dor"))

		originalHash := trie.Hash()

		_ = trie.Commit()

		newHash := trie.Hash()

		assert.Equal(t, originalHash, newHash)
	})

	t.Run("proof should be valid before and after commit", func(t *testing.T) {
		t.Parallel()

		db := mpt.NewMPTMemoryStorage()
		trie := NewTrie(db)

		trie.Put([]byte("dog"), []byte("dog"))
		trie.Put([]byte("dor"), []byte("dor"))

		_ = trie.Commit()

		keyToProof := []byte("dog")
		proof, err := trie.Proof(keyToProof)

		assert.Nil(t, err, "Failed to generate proof")
		hashByte := common.BytesToHash(trie.Hash())

		value, err := ethereumTrie.VerifyProof(hashByte, keyToProof, proof)

		assert.Nil(t, err, "Proof verification failed")
		expectedValue, _ := trie.Get(keyToProof)
		assert.Equal(t, expectedValue, value, "Mismatch in expected and retrieved value")
	})

	t.Run("insert should be valid before and after commit", func(t *testing.T) {
		t.Parallel()

		db := mpt.NewMPTMemoryStorage()
		originalTrie := NewTrie(db)

		originalTrie.Put([]byte("dog"), []byte("dog"))
		originalTrie.Put([]byte("dor"), []byte("dor"))
		originalTrie.Put([]byte("dogger"), []byte("dor"))
		originalTrie.Put([]byte("cat"), []byte("cat"))

		originalHash := originalTrie.Hash()

		trie := NewTrie(db)

		trie.Put([]byte("dog"), []byte("dog"))
		trie.Put([]byte("dor"), []byte("dor"))
		trie.Commit()
		trie.Put([]byte("dogger"), []byte("dor"))
		trie.Commit()
		trie.Put([]byte("cat"), []byte("cat"))

		newHash := trie.Hash()
		assert.Equal(t, originalHash, newHash, "Mismatch in oriignal and new hash")
	})

	t.Run("delete should have affect after commit", func(t *testing.T) {
		t.Parallel()

		db := mpt.NewMPTMemoryStorage()
		originalTrie := NewTrie(db)

		originalTrie.Put([]byte("dog"), []byte("dog"))
		originalTrie.Put([]byte("dor"), []byte("dor"))

		originalHash := originalTrie.Hash()

		trie := NewTrie(db)

		trie.Put([]byte("dog"), []byte("dog"))
		trie.Put([]byte("dor"), []byte("dor"))
		trie.Put([]byte("dogger"), []byte("dor"))
		trie.Commit()

		err := trie.Del([]byte("dogger"))
		require.Nil(t, err)

		newHash := trie.Hash()
		assert.Equal(t, originalHash, newHash, "Mismatch in oriignal and new hash")
	})
}
