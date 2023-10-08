package trie

import (
	"github.com/Aleksao998/Merkle-Patricia-Trie/core/trie/crypto"
	"github.com/Aleksao998/Merkle-Patricia-Trie/core/trie/nibble"
	"github.com/Aleksao998/Merkle-Patricia-Trie/core/trie/nodes"
	"github.com/ethereum/go-ethereum/rlp"
)

func (t *Trie) NodeHash(node nodes.Node) []byte {
	rlp, err := rlp.EncodeToBytes(t.NodeRaw(node))
	if err != nil {
		panic(err)
	}
	return crypto.Keccak256(rlp)
}

func (t *Trie) NodeRaw(node nodes.Node) interface{} {
	switch n := node.(type) {
	case nil:
		return []byte{}
	case *nodes.LeafNode:
		return []interface{}{
			nibble.ToBytes(nibble.CompactEncoding(n.Path, true)),
			n.Value,
		}
	case *nodes.ExtensionNode:
		nextData := t.NodeRaw(n.Node)
		encodedNextData, _ := rlp.EncodeToBytes(nextData)
		if len(encodedNextData) >= 32 {
			nextData = t.NodeHash(n.Node)
		}
		return []interface{}{
			nibble.ToBytes(nibble.CompactEncoding(n.Path, false)),
			nextData,
		}
	case *nodes.BranchNode:
		var childHashes [16]interface{}
		for i, child := range n.Children {
			if child != nil {
				childData := t.NodeRaw(child)
				encodedChildData, _ := rlp.EncodeToBytes(childData)
				if len(encodedChildData) >= 32 {
					childHashes[i] = t.NodeHash(child)
				} else {
					childHashes[i] = childData
				}
			} else {
				childHashes[i] = []byte{}
			}
		}
		return append(childHashes[:], n.Value)
	case *nodes.HashNode:
		actualNode, err := t.DecodeNode(n.Hash)
		if err != nil {
			panic(err)
		}
		return t.NodeRaw(actualNode) // Recursively get the raw data of the decoded node
	default:
		panic("Unknown node type")
	}
}
