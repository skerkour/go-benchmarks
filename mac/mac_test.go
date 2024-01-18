package mac

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"testing"

	"github.com/skerkour/go-benchmarks/utils"
	zeeboblake3 "github.com/zeebo/blake3"
	lukechampineblake3 "lukechampine.com/blake3"
)

type Mac interface {
	Mac(key, input []byte)
}

func BenchmarkMac(b *testing.B) {
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
		benchmarkMac(size, "sha256", sha256Mac{}, b)
		// benchmarkMac(size, "blake2b_256", blake2bHasher{}, b)
		// benchmarkMac(size, "blake2s_256", blake2sHasher{}, b)
		// benchmarkMac("sha512/256", sha512_256Hasher{}, b)
		// benchmarkMac(size, "sha3", sha3Hasher{}, b)
		benchmarkMac(size, "lukechampine_blake3_256", lukechampineBlake3Mac{}, b)
		benchmarkMac(size, "zeebo_blake3_256", zeeboBlake3Mac{}, b)

		benchmarkMac(size, "sha2_512", sha512Hasher{}, b)
		// benchmarkMac(size, "blake2b_512", blake2b512Hasher{}, b)
		// benchmarkMac(size, "sha3_512", sha3_512Hasher{}, b)
		benchmarkMac(size, "lukechampine_blake3_512", lukechampineBlake3_512Mac{}, b)
		benchmarkMac(size, "zeebo_blake3_512", zeeboBlake3_512Mac{}, b)
	}
}

func benchmarkMac[H Mac](size int64, algorithm string, hasher H, b *testing.B) {
	b.Run(fmt.Sprintf("%s-%s", utils.BytesCount(size), algorithm), func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(size)
		key := utils.RandBytes(b, 32)
		buf := utils.RandBytes(b, size)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			hasher.Mac(key, buf)
		}
	})
}

type lukechampineBlake3Mac struct{}

func (lukechampineBlake3Mac) Mac(key, input []byte) {
	out := make([]byte, 64)
	hasher := lukechampineblake3.New(32, key)
	hasher.Write(input)
	hasher.Sum(out)
}

type lukechampineBlake3_512Mac struct{}

func (lukechampineBlake3_512Mac) Mac(key, input []byte) {
	out := make([]byte, 64)
	hasher := lukechampineblake3.New(64, key)
	hasher.Write(input)
	hasher.Sum(out)
}

type zeeboBlake3Mac struct{}

func (zeeboBlake3Mac) Mac(key, input []byte) {
	out := make([]byte, 64)
	hasher, _ := zeeboblake3.NewKeyed(key)
	hasher.Write(input)
	hasher.Sum(out)
}

type zeeboBlake3_512Mac struct{}

func (zeeboBlake3_512Mac) Mac(key, input []byte) {
	out := make([]byte, 64)
	hasher, _ := zeeboblake3.NewKeyed(key)
	hasher.Write(input)
	digest := hasher.Digest()
	digest.Read(out)
}

// type blake2sHasher struct{}

// func (blake2sHasher) Hash(input []byte) {
// 	blake2s.Sum256(input)
// }

// type blake2bHasher struct{}

// func (blake2bHasher) Hash(input []byte) {
// 	blake2b.Sum256(input)
// }

// type blake2b512Hasher struct{}

// func (blake2b512Hasher) Hash(input []byte) {
// 	blake2b.Sum512(input)
// }

type sha256Mac struct{}

func (sha256Mac) Mac(key, input []byte) {
	out := make([]byte, 64)
	hmac := hmac.New(sha256.New, key)
	hmac.Write(input)
	hmac.Sum(out)
}

type sha512Hasher struct{}

func (sha512Hasher) Mac(key, input []byte) {
	out := make([]byte, 64)
	hmac := hmac.New(sha512.New, key)
	hmac.Write(input)
	hmac.Sum(out)
}

// type sha512_256Hasher struct{}
// func (sha512_256Hasher) Hash(input []byte) {
// 	sha512.Sum512_256(input)
// }

// type sha3Hasher struct{}

// func (sha3Hasher) Hash(input []byte) {
// 	sha3.Sum256(input)
// }

// type sha3_512Hasher struct{}

// func (sha3_512Hasher) Hash(input []byte) {
// 	sha3.Sum512(input)
// }