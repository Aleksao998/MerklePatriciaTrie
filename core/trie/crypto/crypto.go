package crypto

import (
	"golang.org/x/crypto/sha3"
	"log"
)

// Keccak256 calculates and returns the Keccak256 hash of the input data
func Keccak256(data []byte) []byte {
	d := sha3.NewLegacyKeccak256()

	_, err := d.Write(data)
	if err != nil {
		log.Fatalf("Failed to write data to hash: %v", err)
	}

	return d.Sum(nil)
}
