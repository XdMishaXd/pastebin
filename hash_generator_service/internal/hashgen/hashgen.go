package hashgen

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
)

type HashGenerator struct {
	id      int
	hashLen int
}

// * New - конструктор для HashGenerator
func New(workerID int, hashLen int) *HashGenerator {
	return &HashGenerator{
		id:      workerID,
		hashLen: hashLen,
	}
}

// * Generate генерирует хэш с длинной hashLen
func (g *HashGenerator) Generate(hashLen int) (string, error) {
	data := make([]byte, hashLen)
	_, err := rand.Read(data)

	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	hash := sha512.Sum512(data)

	hexHash := hex.EncodeToString(hash[:])
	hexHash = hexHash[:hashLen]

	return hexHash, nil
}
