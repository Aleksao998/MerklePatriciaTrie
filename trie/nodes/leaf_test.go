package nodes

import (
	"testing"

	"github.com/Aleksao998/Merkle-Patricia-Trie/trie/nibble"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewLeafNode_Basic tests the basic creation of a LeafNode with given path and value.
func TestNewLeafNode_Basic(t *testing.T) {
	t.Parallel()

	path := []nibble.Nibble{1, 2, 3, 4}
	value := []byte{5, 6, 7, 8}

	leaf := NewLeafNode(path, value)

	assert.Equal(t, path, leaf.Path, "The paths should be equal")
	assert.Equal(t, value, leaf.Value, "The values should be equal")
	assert.True(t, leaf.Dirty, "The values should be equal")
}

// TestNewLeafNode_Empty tests the creation of a LeafNode with an empty path and value.
func TestNewLeafNode_Empty(t *testing.T) {
	t.Parallel()

	path := []nibble.Nibble{}
	value := []byte{}

	leaf := NewLeafNode(path, value)

	require.Len(t, leaf.Path, 0, "The path should be empty")
	require.Len(t, leaf.Value, 0, "The value should be empty")
}
