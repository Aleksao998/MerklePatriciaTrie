package nibble

type Nibble byte

// FromBytes converts a byte slice to a slice of Nibbles.
func FromBytes(bytes []byte) []Nibble {
	nibbles := make([]Nibble, len(bytes)*2)

	for i, b := range bytes {
		high, low := b>>4, b&0x0F
		nibbles[i*2] = Nibble(high)
		nibbles[i*2+1] = Nibble(low)
	}

	return nibbles
}

// CommonPrefixLength returns the length of the common prefix shared between two slices of Nibbles
func CommonPrefixLength(node1 []Nibble, node2 []Nibble) int {
	minLen := len(node1)
	if len(node2) < minLen {
		minLen = len(node2)
	}

	for i := 0; i < minLen; i++ {
		if node1[i] != node2[i] {
			return i
		}
	}

	return minLen
}

// Equal returns if 2 nibble arrays are equal
func Equal(n1, n2 []Nibble) bool {
	if len(n1) != len(n2) {
		return false
	}

	for i, v := range n1 {
		if v != n2[i] {
			return false
		}
	}

	return true
}

// CompactEncoding adds a nibble prefix to a slice of nibbles to indicate its type and make its length even
func CompactEncoding(ns []Nibble, isLeafNode bool) []Nibble {
	var prefix []Nibble

	if isLeafNode {
		if len(ns)%2 == 0 {
			prefix = []Nibble{2, 0} // leaf node with even path length
		} else {
			prefix = []Nibble{3} // leaf node with odd path length
		}
	} else {
		if len(ns)%2 == 0 {
			prefix = []Nibble{0, 0} // extension node with even path length
		} else {
			prefix = []Nibble{1} // extension node with odd path length
		}
	}

	// Append prefix to the nibbles slice
	return append(prefix, ns...)
}

// ToBytes converts a slice of nibbles to a byte slice
func ToBytes(ns []Nibble) []byte {
	buf := make([]byte, 0, len(ns)/2)

	for i := 0; i < len(ns); i += 2 {
		b := byte(ns[i]<<4) + byte(ns[i+1])
		buf = append(buf, b)
	}

	return buf
}

// IsLeaf checks if the compact-encoded path is of a leaf node
func IsLeaf(encodedPath []Nibble) bool {
	firstNibble := encodedPath[0]

	return firstNibble == 2 || firstNibble == 3
}

// RemoveCompactEncoding removes the compact encoding from the nibble slice
func RemoveCompactEncoding(encodedPath []Nibble) []Nibble {
	if len(encodedPath) == 0 {
		return encodedPath
	}

	switch encodedPath[0] {
	case 0, 2:
		return encodedPath[2:]
	case 1, 3:
		return encodedPath[1:]
	default:
		return encodedPath // Invalid case, but just return as-is for this example
	}
}
