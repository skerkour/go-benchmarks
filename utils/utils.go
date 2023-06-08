package utils

import (
	"crypto/rand"
	"fmt"
	"testing"
)

func RandBytes(b *testing.B, n int64) []byte {
	buff := make([]byte, n)

	_, err := rand.Read(buff)
	if err != nil {
		b.Error(err)
	}

	return buff
}

func BytesCount(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%dB", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.0f%ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}
