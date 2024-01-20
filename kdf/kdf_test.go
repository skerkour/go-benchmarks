package kdf

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"io"
	"testing"

	"github.com/bloom42/stdx/crypto/chacha20"
	"github.com/skerkour/go-benchmarks/utils"
	zeeboblake3 "github.com/zeebo/blake3"
	"golang.org/x/crypto/hkdf"
	lukechampineblake3 "lukechampine.com/blake3"
)

type KDF interface {
	DeriveKey(secret, info, out []byte)
}

func BenchmarkKDF(b *testing.B) {
	benchmarks := []int64{
		32,
		64,
		128,
		256,
	}

	info := utils.RandBytes(b, 24)
	key := utils.RandBytes(b, 32)
	output256 := make([]byte, 32, 256)
	output512 := make([]byte, 64, 256)

	for _, size := range benchmarks {
		benchmarkKDF(size, "hkdf_sha256", sha256KDF{}, key, info, output256, b)
		benchmarkKDF(size, "zeebo_blake3_256", zeeboBlake3KDF{}, key, info, output256, b)
		benchmarkKDF(size, "lukechampine_blake3_256", lukechampineBlake3KDF{}, key, info, output256, b)
		benchmarkKDF(size, "chacha20", chacha20KDF{}, key, info, output256, b)
		// benchmarkHasher(size, "blake2b_256", blake2bHasher{}, b)
		// benchmarkHasher(size, "blake2s_256", blake2sHasher{}, b)
		// benchmarkHasher("sha512/256", sha512_256Hasher{}, b)
		// benchmarkHasher(size, "sha3", sha3Hasher{}, b)

		benchmarkKDF(size, "hkdf_sha2_512", sha512KDF{}, key, info, output512, b)
		benchmarkKDF(size, "zeebo_blake3_512", zeeboBlake3_512KDF{}, key, info, output512, b)
		benchmarkKDF(size, "lukechampine_blake3_512", lukechampineBlake3_512KDF{}, key, info, output512, b)
		// benchmarkHasher(size, "blake2b_512", blake2b512Hasher{}, b)
		// benchmarkHasher(size, "sha3_512", sha3_512Hasher{}, b)
	}
}

func benchmarkKDF[H KDF](size int64, algorithm string, kdf H, key, info, output []byte, b *testing.B) {
	b.Run(fmt.Sprintf("%s-%s", utils.BytesCount(size), algorithm), func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(size)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			kdf.DeriveKey(key, info, output)
		}
	})
}

type lukechampineBlake3KDF struct{}

func (lukechampineBlake3KDF) DeriveKey(secret, info, out []byte) {
	lukechampineblake3.DeriveKey(out, string(info), secret)
}

type lukechampineBlake3_512KDF struct{}

func (lukechampineBlake3_512KDF) DeriveKey(secret, info, out []byte) {
	lukechampineblake3.DeriveKey(out, string(info), secret)
}

type zeeboBlake3KDF struct{}

func (zeeboBlake3KDF) DeriveKey(secret, info, out []byte) {
	zeeboblake3.DeriveKey(string(info), secret, out)
}

type zeeboBlake3_512KDF struct{}

func (zeeboBlake3_512KDF) DeriveKey(secret, info, out []byte) {
	zeeboblake3.DeriveKey(string(info), secret, out)
}

type chacha20KDF struct{}

func (chacha20KDF) DeriveKey(secret, info, out []byte) {
	cipher, _ := chacha20.New(secret, info)
	cipher.XORKeyStream(out[:], out[:])
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

type sha256KDF struct{}

func (sha256KDF) DeriveKey(secret, info, out []byte) {
	hkdf := hkdf.New(sha256.New, secret, nil, info)
	_, _ = io.ReadFull(hkdf, out)
}

type sha512KDF struct{}

func (sha512KDF) DeriveKey(secret, info, out []byte) {
	hkdf := hkdf.New(sha512.New, secret, nil, info)
	_, _ = io.ReadFull(hkdf, out)
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
