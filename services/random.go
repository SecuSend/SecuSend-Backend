package services

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

func GenerateUniqueKey() (string, error) {
	// Generate 16 bytes of random data
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	// Append the current Unix timestamp in nanoseconds
	timestampBytes := make([]byte, 8)
	timestamp := time.Now().UnixNano()
	for i := 0; i < 8; i++ {
		timestampBytes[i] = byte(timestamp >> (i * 8))
	}

	keyBytes := append(randomBytes, timestampBytes...)
	return hex.EncodeToString(keyBytes), nil
}
