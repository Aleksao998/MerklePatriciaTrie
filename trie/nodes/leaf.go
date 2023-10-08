package nodes

import (
	"github.com/Aleksao998/Merkle-Patricia-Trie/trie/nibble"
)

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
