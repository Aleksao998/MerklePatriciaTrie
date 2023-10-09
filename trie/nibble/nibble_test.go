package nibble

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFromBytes tests converting byte array to nibble array
func TestFromBytes(t *testing.T) {
	tests := []struct {
		name    string
		bytes   []byte
		nibbles []Nibble
	}{
		{"Byte array", []byte{0xAB, 0xCD}, []Nibble{0xA, 0xB, 0xC, 0xD}},
		{"Single item in byte array", []byte{0x12}, []Nibble{0x1, 0x2}},
		{"Empty byte array", []byte{}, []Nibble{}},
		{"Byte with leading zeros", []byte{0x01, 0x05}, []Nibble{0x0, 0x1, 0x0, 0x5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nibbles := FromBytes(tt.bytes)
			for i, nibble := range nibbles {
				assert.Equal(t, tt.nibbles[i], nibble)
			}
		})
	}
}

// TestCommonPrefixLength tests the length of the common prefix shared between two slices of Nibbles
func TestCommonPrefixLength(t *testing.T) {
	tests := []struct {
		name           string
		node1          []Nibble
		node2          []Nibble
		expectedLength int
	}{
		{"Matching nibbles", []Nibble{0xA, 0xB, 0xC, 0xD}, []Nibble{0xA, 0xB, 0xC}, 3},
		{"Mismatch at start", []Nibble{0x1, 0x2, 0x3}, []Nibble{0x4, 0x5, 0x6}, 0},
		{"Empty slice", []Nibble{}, []Nibble{0x7, 0x8, 0x9}, 0},
		{"First nibble array longer", []Nibble{0xA, 0xB, 0xC, 0xD}, []Nibble{0xA, 0xB}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			length := CommonPrefixLength(tt.node1, tt.node2)
			assert.Equal(t, tt.expectedLength, length)
		})
	}
}

// TestEqual tests if two slices of Nibbles are equal
func TestEqual(t *testing.T) {
	tests := []struct {
		name     string
		node1    []Nibble
		node2    []Nibble
		expected bool
	}{
		{"Equal nibbles", []Nibble{0xA, 0xB, 0xC, 0xD}, []Nibble{0xA, 0xB, 0xC, 0xD}, true},
		{"Different lengths", []Nibble{0xA, 0xB, 0xC, 0xD}, []Nibble{0xA, 0xB, 0xC}, false},
		{"Mismatch at start", []Nibble{0x1, 0x2, 0x3}, []Nibble{0x4, 0x5, 0x6}, false},
		{"Empty slice vs non-empty", []Nibble{}, []Nibble{0x7, 0x8, 0x9}, false},
		{"Both empty slices", []Nibble{}, []Nibble{}, true},
		{"Length mismatch", []Nibble{0xA, 0xB}, []Nibble{0xA, 0xB, 0xC}, false},
		{"Mismatch in middle", []Nibble{0xA, 0xB, 0xC}, []Nibble{0xA, 0x1, 0xC}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Equal(tt.node1, tt.node2)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestCompactEncoding tests the CompactEncoding function
func TestCompactEncoding(t *testing.T) {
	tests := []struct {
		name       string
		ns         []Nibble
		isLeafNode bool
		expected   []Nibble
	}{
		{"Leaf even length", []Nibble{1, 2, 3, 4}, true, []Nibble{2, 0, 1, 2, 3, 4}},
		{"Leaf odd length", []Nibble{1, 2, 3}, true, []Nibble{3, 1, 2, 3}},
		{"Extension even length", []Nibble{1, 2, 3, 4}, false, []Nibble{0, 0, 1, 2, 3, 4}},
		{"Extension odd length", []Nibble{1, 2, 3}, false, []Nibble{1, 1, 2, 3}},
		{"Leaf node with even path length", []Nibble{1, 2, 3, 4}, true, []Nibble{2, 0, 1, 2, 3, 4}},
		{"Extension node with even path length", []Nibble{1, 2, 3, 4}, false, []Nibble{0, 0, 1, 2, 3, 4}},
		{"Leaf node with even path length", []Nibble{0xA, 0xB, 0xC, 0xD}, true, []Nibble{2, 0, 0xA, 0xB, 0xC, 0xD}},
		{"Extension node with odd path length", []Nibble{0xB}, false, []Nibble{1, 0xB}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompactEncoding(tt.ns, tt.isLeafNode)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestToBytes tests the ToBytes function
func TestToBytes(t *testing.T) {
	tests := []struct {
		name     string
		ns       []Nibble
		expected []byte
	}{
		{"Convert nibbles to bytes", []Nibble{0xA, 0xB, 0xC, 0xD}, []byte{0xAB, 0xCD}},
		{"Empty nibbles", []Nibble{}, []byte{}},
		{"Nibbles with leading zeros", []Nibble{0x0, 0x1, 0x0, 0x5}, []byte{0x01, 0x05}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToBytes(tt.ns)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestIsLeaf tests the IsLeaf function
func TestIsLeaf(t *testing.T) {
	tests := []struct {
		name     string
		ns       []Nibble
		expected bool
	}{
		{"Leaf with even path", []Nibble{2, 0, 1, 2, 3, 4}, true},
		{"Leaf with odd path", []Nibble{3, 1, 2, 3}, true},
		{"Extension with even path", []Nibble{0, 0, 1, 2, 3, 4}, false},
		{"Extension with odd path", []Nibble{1, 1, 2, 3}, false},
		{"Invalid prefix", []Nibble{4, 5, 6, 7}, false}, // for this example
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsLeaf(tt.ns)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestRemoveCompactEncoding tests the RemoveCompactEncoding function
func TestRemoveCompactEncoding(t *testing.T) {
	tests := []struct {
		name     string
		ns       []Nibble
		expected []Nibble
	}{
		{"Leaf with even path", []Nibble{2, 0, 1, 2, 3, 4}, []Nibble{1, 2, 3, 4}},
		{"Leaf with odd path", []Nibble{3, 1, 2, 3}, []Nibble{1, 2, 3}},
		{"Extension with even path", []Nibble{0, 0, 1, 2, 3, 4}, []Nibble{1, 2, 3, 4}},
		{"Extension with odd path", []Nibble{1, 1, 2, 3}, []Nibble{1, 2, 3}},
		{"Empty encoded path", []Nibble{}, []Nibble{}},
		{"Encoded path with prefix 0", []Nibble{0, 0, 1, 2, 3}, []Nibble{1, 2, 3}},
		{"Empty encoded path", []Nibble{}, []Nibble{}},
		{"Encoded path with prefix 1", []Nibble{1, 0xA, 0xB}, []Nibble{0xA, 0xB}},
		{"Encoded path with prefix 3", []Nibble{3, 0xA, 0xB, 0xC}, []Nibble{0xA, 0xB, 0xC}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveCompactEncoding(tt.ns)
			assert.Equal(t, tt.expected, result)
		})
	}
}
