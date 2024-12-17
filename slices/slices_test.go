package slices

import (
	"crypto/rand"
	"testing"
)

func BenchmarkSlices(b *testing.B) {
	source := make([]byte, 10_000)
	rand.Read(source)

	b.Run("append", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(source)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			dest := make([]byte, 0, len(source))
			dest = append(dest, source...)
			_ = dest
		}
	})

	b.Run("copy", func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(source)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			dest := make([]byte, len(source))
			copy(dest, source)
		}
	})
}
