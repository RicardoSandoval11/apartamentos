package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func GenerateHash(input string) string {
	hasher := sha256.New()

	hasher.Write([]byte(input))

	hashBytes := hasher.Sum(nil)

	return hex.EncodeToString(hashBytes)
}
