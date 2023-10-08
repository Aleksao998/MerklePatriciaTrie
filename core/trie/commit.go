package trie

import (
	"fmt"
	"github.com/Aleksao998/Merkle-Patricia-Trie/core/trie/nibble"
	"github.com/Aleksao998/Merkle-Patricia-Trie/core/trie/nodes"
	"github.com/ethereum/go-ethereum/rlp"
)

const rootHashKey = "rootHash"

func (t *Trie) commit(node nodes.Node) ([]byte, error) {
	switch n := node.(type) {
	case *nodes.LeafNode:
		return t.handleLeafNode(n)
	case *nodes.ExtensionNode:
		return t.handleExtensionNode(n)
	case *nodes.BranchNode:
		return t.handleBranchNode(n)
	case *nodes.HashNode:
		return n.Hash, nil
	default:
		panic("Unknown node type")
	}
}

func (t *Trie) handleLeafNode(n *nodes.LeafNode) ([]byte, error) {
	raw := t.NodeRaw(n)
	encoded, err := rlp.EncodeToBytes(raw)
	if err != nil {
		return nil, err
	}

	hash := t.NodeHash(n)
	err = t.storage.Put(hash, encoded)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func (t *Trie) handleExtensionNode(n *nodes.ExtensionNode) ([]byte, error) {
	childHash, err := t.commit(n.Node)
	if err != nil {
		return nil, err
	}

	// replace the node with its hash node
	n.Node = nodes.NewHashNode(childHash)

	raw := t.NodeRaw(n)
	encoded, err := rlp.EncodeToBytes(raw)
	if err != nil {
		return nil, err
	}

	hash := t.NodeHash(n)
	err = t.storage.Put(hash, encoded)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func (t *Trie) handleBranchNode(n *nodes.BranchNode) ([]byte, error) {
	for index, child := range n.Children {
		if child != nil {
			childHash, err := t.commit(child)
			if err != nil {
				return nil, err
			}
			n.Children[index] = nodes.NewHashNode(childHash)
		}
	}

	raw := t.NodeRaw(n)
	encoded, err := rlp.EncodeToBytes(raw)
	if err != nil {
		return nil, err
	}

	hash := t.NodeHash(n)
	err = t.storage.Put(hash, encoded)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func (t *Trie) DecodeNode(hash []byte) (nodes.Node, error) {
	data, err := t.storage.Get(hash)
	if err != nil {
		return nil, err
	}

	raw := []interface{}{}
	if err := rlp.DecodeBytes(data, &raw); err != nil {
		return nil, err
	}

	return t.reconstructNode(raw)
}

func (t *Trie) reconstructNode(raw []interface{}) (nodes.Node, error) {
	switch len(raw) {
	case 2: // Could be LeafNode or ExtensionNode
		path := nibble.FromBytes(raw[0].([]byte))
		isLeaf := nibble.IsLeaf(path)
		path = nibble.RemoveCompactEncoding(path)

		if isLeaf {
			return &nodes.LeafNode{
				Path:  path,
				Value: raw[1].([]byte),
				Dirty: false,
			}, nil
		}

		// Handle ExtensionNode's child
		child, err := t.decodeChild(raw[1])
		if err != nil {
			return nil, err
		}

		return &nodes.ExtensionNode{
			Path: path,
			Node: child,
		}, nil

	case 17: // BranchNode
		branch := &nodes.BranchNode{Dirty: false}
		for i := 0; i < 16; i++ {
			child, err := t.decodeChild(raw[i])
			if err != nil {
				return nil, err
			}
			branch.Children[i] = child
		}
		branch.Value = raw[16].([]byte)
		return branch, nil

	default:
		return nil, fmt.Errorf("Unknown node type")
	}
}

func (t *Trie) decodeChild(data interface{}) (nodes.Node, error) {
	switch v := data.(type) {
	case []byte:
		if len(v) == 32 { // hash length
			return &nodes.HashNode{Hash: v}, nil
		}
		return nil, nil
	case []interface{}:
		return t.reconstructNode(v)
	default:
		return nil, fmt.Errorf("Unexpected child data type")
	}
}

// SetRootHash saves the root hash in the Committer and also in the key-value storage.
func (t *Trie) SetRootHash(hash []byte) error {
	// Update in-memory representation
	t.rootHash = hash

	// Persist to the key-value storage
	err := t.storage.Put([]byte(rootHashKey), hash)
	if err != nil {
		return fmt.Errorf("failed to set root hash in storage: %v", err)
	}

	return nil
}

// GetRootHash retrieves the root hash from the Committer. If it's not present in memory,
// it tries to fetch from the key-value storage.
func (t *Trie) GetRootHash() ([]byte, error) {
	if t.rootHash != nil {
		return t.rootHash, nil
	}

	// If the rootHash is nil, try fetching from the key-value storage
	value, err := t.storage.Get([]byte(rootHashKey))
	if err != nil {
		return nil, fmt.Errorf("failed to get root hash from storage: %v", err)
	}

	// Update in-memory representation
	t.rootHash = value

	return value, nil
}
