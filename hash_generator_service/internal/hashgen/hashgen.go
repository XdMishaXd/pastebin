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

func New(workerID int, hashLen int) *HashGenerator {
	return &HashGenerator{
		id:      workerID,
		hashLen: hashLen,
	}
}

func (g *HashGenerator) Generate(hashLen int) (string, error) {
	data := make([]byte, hashLen)
	_, err := rand.Read(data)

	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return base64.StdEncoding.EncodeToString(data), nil
}
