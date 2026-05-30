package a_crypto

import (
	"crypto/rand"
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
