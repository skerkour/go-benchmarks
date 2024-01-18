package hashing

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"testing"

	"github.com/skerkour/go-benchmarks/utils"
	zeeboblake3 "github.com/zeebo/blake3"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/sha3"
	lukechampineblake3 "lukechampine.com/blake3"
)

type Hasher interface {
	Hash(input []byte)
}

func BenchmarkHashing(b *testing.B) {
	benchmarks := []int64{
		64,
		1024,
		16 * 1024,
		64 * 1024,
		1024 * 1024,
		10 * 1024 * 1024,
		1024 * 1024 * 1024,
	}

	for _, size := range benchmarks {
		benchmarkHasher(size, "sha256", sha256Hasher{}, b)
		benchmarkHasher(size, "zeebo_blake3_256", zeeboBlake3Hasher{}, b)
		benchmarkHasher(size, "lukechampine_blake3_256", lukechampineBlake3Hasher{}, b)
		benchmarkHasher(size, "blake2b_256", blake2bHasher{}, b)
		benchmarkHasher(size, "blake2s_256", blake2sHasher{}, b)
		// benchmarkHasher("sha512/256", sha512_256Hasher{}, b)
		benchmarkHasher(size, "sha3", sha3Hasher{}, b)
		benchmarkHasher(size, "sha1", sha1Hasher{}, b)

		benchmarkHasher(size, "sha2_512", sha512Hasher{}, b)
		benchmarkHasher(size, "zeebo_blake3_512", zeeboBlake3_512Hasher{}, b)
		benchmarkHasher(size, "lukechampine_blake3_512", lukechampineBlake3_512Hasher{}, b)
		benchmarkHasher(size, "blake2b_512", blake2b512Hasher{}, b)
		benchmarkHasher(size, "sha3_512", sha3_512Hasher{}, b)
	}
}

func benchmarkHasher[H Hasher](size int64, algorithm string, hasher H, b *testing.B) {
	b.Run(fmt.Sprintf("%s-%s", utils.BytesCount(size), algorithm), func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(size)
		buf := utils.RandBytes(b, size)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			hasher.Hash(buf)
		}
	})
}

type lukechampineBlake3Hasher struct{}

func (lukechampineBlake3Hasher) Hash(input []byte) {
	lukechampineblake3.Sum256(input)
}

type lukechampineBlake3_512Hasher struct{}

func (lukechampineBlake3_512Hasher) Hash(input []byte) {
	lukechampineblake3.Sum512(input)
}

type zeeboBlake3Hasher struct{}

func (zeeboBlake3Hasher) Hash(input []byte) {
	zeeboblake3.Sum256(input)
}

type zeeboBlake3_512Hasher struct{}

func (zeeboBlake3_512Hasher) Hash(input []byte) {
	zeeboblake3.Sum512(input)
}

type blake2sHasher struct{}

func (blake2sHasher) Hash(input []byte) {
	blake2s.Sum256(input)
}

type blake2bHasher struct{}

func (blake2bHasher) Hash(input []byte) {
	blake2b.Sum256(input)
}

type blake2b512Hasher struct{}

func (blake2b512Hasher) Hash(input []byte) {
	blake2b.Sum512(input)
}

type sha256Hasher struct{}

func (sha256Hasher) Hash(input []byte) {
	sha256.Sum256(input)
}

type sha512Hasher struct{}

func (sha512Hasher) Hash(input []byte) {
	sha512.Sum512(input)
}

// type sha512_256Hasher struct{}
// func (sha512_256Hasher) Hash(input []byte) {
// 	sha512.Sum512_256(input)
// }

type sha1Hasher struct{}

func (sha1Hasher) Hash(input []byte) {
	sha1.Sum(input)
}

type sha3Hasher struct{}

func (sha3Hasher) Hash(input []byte) {
	sha3.Sum256(input)
}

type sha3_512Hasher struct{}

func (sha3_512Hasher) Hash(input []byte) {
	sha3.Sum512(input)
}
