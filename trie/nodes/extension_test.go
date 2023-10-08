package nodes

import (
	"github.com/Aleksao998/Merkle-Patricia-Trie/trie/nibble"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestNewExtension_Basic tests the creation of an ExtensionNode and its properties.
func TestNewExtension_Basic(t *testing.T) {
	t.Parallel()

	path := []nibble.Nibble{1, 2, 3, 4}
	leaf := NewLeafNode([]nibble.Nibble{5, 6, 7, 8}, []byte{9, 10, 11, 12})
	extNode := NewExtension(path, leaf)

	assert.Equal(t, path, extNode.Path, "The paths should be equal")
	assert.Equal(t, leaf, extNode.Node, "The nodes should be equal")
	assert.True(t, leaf.Dirty, "The values should be equal")
}

// TestNewExtension_EmptyPath tests the creation of an ExtensionNode with an empty path.
func TestNewExtension_EmptyPath(t *testing.T) {
	t.Parallel()

	path := []nibble.Nibble{}
	leaf := NewLeafNode([]nibble.Nibble{5, 6, 7, 8}, []byte{9, 10, 11, 12})
	extNode := NewExtension(path, leaf)

	assert.Empty(t, extNode.Path, "The path should be empty")
	assert.Equal(t, leaf, extNode.Node, "The nodes should be equal")
}

// TestNewExtension_EmptyNode tests the creation of an ExtensionNode with an empty (nil) node.
func TestNewExtension_EmptyNode(t *testing.T) {
	t.Parallel()

	path := []nibble.Nibble{1, 2, 3, 4}
	extNode := NewExtension(path, nil)

	assert.Equal(t, path, extNode.Path, "The paths should be equal")
	assert.Nil(t, extNode.Node, "The node should be nil")
}
