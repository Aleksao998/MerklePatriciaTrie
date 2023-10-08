package nodes

// Package nodes represents the various types of nodes used in a Merkle Patricia Trie (MPT).
// In the context of an MPT, there are three primary types of nodes:
//
//  1. Leaf: A node that contains the end of a key and its associated value. Leaf nodes
//     signify the end of a path in the trie and store the final data. The key is used
//     to navigate through the trie, and when the end of the key is reached, the leaf's
//     value is the data associated with that key.
//
//  2. Extension: An intermediary node that has a shared key part and a single child.
//     Extension nodes are used to represent the shared parts of keys to save space.
//     They serve as a path compression mechanism to ensure the trie remains efficient
//     and doesn't grow unnecessarily large with redundant information.
//
//  3. Branch: A node that can have up to 16 children (based on hex values). This node type
//     is used when there's a divergence in paths. Each child represents a potential
//     continuation of the key, and traversal through the trie follows the path determined
//     by the key's hex value at that depth. Additionally, a Branch node can also signify
//     the end of a key if it holds a value directly, making it serve dual roles in the trie's
//     structure.
//
// Using these nodes, the MPT structures data in a way that allows for efficient
// lookups, updates, and deletions, while also ensuring cryptographic security and data
// integrity by means of hashing.
type Node interface{}
