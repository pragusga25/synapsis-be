package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateID generates a unique string ID of the specified length
func GenerateID(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}
