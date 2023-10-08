package nodes

import (
	"github.com/Aleksao998/Merkle-Patricia-Trie/trie/nibble"
)

type ExtensionNode struct {
	Path  []nibble.Nibble
	Node  Node
	Dirty bool
}

func NewExtension(path []nibble.Nibble, node Node) *ExtensionNode {
	return &ExtensionNode{
		Path:  path,
		Node:  node,
		Dirty: true,
	}
}
