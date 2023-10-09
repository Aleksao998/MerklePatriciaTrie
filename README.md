<h1 align="center">🌴 MerklePatriciaTrie</h1>

MerklePatriciaTrie is an efficient and robust implementation of the trie data structure in Go. This trie is tailored for Ethereum-like systems but can be used in a variety of applications where data integrity, proof generation, and optimized storage are crucial.

## ⚠️ Disclaimer
1. **Not Production Ready:** This codebase is still in its development phase and should not be considered production-ready.

2. **Documentation in Progress:** While we strive to provide comprehensive documentation, it's not final. We appreciate your patience and any contributions to improve it.

3. **Optimization Pending:** The code, while functional, is yet to be optimized for performance. Future releases will focus on this aspect.

## 🌟 Features
### Efficient Storage Options:
   1. **PebbleDB:** A lightweight key-value store integrated for managing and preserving the data.
   2. **MPTMemoryStorage:** A custom in-memory storage solution, handy for generating proofs and extremely beneficial during unit testing.

### Comprehensive Operations:
Our trie supports various operations, like:
1. **Put:** To insert a key-value pair. 
2. **Get:** To retrieve the value for a given key. 
3. **Hash:** To calculate the hash of the entire trie. 
4. **Proof:** To generate a proof of inclusion for a specific key. 
5. **Commit:** To make all the changes permanent and return the root hash. 
6. **Del:** To delete a key-value pair from the trie.

## 📚 Resources & Documentation
[Official Documentation](https://gotolabs.gitbook.io/merklepatriciatrie/)
