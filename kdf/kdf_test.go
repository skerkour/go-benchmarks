package kdf

import (
	"crypto/hkdf"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha3"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/skerkour/go-benchmarks/crypto/kmac"
	"github.com/skerkour/go-benchmarks/utils"
	"github.com/skerkour/stdx-go/crypto/chacha20"
	zeeboblake3 "github.com/zeebo/blake3"
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

	info := []byte(base64.StdEncoding.EncodeToString(utils.RandBytes(b, 30)))
	key := utils.RandBytes(b, 32)

	for _, size := range benchmarks {
		benchmarkKDF(size, "HKDF-SHA2-256", sha256KDF{}, key, info, b)
		benchmarkKDF(size, "HKDF-SHA2-512", sha512KDF{}, key, info, b)
		benchmarkKDF(size, "SHAKE-256", shake256Kdf{}, key, info, b)

		benchmarkKDF(size, "KMAC-128", kmac128{}, key, info, b)
		benchmarkKDF(size, "KMAC-256", kmac256{}, key, info, b)

		benchmarkKDF(size, "BLAKE3_zeebo", zeeboBlake3KDF{}, key, info, b)
		// benchmarkKDF(size, "BLAKE3-512_zeebo", zeeboBlake3_512KDF{}, key, info, output512, b)
		benchmarkKDF(size, "BLAKE3_lukechampine", lukechampineBlake3KDF{}, key, info, b)
		// benchmarkKDF(size, "BLAKE3-512_lukechampine", lukechampineBlake3_512KDF{}, key, info, output512, b)

		benchmarkKDF(size, "ChaCha20", newChacha20KDF(key), key, info, b)
		// benchmarkHasher(size, "blake2b_256", blake2bHasher{}, b)
		// benchmarkHasher(size, "blake2s_256", blake2sHasher{}, b)
		// benchmarkHasher("sha512/256", sha512_256Hasher{}, b)
		// benchmarkHasher(size, "sha3", sha3Hasher{}, b)

		// benchmarkHasher(size, "blake2b_512", blake2b512Hasher{}, b)
		// benchmarkHasher(size, "sha3_512", sha3_512Hasher{}, b)
	}
}

func benchmarkKDF[H KDF](size int64, algorithm string, kdf H, key, info []byte, b *testing.B) {
	b.Run(fmt.Sprintf("%s-%s", utils.BytesCount(size), algorithm), func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(size)
		b.ResetTimer()
		output := make([]byte, size)
		for i := 0; i < b.N; i++ {
			kdf.DeriveKey(key, info, output)
		}
	})
}

type lukechampineBlake3KDF struct{}

func (lukechampineBlake3KDF) DeriveKey(secret, info, out []byte) {
	lukechampineblake3.DeriveKey(out, string(info), secret)
}

// type lukechampineBlake3_512KDF struct{}

// func (lukechampineBlake3_512KDF) DeriveKey(secret, info, out []byte) {
// 	lukechampineblake3.DeriveKey(out, string(info), secret)
// }

type zeeboBlake3KDF struct{}

func (zeeboBlake3KDF) DeriveKey(secret, info, out []byte) {
	zeeboblake3.DeriveKey(string(info), secret, out[:0])
}

// type zeeboBlake3_512KDF struct{}

// func (zeeboBlake3_512KDF) DeriveKey(secret, info, out []byte) {
// 	zeeboblake3.DeriveKey(string(info), secret, out)
// }

type chacha20KDF struct {
	cipher chacha20.StreamCipher
}

func newChacha20KDF(secret []byte) chacha20KDF {
	var nonce [8]byte
	rand.Read(nonce[:])

	cipher, err := chacha20.New(secret, nonce[:])
	if err != nil {
		panic(err)
	}

	return chacha20KDF{
		cipher: cipher,
	}
}

func (kdf chacha20KDF) DeriveKey(secret, info, out []byte) {
	kdf.cipher.XORKeyStream(out[:], out[:])
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
	out, err := hkdf.Key(sha256.New, secret, nil, string(info), len(out))
	if err != nil {
		panic(err)
	}
}

type sha512KDF struct{}

func (sha512KDF) DeriveKey(secret, info, out []byte) {
	out, err := hkdf.Key(sha512.New, secret, nil, string(info), len(out))
	if err != nil {
		panic(err)
	}
}

type shake256Kdf struct{}

func (shake256Kdf) DeriveKey(secret, info, out []byte) {
	hasher := sha3.NewSHAKE256()
	hasher.Write(info)
	hasher.Write(secret)
	hasher.Read(out)
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

type kmac128 struct{}

func (kmac128) DeriveKey(key, input, output []byte) {
	hasher := kmac.NewKMAC128(key, 64, []byte("KDF"))
	hasher.Write(input)
	hasher.Sum(output[:0])
}

type kmac256 struct{}

func (kmac256) DeriveKey(key, input, output []byte) {
	hasher := kmac.NewKMAC256(key, len(output), []byte("KDF"))
	hasher.Write(input)
	hasher.Sum(output[:0])
}
