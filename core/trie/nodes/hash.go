package nodes

type HashNode struct {
	Hash []byte
}

func NewHashNode(hash []byte) *HashNode {
	return &HashNode{
		Hash: hash,
	}
}
