package signatures

import (
	"crypto/ed25519"
	"fmt"
	"testing"

	"github.com/skerkour/go-benchmarks/utils"
)

type Signer interface {
	Sign(message []byte) []byte
	Verify(message, signature []byte) bool
}

func BenchmarkSign(b *testing.B) {
	benchmarks := []int64{
		64,
		1024,
		64 * 1024,
		1 * 1024 * 1024,
		1024 * 1024 * 1024,
	}

	for _, size := range benchmarks {
		benchmarkSign(size, "ed25519", newEd25519Signer(b), b)
	}
}

func BenchmarkVerify(b *testing.B) {
	benchmarks := []int64{
		64,
		1024,
		64 * 1024,
		1 * 1024 * 1024,
		1024 * 1024 * 1024,
	}

	for _, size := range benchmarks {
		benchmarkVerify(size, "ed25519", newEd25519Signer(b), b)
	}
}

func benchmarkSign[S Signer](size int64, algorithm string, signer S, b *testing.B) {
	b.Run(fmt.Sprintf("%s-%s", utils.BytesCount(size), algorithm), func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(size)
		buf := utils.RandBytes(b, size)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			signer.Sign(buf)
		}
	})
}

func benchmarkVerify[S Signer](size int64, algorithm string, signer S, b *testing.B) {
	b.Run(fmt.Sprintf("%s-%s", utils.BytesCount(size), algorithm), func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(size)
		buf := utils.RandBytes(b, size)
		signature := signer.Sign(buf)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			signer.Verify(buf, signature)
		}
	})
}

type ed25519Signer struct {
	privakeKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
}

func newEd25519Signer(b *testing.B) (signer ed25519Signer) {
	public, private, err := ed25519.GenerateKey(nil)
	if err != nil {
		b.Error(err)
	}

	signer = ed25519Signer{
		privakeKey: private,
		publicKey:  public,
	}
	return
}

func (signer ed25519Signer) Sign(message []byte) []byte {
	return ed25519.Sign(signer.privakeKey, message)
}

func (signer ed25519Signer) Verify(message, signature []byte) bool {
	return ed25519.Verify(signer.publicKey, message, signature)
}
