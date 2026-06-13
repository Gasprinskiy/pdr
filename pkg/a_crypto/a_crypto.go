package a_crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

type ACrypto struct {
	IDComplexity int
}

func NewACrypto(idComplexity int) *ACrypto {
	return &ACrypto{
		IDComplexity: idComplexity,
	}
}

func (c *ACrypto) GenerateHexID() string {
	b := make([]byte, c.IDComplexity)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (c *ACrypto) GenerateHashFromBytes(b []byte) string {
	hash := sha256.Sum256(b)
	return hex.EncodeToString(hash[:])
}
