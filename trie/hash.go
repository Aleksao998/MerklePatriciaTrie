package trie

import (
	"github.com/Aleksao998/Merkle-Patricia-Trie/trie/crypto"
	"github.com/Aleksao998/Merkle-Patricia-Trie/trie/nibble"
	nodes2 "github.com/Aleksao998/Merkle-Patricia-Trie/trie/nodes"
	"github.com/ethereum/go-ethereum/rlp"
)

func (t *Trie) NodeHash(node nodes2.Node) []byte {
	rlp, err := rlp.EncodeToBytes(t.NodeRaw(node, true))
	if err != nil {
		panic(err)
	}

	return crypto.Keccak256(rlp)
}

func (t *Trie) NodeRaw(node nodes2.Node, forHashing bool) interface{} {
	switch n := node.(type) {
	case nil:
		return []byte{}
	case *nodes2.LeafNode:
		return []interface{}{
			nibble.ToBytes(nibble.CompactEncoding(n.Path, true)),
			n.Value,
		}
	case *nodes2.ExtensionNode:
		nextData := t.NodeRaw(n.Node, forHashing)

		encodedNextData, _ := rlp.EncodeToBytes(nextData)
		if len(encodedNextData) >= 32 {
			nextData = t.NodeHash(n.Node)
		}

		return []interface{}{
			nibble.ToBytes(nibble.CompactEncoding(n.Path, false)),
			nextData,
		}
	case *nodes2.BranchNode:
		var childHashes [16]interface{}

		for i, child := range n.Children {
			if child != nil {
				childData := t.NodeRaw(child, forHashing)

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
	case *nodes2.HashNode:
		if !forHashing {
			return n.Hash // just return the hash if we're hashing
		}

		actualNode, err := t.DecodeNode(n.Hash)
		if err != nil {
			panic(err)
		}

		return t.NodeRaw(actualNode, forHashing)
	default:
		panic("Unknown node type")
	}
}
