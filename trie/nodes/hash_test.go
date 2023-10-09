package nodes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewHashNode_Basic tests the basic creation of a HashNode with given hash.
func TestNewHashNode_Basic(t *testing.T) {
	t.Parallel()

	value := []byte{5, 6, 7, 8}

	leaf := NewHashNode(value)

	assert.Equal(t, value, leaf.Hash, "The values should be equal")
}
