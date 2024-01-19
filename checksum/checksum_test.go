package checksum

import (
	"fmt"
	"hash/crc32"
	"hash/crc64"
	"testing"

	"github.com/cespare/xxhash/v2"
	"github.com/skerkour/go-benchmarks/utils"
	"github.com/zeebo/xxh3"
)

type Checksumer interface {
	Checksum(input []byte)
}

func BenchmarkChecksum(b *testing.B) {
	benchmarks := []int64{
		64,
		1024,
		16 * 1024,
		64 * 1024,
		1024 * 1024,
		10 * 1024 * 1024,
		100 * 1024 * 1024,
		1024 * 1024 * 1024,
	}

	for _, size := range benchmarks {
		benchmarkChecksumer(size, "crc32", crc32Checksumer{}, b)
		benchmarkChecksumer(size, "crc64", NewCrc64Checksumer(), b)
		benchmarkChecksumer(size, "xxh3", xxh3Checksumer{}, b)
		benchmarkChecksumer(size, "xxh3_128", xxh3_128Checksumer{}, b)
		benchmarkChecksumer(size, "xxhash", xxh3_128Checksumer{}, b)
	}
}

func benchmarkChecksumer[C Checksumer](size int64, algorithm string, checksumer C, b *testing.B) {
	b.Run(fmt.Sprintf("%s-%s", utils.BytesCount(size), algorithm), func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(size)
		buf := utils.RandBytes(b, size)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			checksumer.Checksum(buf)
		}
	})
}

type crc32Checksumer struct{}

func (crc32Checksumer) Checksum(input []byte) {
	crc32.ChecksumIEEE(input)
}

type crc64Checksumer struct {
	table *crc64.Table
}

func NewCrc64Checksumer() crc64Checksumer {
	return crc64Checksumer{
		table: crc64.MakeTable(crc64.ISO),
	}
}

func (checksumer crc64Checksumer) Checksum(input []byte) {
	crc64.Checksum(input, checksumer.table)
}

type xxh3Checksumer struct{}

func (xxh3Checksumer) Checksum(input []byte) {
	xxh3.Hash(input)
}

type xxh3_128Checksumer struct{}

func (xxh3_128Checksumer) Checksum(input []byte) {
	xxh3.Hash128(input)
}

type xxhashChecksummer struct{}

func (xxhashChecksummer) Checksum(input []byte) {
	xxhash.Sum64(input)
}
