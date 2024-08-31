package strgen

import (
	"strings"

	"github.com/google/uuid"
)

func UniqueID() string {
	return uuid.New().String()
}

// Unique 보장이 떨어짐.
func ShortUniqueID() string {
	u := uuid.New()
	uuidStr := u.String()

	var result strings.Builder
	for _, char := range uuidStr {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') {
			result.WriteRune(char)
		}
	}

	shortUUID := result.String()[:12]

	return shortUUID
}
