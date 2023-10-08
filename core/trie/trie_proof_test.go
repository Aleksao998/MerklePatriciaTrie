package trie

import (
	"github.com/Aleksao998/Merkle-Patricia-Trie/core/storage/mockStorage"
	"github.com/ethereum/go-ethereum/common"
	ethereumTrie "github.com/ethereum/go-ethereum/trie"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestProofVerification tests the MPT's proof generation and verification mechanism
// for a key that exists in the trie. It ensures that the trie produces accurate proofs
// that can be verified using the go-ethereum implementation.
func TestProofVerification(t *testing.T) {
	t.Parallel()

	db := &mockStorage.MockStorage{}
	trie := NewTrie(db)

	keys := [][]byte{
		[]byte("doe"), []byte("dog"), []byte("dogglesworth"),
		[]byte("cat"), []byte("car"), []byte("do"),
	}
	values := [][]byte{
		[]byte("reindeer"), []byte("puppy"), []byte("cat"),
		[]byte("cat"), []byte("car"), []byte("do"),
	}

	for i := range keys {
		trie.Put(keys[i], values[i])
	}

	keyToProof := []byte("dog")
	proof, err := trie.Proof(keyToProof)

	assert.Nil(t, err, "Failed to generate proof")
	hashByte := common.BytesToHash(trie.Hash())

	value, err := ethereumTrie.VerifyProof(hashByte, keyToProof, proof)

	assert.Nil(t, err, "Proof verification failed")
	expectedValue, _ := trie.Get(keyToProof)
	assert.Equal(t, expectedValue, value, "Mismatch in expected and retrieved value")
}

// TestProofVerificationForNonExistentKey tests the MPT's behavior when generating a proof
// for a key that was never added to the trie. The test ensures that the trie's proof
// generation mechanism accurately represents the absence of a key.
func TestProofVerificationForNonExistentKey(t *testing.T) {
	t.Parallel()

	db := &mockStorage.MockStorage{}
	trie := NewTrie(db)

	key := []byte("neverExists")

	proof, err := trie.Proof(key)
	assert.Error(t, err)
	hashByte := common.BytesToHash(trie.Hash())

	_, err = ethereumTrie.VerifyProof(hashByte, key, proof)
	assert.NotNil(t, err, "Proof verification should have failed for a non-existent key")
}

// TestProofVerificationForOverwrittenKey tests the MPT's ability to produce a valid proof
// for a key whose associated value has been overwritten. The test ensures that the proof
// represents the latest value associated with the key.
func TestProofVerificationForOverwrittenKey(t *testing.T) {
	t.Parallel()

	db := &mockStorage.MockStorage{}
	trie := NewTrie(db)

	key := []byte("overwriteMe")
	value1 := []byte("originalValue")
	value2 := []byte("newValue")

	trie.Put(key, value1)
	trie.Put(key, value2) // Overwrite

	proof, err := trie.Proof(key)
	assert.Nil(t, err, "Failed to generate proof")

	hashByte := common.BytesToHash(trie.Hash())
	value, err := ethereumTrie.VerifyProof(hashByte, key, proof)

	assert.Nil(t, err, "Proof verification failed")
	assert.Equal(t, value2, value, "Mismatch in expected and retrieved value after overwriting")
}

// TestProofVerificationAfterDeletion tests the MPT's ability to produce
// a valid proof for a key even after its associated value has been deleted. This ensures
// the trie's integrity and demonstrates the immutability characteristic of the MPT.
func TestProofVerificationAfterDeletion(t *testing.T) {
	t.Parallel()

	db := &mockStorage.MockStorage{}
	trie := NewTrie(db)

	key := []byte("deleteMe")
	value := []byte("toBeDeleted")

	trie.Put(key, value)
	trie.Del(key)

	proof, err := trie.Proof(key)
	assert.Error(t, err)

	hashByte := common.BytesToHash(trie.Hash())

	_, err = ethereumTrie.VerifyProof(hashByte, key, proof)
	assert.NotNil(t, err, "Proof verification should have failed for a deleted key")
}
