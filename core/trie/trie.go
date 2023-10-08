package trie

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Aleksao998/Merkle-Patricia-Trie/core/storage"
	"github.com/Aleksao998/Merkle-Patricia-Trie/core/trie/nibble"
	"github.com/Aleksao998/Merkle-Patricia-Trie/core/trie/nodes"
	"sync"
)

var (
	errKeyNotFound = errors.New("key not found")
)

type Trie struct {
	root     nodes.Node
	storage  storage.Storage
	mu       sync.RWMutex
	rootHash []byte
}

func NewTrie(storage storage.Storage) *Trie {
	return &Trie{
		storage: storage,
	}
}

func (t *Trie) Hash() []byte {
	t.getRootHash()

	if t.root == nil {
		return nil // Empty Trie
	}
	return t.NodeHash(t.root)
}

// Proof returns the Merkle-proof associated with
// a node. An error is returned if the node is not found.
func (t *Trie) Proof(key []byte) (storage.Storage, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.getRootHash()

	return t.GenerateProof(t.root, key)
}

// Get retrieves the value associated with a given key in the trie
func (t *Trie) Get(key []byte) ([]byte, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	// convert the byte key to a nibble path for easier traversal
	nibblePath := nibble.FromBytes(key)

	currentNode := &t.root

	t.getRootHash()

	// loop until a value is found, or it's determined the key is not in the trie
	for {
		switch node := (*currentNode).(type) {
		case nil:
			// if a nil node is encountered, the key isn't in the trie
			return nil, errKeyNotFound
		case *nodes.HashNode:
			// If a HashNode is encountered, fetch the actual node from storage
			actualNode, err := t.DecodeNode(node.Hash)
			if err != nil {
				return nil, err
			}
			*currentNode = actualNode
			continue
		case *nodes.LeafNode:
			// calculate the length of the common prefix with the leaf node's path
			commonLength := nibble.CommonPrefixLength(node.Path, nibblePath)
			if commonLength != len(node.Path) || commonLength != len(nibblePath) {
				// if they don't match exactly, the key isn't in the trie
				return nil, errKeyNotFound
			}
			// if they do match, return the leaf node's value
			return node.Value, nil
		case *nodes.BranchNode:
			// if there's no remaining path, try to get the value directly from the branch node
			if len(nibblePath) == 0 {
				value, found := node.GetValue()
				if found {
					return value, nil
				}
				return nil, errKeyNotFound
			}

			// otherwise, extract the next child nibble and the remaining path
			child, remaining := nibblePath[0], nibblePath[1:]
			nibblePath = remaining
			// move to the child node based on the nibble
			currentNode = &node.Children[child]

			// continue the loop with the child node
			continue
		case *nodes.ExtensionNode:
			// calculate the length of the common prefix with the extension node's path
			commonLength := nibble.CommonPrefixLength(node.Path, nibblePath)
			if commonLength < len(node.Path) {
				// if they don't share the full extension path, the key isn't in the trie
				return nil, errKeyNotFound
			}

			// move to the next segment of the nibble path
			nibblePath = nibblePath[commonLength:]
			// move to the child node of the extension
			currentNode = &node.Node
		default:
			// if an unexpected node type is encountered, panic
			panic("Unexpected node type encountered while traversing the trie")
		}
	}
}

// Put inserts or updates a value associated with a given key in the trie
func (t *Trie) Put(key []byte, value []byte) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// convert the byte key to a nibble path for easier traversal
	nibblePath := nibble.FromBytes(key)
	currentNode := &t.root

	t.getRootHash()

	// loop until a value is set or updated
	for {
		switch node := (*currentNode).(type) {
		case nil:
			// if current node is nil, create a new leaf node with the remaining nibble path and value
			*currentNode = nodes.NewLeafNode(nibblePath, value)
			return

		case *nodes.LeafNode:
			// handle the logic of inserting a key-value pair when encountering a leaf node
			t.handleLeafNodeInsert(currentNode, node, nibblePath, value)
			return

		case *nodes.BranchNode:
			node.Dirty = true

			// if there's no remaining path, set the value directly on the branch node
			if len(nibblePath) == 0 {
				node.SetValue(value)
				return
			}
			// update the current node to the child pointed by the next nibble and continue to next segment
			currentNode = &node.Children[nibblePath[0]]
			nibblePath = nibblePath[1:]

		case *nodes.ExtensionNode:
			node.Dirty = true

			// calculate the length of the common prefix with the extension node's path
			commonLength := nibble.CommonPrefixLength(node.Path, nibblePath)
			// if they don't share the full extension path, handle the logic of inserting in such scenario
			if commonLength < len(node.Path) {
				t.handleExtensionNodeInsert(currentNode, node, nibblePath, value, commonLength)
				return
			}
			// move to the next segment of the nibble path and the child node of the extension
			nibblePath = nibblePath[commonLength:]
			currentNode = &node.Node
		case *nodes.HashNode:
			actualNode, err := t.DecodeNode(node.Hash)
			if err != nil {
				panic(err.Error())
			}

			*currentNode = actualNode
		default:
			// if an unexpected node type is encountered, panic
			panic("Unexpected node type encountered while traversing the trie")
		}
	}
}

// Commit saves the trie in persistent storage
// and returns the trie root key.
func (t *Trie) Commit() []byte {
	t.mu.Lock()
	defer t.mu.Unlock()

	rootKey, err := t.commit(t.root)
	if err != nil {
		panic("Failed to commit the trie: " + err.Error())
	}

	if err := t.SetRootHash(rootKey); err != nil {
		panic("Failed to set root hash: " + err.Error())
	}

	// Set the root to nil to release the in-memory storage of the trie.
	t.root = nil

	return rootKey
}

// Del removes the key from the trie
func (t *Trie) Del(key []byte) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// convert the byte key to a nibble path for easier traversal
	nibblePath := nibble.FromBytes(key)
	// use pathStack to keep track of nodes for potential path compression later
	var pathStack []*nodes.Node
	currentNode := &t.root

	t.getRootHash()

	// loop until the key is found and removed or until it's clear the key doesn't exist
	for {
		switch node := (*currentNode).(type) {
		case nil:
			// if the currentNode is nil, the key is not in the trie
			return errKeyNotFound
		case *nodes.HashNode:
			actualNode, err := t.DecodeNode(node.Hash)
			if err != nil {
				panic(err.Error())
			}

			*currentNode = actualNode
		case *nodes.LeafNode:
			// if the key matches with the leaf node's path, delete the leaf node
			if nibble.Equal(node.Path, nibblePath) {
				*currentNode = nil
				t.compressPath(pathStack)
				return nil
			}
			return errKeyNotFound
		case *nodes.BranchNode:
			// if there's no remaining path and the branch node has the value, delete the value
			if len(nibblePath) == 0 {
				if !node.HasValue() {
					return errKeyNotFound
				}

				node.ClearValue()
				if node.ChildCount() == 1 {
					t.compressBranchNode(node, currentNode)
				}

				node.Dirty = true
				t.compressPath(pathStack)
				return nil
			}
			// update the current node and path and keep track of the nodes encountered
			pathStack = append(pathStack, currentNode)
			childNibble := nibblePath[0]
			nibblePath = nibblePath[1:]
			currentNode = &node.Children[childNibble]
		case *nodes.ExtensionNode:
			// if the key doesn't share the full extension path, return key not found
			commonLength := nibble.CommonPrefixLength(node.Path, nibblePath)
			if commonLength < len(node.Path) {
				return errKeyNotFound
			}

			// update the current node and path and keep track of the nodes encountered
			pathStack = append(pathStack, currentNode)
			nibblePath = nibblePath[commonLength:]
			currentNode = &node.Node
		default:
			// if an unexpected node type is encountered, panic
			panic("Unexpected node type encountered while traversing the trie")
		}
	}
}

// compressPath compresses the path after deletion if possible
func (t *Trie) compressPath(pathStack []*nodes.Node) {
	for len(pathStack) > 0 {
		node := pathStack[len(pathStack)-1]

		switch n := (*node).(type) {
		case *nodes.BranchNode:
			// compress the branch node if it has only one child left and no value
			if n.ChildCount() == 1 && !n.HasValue() {
				t.compressBranchNode(n, node)
			} else {
				break
			}

		case *nodes.ExtensionNode:
			// merge consecutive extension nodes if found
			childNode := n.Node
			switch child := childNode.(type) {
			case *nodes.ExtensionNode:
				mergedPath := append(n.Path, child.Path...)
				newNode := nodes.NewExtension(mergedPath, child.Node)
				*node = newNode
			case *nodes.LeafNode:
				mergedPath := append(n.Path, child.Path...)
				newNode := nodes.NewLeafNode(mergedPath, child.Value)
				*node = newNode
			case *nodes.BranchNode:
				if child.ChildCount() == 0 && child.HasValue() {
					newNode := nodes.NewLeafNode(n.Path, child.Value)
					*node = newNode
				}
			default:
				break
			}
		}

		pathStack = pathStack[:len(pathStack)-1]
	}
}

// compressBranchNode compresses a branch node into a leaf or extension node
// This happens when a branch node has only one child. Instead of keeping
// the branch node structure, the trie can be made more efficient by
// compressing the branch node
func (t *Trie) compressBranchNode(node *nodes.BranchNode, parentNode *nodes.Node) {
	// iterate over the children of the branch node
	for i, child := range node.Children {
		if child != nil {
			// check the type of the child node
			switch c := child.(type) {
			case *nodes.LeafNode:
				// if the child is a leaf node, merge the branch node and leaf node paths
				mergedPath := append([]nibble.Nibble{nibble.Nibble(i)}, c.Path...)
				*parentNode = nodes.NewLeafNode(mergedPath, c.Value)
			default:
				// if the child is any other type, create a new extension node with the child
				*parentNode = nodes.NewExtension([]nibble.Nibble{nibble.Nibble(i)}, child)
			}
			// exit the loop once the compression is done for the only child
			break
		}
	}
}

// getRootHash gets root hash from storage if nil and if available
func (t *Trie) getRootHash() {
	// If root is nil, attempt to fetch root node from storage
	if t.root == nil {
		rootHash, _ := t.GetRootHash()

		// If rootHash is nil, it indicates an empty trie and we can just return
		if rootHash == nil {
			return
		}

		rootNode, err := t.DecodeNode(rootHash)
		if err != nil {
			panic(err.Error())
		}
		t.root = rootNode
	}
}

// handleLeafNodeInsert handles the insertion logic when encountering a leaf node in the trie
func (t *Trie) handleLeafNodeInsert(currentNode *nodes.Node, leafNode *nodes.LeafNode, nibblePath []nibble.Nibble, value []byte) {
	// calculate the length of the common prefix between the leaf node path and the input nibble path
	commonLength := nibble.CommonPrefixLength(leafNode.Path, nibblePath)

	// check if the leaf node's path is the same as the input nibble path
	if commonLength == len(nibblePath) && commonLength == len(leafNode.Path) && !bytes.Equal(leafNode.Value, value) {
		// if they're the same and the values differ, update the current node to the new value
		*currentNode = nodes.NewLeafNode(nibblePath, value)
		return
	}

	// create a new branch node to handle the split
	branchNode := nodes.NewBranchNode()

	// setBranch function sets the child or value of the branch node based on the path
	setBranch := func(path []nibble.Nibble, val []byte) {
		// if there are more nibbles beyond the common length, add a new leaf node as a child
		if commonLength < len(path) {
			branchNibble, leafNibbles := path[commonLength], path[commonLength+1:]
			newLeaf := nodes.NewLeafNode(leafNibbles, val)
			branchNode.SetChild(branchNibble, newLeaf)
		} else {
			// if the path is entirely common, set the value of the branch node
			branchNode.SetValue(val)
		}
	}

	// set values or children for both the leaf node and the input path
	setBranch(leafNode.Path, leafNode.Value)
	setBranch(nibblePath, value)

	// if there's a common prefix, create an extension node with the branch node as a child
	if commonLength > 0 {
		*currentNode = nodes.NewExtension(nibblePath[:commonLength], branchNode)
	} else {
		// if no common prefix, set the current node to the branch node directly
		*currentNode = branchNode
	}
}

// handleExtensionNodeInsert handles the insertion logic when encountering an extension node in the trie
func (t *Trie) handleExtensionNodeInsert(currentNode *nodes.Node, extNode *nodes.ExtensionNode, nibblePath []nibble.Nibble, value []byte, commonLength int) {
	// derive the common nibbles, the branching nibble, and the remaining nibbles from the extension node's path
	extNibbles := extNode.Path[:commonLength]
	branchNibble := extNode.Path[commonLength]
	extRemainingNibbles := extNode.Path[commonLength+1:]

	// create a new branch node to handle the split
	branchNode := nodes.NewBranchNode()

	// if the extension node doesn't have any remaining nibbles after the common ones
	if len(extRemainingNibbles) == 0 {
		branchNode.SetChild(branchNibble, extNode.Node)
	} else {
		// if there are remaining nibbles, create a new extension node and set it as a child of the branch node
		newExtension := nodes.NewExtension(extRemainingNibbles, extNode.Node)
		branchNode.SetChild(branchNibble, newExtension)
	}

	// check the length of the given nibble path in comparison to the common length
	if commonLength < len(nibblePath) {
		nodeBranchNibble, nodeLeafNibbles := nibblePath[commonLength], nibblePath[commonLength+1:]
		remainingLeaf := nodes.NewLeafNode(nodeLeafNibbles, value)
		branchNode.SetChild(nodeBranchNibble, remainingLeaf)
	} else if commonLength == len(nibblePath) {
		// if they are the same length, set the value on the branch node
		branchNode.SetValue(value)
	} else {
		// if there's an unexpected match of more nibbles than provided, panic
		panic(fmt.Sprintf("too many matched (%v > %v)", commonLength, len(nibblePath)))
	}

	// if there are no common nibbles, set the current node to the branch node
	if len(extNibbles) == 0 {
		*currentNode = branchNode
	} else {
		// if there are common nibbles, create an extension node with the branch node as a child
		*currentNode = nodes.NewExtension(extNibbles, branchNode)
	}
}
