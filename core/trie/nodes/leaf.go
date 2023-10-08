package nodes

import "github.com/Aleksao998/Merkle-Patricia-Trie/core/trie/nibble"

type LeafNode struct {
	Path  []nibble.Nibble
	Value []byte
	Dirty bool
}

func NewLeafNode(path []nibble.Nibble, value []byte) *LeafNode {
	return &LeafNode{
		Path:  path,
		Value: value,
		Dirty: true,
	}
}

// IsDirty returns if the node was changed from last commit
func (l *LeafNode) IsDirty() bool {
	return l.Dirty
}
