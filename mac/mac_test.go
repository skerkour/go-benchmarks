package mac

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"testing"

	"github.com/skerkour/go-benchmarks/utils"
	zeeboblake3 "github.com/zeebo/blake3"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/poly1305"
	"golang.org/x/crypto/sha3"
	lukechampineblake3 "lukechampine.com/blake3"
)

type Mac interface {
	Mac(key, input, output []byte)
}

func BenchmarkMac(b *testing.B) {
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

	output128 := make([]byte, 16, 256)
	output256 := make([]byte, 32, 256)
	output512 := make([]byte, 64, 256)

	for _, size := range benchmarks {
		benchmarkMac(size, "sha256", sha256Mac{}, output256, b)
		benchmarkMac(size, "zeebo_blake3_256", zeeboBlake3Mac{}, output256, b)
		benchmarkMac(size, "lukechampine_blake3_256", lukechampineBlake3Mac{}, output256, b)
		benchmarkMac(size, "blake2b_256", blake2bMac{}, output256, b)
		benchmarkMac(size, "blake2s_256", blake2sMac{}, output256, b)
		// benchmarkMac("sha512/256", sha512_256Hasher{}, b)
		benchmarkMac(size, "sha3", sha3Mac{}, output256, b)
		benchmarkMac(size, "poly1305", poly1305Mac{}, output128, b)

		benchmarkMac(size, "sha2_512", sha512Hasher{}, output512, b)
		benchmarkMac(size, "zeebo_blake3_512", zeeboBlake3_512Mac{}, output512, b)
		benchmarkMac(size, "lukechampine_blake3_512", lukechampineBlake3_512Mac{}, output512, b)
		// benchmarkMac(size, "blake2b_512", blake2b512Hasher{}, b)
		benchmarkMac(size, "sha3_512", sha3_512Mac{}, output512, b)
	}
}

func benchmarkMac[H Mac](size int64, algorithm string, hasher H, output []byte, b *testing.B) {
	b.Run(fmt.Sprintf("%s-%s", utils.BytesCount(size), algorithm), func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(size)
		key := utils.RandBytes(b, 32)
		buf := utils.RandBytes(b, size)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			hasher.Mac(key, buf, output)
		}
	})
}

type lukechampineBlake3Mac struct{}

func (lukechampineBlake3Mac) Mac(key, input, output []byte) {
	hasher := lukechampineblake3.New(32, key)
	hasher.Write(input)
	hasher.Sum(output)
}

type lukechampineBlake3_512Mac struct{}

func (lukechampineBlake3_512Mac) Mac(key, input, output []byte) {
	hasher := lukechampineblake3.New(64, key)
	hasher.Write(input)
	hasher.Sum(output)
}

type zeeboBlake3Mac struct{}

func (zeeboBlake3Mac) Mac(key, input, output []byte) {
	hasher, _ := zeeboblake3.NewKeyed(key)
	hasher.Write(input)
	hasher.Sum(output)
}

type zeeboBlake3_512Mac struct{}

func (zeeboBlake3_512Mac) Mac(key, input, output []byte) {
	hasher, _ := zeeboblake3.NewKeyed(key)
	hasher.Write(input)
	digest := hasher.Digest()
	digest.Read(output)
}

type blake2sMac struct{}

func (blake2sMac) Mac(key, input, output []byte) {
	hasher, _ := blake2s.New256(key)
	hasher.Write(input)
	hasher.Sum(output)
}

type blake2bMac struct{}

func (blake2bMac) Mac(key, input, output []byte) {
	hasher, _ := blake2b.New(32, key)
	hasher.Write(input)
	hasher.Sum(output)
}

type poly1305Mac struct{}

func (poly1305Mac) Mac(key, input, output []byte) {
	polyKey := [32]byte(key[0:32])
	hasher := poly1305.New(&polyKey)
	hasher.Write(input)
	hasher.Sum(output)
}

// type blake2b512Hasher struct{}

// func (blake2b512Hasher) Hash(input []byte) {
// 	blake2b.Sum512(input)
// }

type sha256Mac struct{}

func (sha256Mac) Mac(key, input, output []byte) {
	hmac := hmac.New(sha256.New, key)
	hmac.Write(input)
	hmac.Sum(output)
}

type sha512Hasher struct{}

func (sha512Hasher) Mac(key, input, output []byte) {
	hmac := hmac.New(sha512.New, key)
	hmac.Write(input)
	hmac.Sum(output)
}

// type sha512_256Hasher struct{}
// func (sha512_256Hasher) Hash(input []byte) {
// 	sha512.Sum512_256(input)
// }

type sha3Mac struct{}

func (sha3Mac) Mac(key, input, output []byte) {
	hmac := hmac.New(sha3.New256, key)
	hmac.Write(input)
	hmac.Sum(output)
}

type sha3_512Mac struct{}

func (sha3_512Mac) Mac(key, input, output []byte) {
	hmac := hmac.New(sha3.New512, key)
	hmac.Write(input)
	hmac.Sum(output)
}
