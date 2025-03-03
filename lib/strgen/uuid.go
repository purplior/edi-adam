package strgen

import (
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

func UniqueSortableID() (string, error) {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	id, err := ulid.New(ulid.Timestamp(time.Now()), entropy)

	return id.String(), err
}

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
