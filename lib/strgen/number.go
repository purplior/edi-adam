package strgen

import (
	"math/rand"
	"strings"
	"time"
)

func RandomNumber(length int) string {
	digits := "0123456789"

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(digits[r.Intn(len(digits))])
	}

	return sb.String()
}
