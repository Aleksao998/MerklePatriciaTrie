package trie

import (
	"fmt"
	"github.com/Aleksao998/Merkle-Patricia-Trie/storage"
	"github.com/Aleksao998/Merkle-Patricia-Trie/storage/mpt"
	"github.com/Aleksao998/Merkle-Patricia-Trie/trie/nibble"
	nodes2 "github.com/Aleksao998/Merkle-Patricia-Trie/trie/nodes"
	"github.com/ethereum/go-ethereum/rlp"
)

func (t *Trie) GenerateProof(root nodes2.Node, key []byte) (storage.Storage, error) {
	currentNode := root
	nibblePath := nibble.FromBytes(key)
	db := mpt.NewMPTMemoryStorage()

	for {
		switch node := currentNode.(type) {
		case nil:
			// If node is nil, then the path does not exist in the trie
			t.storeNode(db, currentNode)
			return db, errKeyNotFound

		case *nodes2.LeafNode:
			t.storeNode(db, node)
			if nibble.Equal(node.Path, nibblePath) {
				// Key found in trie
				return db, nil
			}
			// Path mismatch
			return db, errKeyNotFound

		case *nodes2.BranchNode:
			t.storeNode(db, node)
			if len(nibblePath) == 0 {
				_, found := node.GetValue()
				if found {
					return db, nil
				}
				return db, errKeyNotFound
			}
			// Move to the next node in the branch
			currentNode = node.Children[nibblePath[0]]
			nibblePath = nibblePath[1:]
			continue

		case *nodes2.ExtensionNode:
			t.storeNode(db, node)
			matchLen := nibble.CommonPrefixLength(node.Path, nibblePath)
			if matchLen < len(node.Path) {
				return db, fmt.Errorf("Key not found in trie")
			}
			nibblePath = nibblePath[matchLen:]
			currentNode = node.Node
			continue
		case *nodes2.HashNode:
			actualNode, err := t.DecodeNode(node.Hash)
			if err != nil {
				return db, err
			}
			currentNode = actualNode
			t.storeNode(db, currentNode)
			continue
		default:
			panic("Unexpected node type encountered while traversing the trie")
		}
	}

	return db, fmt.Errorf("Key not found in trie")
}

func (t *Trie) storeNode(db storage.Storage, node nodes2.Node) {
	rawNode := t.NodeRaw(node)
	encoded, err := rlp.EncodeToBytes(rawNode)
	if err != nil {
		panic(err)
	}
	db.Put(t.NodeHash(node), encoded)
}
