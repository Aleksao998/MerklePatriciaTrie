package nodes

import "github.com/Aleksao998/Merkle-Patricia-Trie/core/trie/nibble"

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

// IsDirty returns if the node was changed from last commit
func (e *ExtensionNode) IsDirty() bool {
	return e.Dirty
}
