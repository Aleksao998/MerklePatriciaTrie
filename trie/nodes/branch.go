package nodes

import (
	"github.com/Aleksao998/Merkle-Patricia-Trie/trie/nibble"
)

const BranchChildrenSize = 16

type BranchNode struct {
	Children   [BranchChildrenSize]Node
	Value      []byte
	childCount int
	Dirty      bool
}

func NewBranchNode() *BranchNode {
	return &BranchNode{
		Children: [BranchChildrenSize]Node{},
		Dirty:    true,
	}
}

// SetValue sets value for branch node
func (b *BranchNode) SetValue(value []byte) {
	b.Dirty = true
	b.Value = value
}

// GetValue returns the value of the branch node and a boolean indicating its existence
func (b *BranchNode) GetValue() (value []byte, exists bool) {
	if b.Value != nil {
		return b.Value, true
	}
	return nil, false
}

// SetChild sets a specific child node in the branch node based on the given nibble
func (b *BranchNode) SetChild(nibble nibble.Nibble, node Node) {
	if int(nibble) < BranchChildrenSize {
		// If the child is being set to nil, and it previously existed, decrement childCount
		b.Dirty = true
		b.Children[int(nibble)] = node
	} else {
		panic("Invalid nibble for BranchNode")
	}
}

// ChildCount returns the number of children the branch node has (i.e., non-nil children)
func (b *BranchNode) ChildCount() int {
	total := 0
	for _, child := range b.Children {
		if child != nil {
			total++
		}
	}
	return total
}

// HasValue checks if the branch node has a value
func (b *BranchNode) HasValue() bool {
	return b.Value != nil
}

// ClearValue clears value
func (b *BranchNode) ClearValue() {
	b.Dirty = true
	b.Value = nil
}

// IsDirty returns if the node was changed from last commit
func (b *BranchNode) IsDirty() bool {
	return b.Dirty
}
