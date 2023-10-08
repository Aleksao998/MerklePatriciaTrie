package crypto

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestKeccak256 tests hashing byte array
func TestKeccak256(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		value []byte
		hash  string
	}{
		{value: []byte("test1"), hash: "6d255fc3390ee6b41191da315958b7d6a1e5b17904cc7683558f98acc57977b4"},
		{value: []byte("randomString"), hash: "17ea9555c69c9fb4ab30d425cf5fe027a27562344ad65ca0d4ed63ceffafcb72"},
		{value: []byte("test1"), hash: "6d255fc3390ee6b41191da315958b7d6a1e5b17904cc7683558f98acc57977b4"},
		{value: []byte{}, hash: "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"},
	}

	for _, tt := range testCases {
		value := hex.EncodeToString(Keccak256(tt.value))
		require.Equal(t, tt.hash, value)
	}
}
