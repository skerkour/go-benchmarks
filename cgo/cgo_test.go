package cgobench

import (
	"testing"
)

func BenchmarkCGO(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	CallCgo(b.N)
}

// BenchmarkGo must be called with `-gcflags -l` to avoid inlining.
func BenchmarkGo(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	CallGo(b.N)
}
