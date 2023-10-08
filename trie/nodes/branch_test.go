package nodes

import (
	"github.com/Aleksao998/Merkle-Patricia-Trie/trie/nibble"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestNewBranchNode_Basic tests the basic creation of a BranchNode.
func TestNewBranchNode_Basic(t *testing.T) {
	t.Parallel()

	branch := NewBranchNode()
	require.NotNil(t, branch, "Branch node should not be nil")
	assert.Empty(t, branch.Value, "New branch node should not have value")
	assert.Equal(t, 0, branch.ChildCount(), "New branch node should not have children")
}

// TestSetValueAndGet_BranchNode tests the setting and getting value of a BranchNode.
func TestSetValueAndGet_BranchNode(t *testing.T) {
	t.Parallel()

	branch := NewBranchNode()
	value := []byte{5, 6, 7, 8}

	branch.SetValue(value)
	gotValue, exists := branch.GetValue()

	assert.True(t, exists, "Value should exist in the branch node")
	assert.Equal(t, value, gotValue, "Values should match")
}

// TestClearValue_BranchNode tests clearing the value from a BranchNode.
func TestClearValue_BranchNode(t *testing.T) {
	t.Parallel()

	branch := NewBranchNode()
	value := []byte{5, 6, 7, 8}
	branch.SetValue(value)

	branch.ClearValue()
	assert.Nil(t, branch.Value, "Value should be cleared")
	assert.False(t, branch.HasValue(), "Branch node should not have a value")
}

// TestSetChild_BranchNode tests the setting of a child in a BranchNode.
func TestSetChild_BranchNode(t *testing.T) {
	t.Parallel()

	branch := NewBranchNode()
	leaf := NewLeafNode([]nibble.Nibble{1, 2}, []byte{3, 4})

	branch.SetChild(5, leaf)
	assert.Equal(t, leaf, branch.Children[5], "Child node should match the set node")
	assert.Equal(t, 1, branch.ChildCount(), "Child count should be 1")
}

// TestChildCount_BranchNode tests the counting of children in a BranchNode.
func TestChildCount_BranchNode(t *testing.T) {
	t.Parallel()

	branch := NewBranchNode()
	leaf1 := NewLeafNode([]nibble.Nibble{1, 2}, []byte{3, 4})
	leaf2 := NewLeafNode([]nibble.Nibble{5, 6}, []byte{7, 8})

	branch.SetChild(5, leaf1)
	branch.SetChild(6, leaf2)
	assert.Equal(t, 2, branch.ChildCount(), "Child count should be 2")
}

// TestSetChildInvalidNibble_BranchNode tests the behavior when setting a child with an invalid nibble.
func TestSetChildInvalidNibble_BranchNode(t *testing.T) {
	t.Parallel()

	branch := NewBranchNode()
	leaf := NewLeafNode([]nibble.Nibble{1, 2}, []byte{3, 4})

	assert.Panics(t, func() { branch.SetChild(16, leaf) }, "Should panic with invalid nibble")
}
