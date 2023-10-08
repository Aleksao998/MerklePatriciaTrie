package trie

import (
	"bytes"
	"github.com/Aleksao998/Merkle-Patricia-Trie/core/storage/mockStorage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestInsert tests the insertion of various nodes into the Merkle Patricia Trie (MPT)
// It covers several edge cases for inserting different nodes and verifies the integrity
// of the trie by comparing its resultant hash against a known value from the go-ethereum
// official MPT implementation. The goal is to ensure that the insertion process is both
// accurate and consistent with established implementations.
func TestInsert(t *testing.T) {
	t.Parallel()
	db := &mockStorage.MockStorage{}
	trie := NewTrie(db)

	trie.Put([]byte("doe"), []byte("reindeer"))
	trie.Put([]byte("dog"), []byte("puppy"))
	trie.Put([]byte("dogglesworth"), []byte("cat"))

	trie.Put([]byte("cat"), []byte("cat"))
	trie.Put([]byte("car"), []byte("car"))

	trie.Put([]byte("do"), []byte("do"))
	trie.Put([]byte("dog"), []byte("dog"))

	trie.Put([]byte("c"), []byte("c"))
	trie.Put([]byte("cat"), []byte("cat1"))

	exp := common.HexToHash("ee85616d8e5799ae1d210f48d4661a9f0287656d8fc552113966a074b6bbf68f")
	assert.Equal(t, exp.Bytes(), trie.Hash(), "Unexpected root hash after insert operations")
}

// TestDelete tests the deletion of various nodes from the Merkle Patricia Trie (MPT)
// It covers several edge cases for deleting different nodes and verifies the integrity
// of the trie by comparing its resultant hash against a known value from the go-ethereum
// official MPT implementation. The goal is to ensure that the deletion process is both
// accurate and consistent with established implementations.
func TestDelete(t *testing.T) {
	t.Parallel()

	db := &mockStorage.MockStorage{}
	trie := NewTrie(db)
	trie.Put([]byte("doe"), []byte("reindeer"))
	trie.Put([]byte("dog"), []byte("puppy"))

	trie.Del([]byte("dog"))

	trie.Put([]byte("dog"), []byte("puppy"))
	trie.Put([]byte("do"), []byte("do"))

	trie.Del([]byte("do"))

	trie.Put([]byte("dor"), []byte("dor"))
	trie.Put([]byte("dos"), []byte("dos"))

	trie.Del([]byte("dos"))

	exp := common.HexToHash("39430803baabb9d662bd2c25905c6bdf9f5f8e1e2aed6cc1af732554816d55e0")

	assert.Equal(t, exp.Bytes(), trie.Hash(), "Unexpected root hash after delete operations")
}

// TestEmptyTree tests operations on an empty Merkle Patricia Trie (MPT)
// TODO in go-ethereum it hashes to some value
func TestEmptyTree(t *testing.T) {
	t.Skip()
	t.Parallel()

	db := &mockStorage.MockStorage{}
	trie := NewTrie(db)

	exp := common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
	root := trie.Hash()
	if !bytes.Equal(root, exp.Bytes()) {
		t.Errorf("case 1: exp %x got %x", exp, root)
	}
}

// TestSameTries tests that 2 trees with same nodes have same root hash
func TestSameTries(t *testing.T) {
	t.Parallel()

	db := &mockStorage.MockStorage{}
	trie := NewTrie(db)
	copyTrie := NewTrie(db)

	trie.Put([]byte("doe"), []byte("reindeer"))
	trie.Put([]byte("dog"), []byte("puppy"))
	trie.Del([]byte("dog"))
	trie.Put([]byte("dog"), []byte("puppy"))
	trie.Put([]byte("do"), []byte("do"))
	trie.Del([]byte("do"))

	copyTrie.Put([]byte("doe"), []byte("reindeer"))
	copyTrie.Put([]byte("dog"), []byte("puppy"))
	copyTrie.Put([]byte("do"), []byte("do"))
	copyTrie.Del([]byte("dog"))
	copyTrie.Del([]byte("do"))
	copyTrie.Put([]byte("dog"), []byte("puppy"))

	assert.Equal(t, trie.Hash(), copyTrie.Hash(), "Tries with same nodes are not identical")
}
