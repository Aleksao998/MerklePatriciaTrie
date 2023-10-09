package trie

import (
	"fmt"

	"github.com/Aleksao998/Merkle-Patricia-Trie/trie/nibble"
	nodes2 "github.com/Aleksao998/Merkle-Patricia-Trie/trie/nodes"
	"github.com/ethereum/go-ethereum/rlp"
)

const rootHashKey = "rootHash"

func (t *Trie) commit(node nodes2.Node) ([]byte, error) {
	switch n := node.(type) {
	case *nodes2.LeafNode:
		return t.handleLeafNode(n)
	case *nodes2.ExtensionNode:
		return t.handleExtensionNode(n)
	case *nodes2.BranchNode:
		return t.handleBranchNode(n)
	case *nodes2.HashNode:
		return n.Hash, nil
	default:
		panic("Unknown node type")
	}
}

func (t *Trie) handleLeafNode(n *nodes2.LeafNode) ([]byte, error) {
	raw := t.NodeRaw(n, false)

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

func (t *Trie) handleExtensionNode(n *nodes2.ExtensionNode) ([]byte, error) {
	childHash, err := t.commit(n.Node)
	if err != nil {
		return nil, err
	}

	// replace the node with its hash node
	n.Node = nodes2.NewHashNode(childHash)

	raw := t.NodeRaw(n, false)

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

func (t *Trie) handleBranchNode(n *nodes2.BranchNode) ([]byte, error) {
	for index, child := range n.Children {
		if child != nil {
			childHash, err := t.commit(child)
			if err != nil {
				return nil, err
			}

			n.Children[index] = nodes2.NewHashNode(childHash)
		}
	}

	raw := t.NodeRaw(n, false)

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

func (t *Trie) DecodeNode(hash []byte) (nodes2.Node, error) {
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

func (t *Trie) reconstructNode(raw []interface{}) (nodes2.Node, error) {
	switch len(raw) {
	case 2: // Could be LeafNode or ExtensionNode
		pathBytes, ok := raw[0].([]byte)
		if !ok {
			// Handle the error. For example:
			return nil, fmt.Errorf("expected raw[0] to be []byte, got %T", raw[0])
		}

		path := nibble.FromBytes(pathBytes)

		isLeaf := nibble.IsLeaf(path)
		path = nibble.RemoveCompactEncoding(path)

		valueBytes, ok := raw[1].([]byte)
		if !ok {
			return nil, fmt.Errorf("expected raw[1] to be []byte, got %T", raw[1])
		}

		if isLeaf {
			return &nodes2.LeafNode{
				Path:  path,
				Value: valueBytes,
				Dirty: false,
			}, nil
		}

		// Handle ExtensionNode's child
		child, err := t.decodeChild(raw[1])
		if err != nil {
			return nil, err
		}

		return &nodes2.ExtensionNode{
			Path: path,
			Node: child,
		}, nil

	case 17: // BranchNode
		branch := &nodes2.BranchNode{Dirty: false}

		for i := 0; i < 16; i++ {
			child, err := t.decodeChild(raw[i])
			if err != nil {
				return nil, err
			}

			branch.Children[i] = child
		}

		branchBytes, ok := raw[16].([]byte)
		if !ok {
			// Handle the error. For example:
			return nil, fmt.Errorf("expected raw[16] to be []byte, got %T", raw[16])
		}

		branch.Value = branchBytes

		return branch, nil

	default:
		return nil, fmt.Errorf("nknown node type")
	}
}

func (t *Trie) decodeChild(data interface{}) (nodes2.Node, error) {
	switch v := data.(type) {
	case []byte:
		if len(v) == 32 { // hash length
			return &nodes2.HashNode{Hash: v}, nil
		}

		return nil, nil
	default:
		return nil, fmt.Errorf("unexpected child data type")
	}
}

// SetRootHash saves the root hash in the Committer and also in the key-value storage.
func (t *Trie) SetRootHash(hash []byte) error {
	// Update in-memory representation
	t.rootHash = hash

	// Persist to the key-value storage
	err := t.storage.Put([]byte(rootHashKey), hash)
	if err != nil {
		return fmt.Errorf("failed to set root hash in storage: %w", err)
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
		return nil, fmt.Errorf("failed to get root hash from storage: %w", err)
	}

	// Update in-memory representation
	t.rootHash = value

	return value, nil
}
