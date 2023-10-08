package trie

import (
	"fmt"
	"github.com/Aleksao998/Merkle-Patricia-Trie/storage/mpt"
	"github.com/Aleksao998/Merkle-Patricia-Trie/trie/nibble"
	nodes2 "github.com/Aleksao998/Merkle-Patricia-Trie/trie/nodes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommitter_SingleLeafNode(t *testing.T) {
	// Initialize memory storage and committer
	storage := mpt.NewMPTMemoryStorage()
	trie := NewTrie(storage)

	// Create a leaf node
	leaf := &nodes2.LeafNode{
		Path:  nibble.FromBytes([]byte{1, 2, 3, 4}),
		Value: []byte("value"),
	}

	// Commit the node
	hash, err := trie.commit(leaf)
	assert.NoError(t, err, "Failed to commit leaf node")

	// Decode the node
	decodedNode, err := trie.DecodeNode(hash)
	assert.NoError(t, err, "Failed to decode node")

	// Compare the original and decoded nodes
	decodedLeaf := decodedNode.(*nodes2.LeafNode)
	assert.Equal(t, leaf.Path, decodedLeaf.Path, "Leaf path mismatch")
	assert.Equal(t, leaf.Value, decodedLeaf.Value, "Leaf value mismatch")
}

func TestCommitter_ExtensionAndBranchNodes(t *testing.T) {
	// Initialize memory storage and committer
	storage := mpt.NewMPTMemoryStorage()
	trie := NewTrie(storage)

	// Create nodes
	leaf1 := &nodes2.LeafNode{Path: nibble.FromBytes([]byte{1, 2}), Value: []byte("leaf1")}
	leaf2 := &nodes2.LeafNode{Path: nibble.FromBytes([]byte{2, 3}), Value: []byte("leaf2")}
	branch := &nodes2.BranchNode{}
	branch.Children[1] = leaf1
	branch.Children[2] = leaf2
	ext := &nodes2.ExtensionNode{
		Path: nibble.FromBytes([]byte{0}),
		Node: branch,
	}

	// Commit the extension node (which also commits the branch and leaf nodes)
	hash, err := trie.commit(ext)
	assert.NoError(t, err, "Failed to commit extension node")

	// Decode the node
	decodedNode, err := trie.DecodeNode(hash)
	assert.NoError(t, err, "Failed to decode node")

	// Compare original and decoded structures
	decodedExt := decodedNode.(*nodes2.ExtensionNode)
	assert.Equal(t, ext.Path, decodedExt.Path, "Extension path mismatch")

	// Fetch branch node if it's a HashNode
	if hashNode, ok := decodedExt.Node.(*nodes2.HashNode); ok {
		decodedExt.Node, err = trie.DecodeNode(hashNode.Hash)
		assert.NoError(t, err, "Failed to decode hash node")
	}

	decodedBranch := decodedExt.Node.(*nodes2.BranchNode)

	// Fetch child nodes if they are HashNodes, then compare
	for i, child := range decodedBranch.Children {
		if child == nil {
			continue
		}
		if hashNode, ok := child.(*nodes2.HashNode); ok {
			decodedChild, err := trie.DecodeNode(hashNode.Hash)
			assert.NoError(t, err, fmt.Sprintf("Failed to decode child node at index %d", i))
			decodedBranch.Children[i] = decodedChild
		}
	}

	assert.Equal(t, leaf1.Path, decodedBranch.Children[1].(*nodes2.LeafNode).Path, "Leaf1 path mismatch")
	assert.Equal(t, leaf1.Value, decodedBranch.Children[1].(*nodes2.LeafNode).Value, "Leaf1 value mismatch")
	assert.Equal(t, leaf2.Path, decodedBranch.Children[2].(*nodes2.LeafNode).Path, "Leaf2 path mismatch")
	assert.Equal(t, leaf2.Value, decodedBranch.Children[2].(*nodes2.LeafNode).Value, "Leaf2 value mismatch")
}
