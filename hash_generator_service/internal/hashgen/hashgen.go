package hashgen

import (
	"crypto/rand"
	"encoding/base64"
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

	return base64.StdEncoding.EncodeToString(data), nil
}
