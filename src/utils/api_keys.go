package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateAPIKey() string {
	bytes := make([]byte, 32) // 256 bits
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
