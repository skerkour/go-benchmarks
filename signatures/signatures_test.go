package signatures

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
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
		benchmarkSign(size, "Ed25519", newEd25519Signer(b), b)
		benchmarkSign(size, "ECDSA-P-256", newP256Signer(), b)
		benchmarkSign(size, "ECDSA-P-384", newP384Signer(), b)
		benchmarkSign(size, "ECDSA-P-521", newP521Signer(), b)
		benchmarkSign(size, "RSA-PKCS-1-v1.5-2048-SHA256", newRsaSha256Signer(b, 4096), b)
		benchmarkSign(size, "RSA-PKCS-1-v1.5-4096-SHA256", newRsaSha256Signer(b, 4096), b)
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
		benchmarkVerify(size, "Ed25519", newEd25519Signer(b), b)
		benchmarkVerify(size, "ECDSA-P-256", newP256Signer(), b)
		benchmarkVerify(size, "ECDSA-P-384", newP384Signer(), b)
		benchmarkVerify(size, "ECDSA-P-521", newP521Signer(), b)
		benchmarkVerify(size, "RSA-PKCS-1-v1.5-2048-SHA256", newRsaSha256Signer(b, 4096), b)
		benchmarkVerify(size, "RSA-PKCS-1-v1.5-4096-SHA256", newRsaSha256Signer(b, 4096), b)
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

type p256Signer struct {
	privakeKey *ecdsa.PrivateKey
	publicKey  ecdsa.PublicKey
}

func newP256Signer() (signer p256Signer) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	signer = p256Signer{
		privakeKey: privateKey,
		publicKey:  privateKey.PublicKey,
	}
	return
}

func (signer p256Signer) Sign(message []byte) []byte {
	hash := sha256.Sum256(message)
	ret, err := ecdsa.SignASN1(rand.Reader, signer.privakeKey, hash[:])
	if err != nil {
		panic(err)
	}
	return ret
}

func (signer p256Signer) Verify(message, signature []byte) bool {
	hash := sha256.Sum256(message)
	return ecdsa.VerifyASN1(&signer.publicKey, hash[:], signature)
}

type p384Signer struct {
	privakeKey *ecdsa.PrivateKey
	publicKey  ecdsa.PublicKey
}

func newP384Signer() (signer p384Signer) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		panic(err)
	}

	signer = p384Signer{
		privakeKey: privateKey,
		publicKey:  privateKey.PublicKey,
	}
	return
}

func (signer p384Signer) Sign(message []byte) []byte {
	hash := sha256.Sum256(message)
	ret, err := ecdsa.SignASN1(rand.Reader, signer.privakeKey, hash[:])
	if err != nil {
		panic(err)
	}
	return ret
}

func (signer p384Signer) Verify(message, signature []byte) bool {
	hash := sha256.Sum256(message)
	return ecdsa.VerifyASN1(&signer.publicKey, hash[:], signature)
}

type p521Signer struct {
	privakeKey *ecdsa.PrivateKey
	publicKey  ecdsa.PublicKey
}

func newP521Signer() (signer p521Signer) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		panic(err)
	}

	signer = p521Signer{
		privakeKey: privateKey,
		publicKey:  privateKey.PublicKey,
	}
	return
}

func (signer p521Signer) Sign(message []byte) []byte {
	hash := sha256.Sum256(message)
	ret, err := ecdsa.SignASN1(rand.Reader, signer.privakeKey, hash[:])
	if err != nil {
		panic(err)
	}
	return ret
}

func (signer p521Signer) Verify(message, signature []byte) bool {
	hash := sha256.Sum256(message)
	return ecdsa.VerifyASN1(&signer.publicKey, hash[:], signature)
}

type rsaSha256Signer struct {
	privakeKey *rsa.PrivateKey
	publicKey  rsa.PublicKey
}

func newRsaSha256Signer(b *testing.B, bits int) (signer rsaSha256Signer) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		b.Error(err)
	}

	signer = rsaSha256Signer{
		privakeKey: privateKey,
		publicKey:  privateKey.PublicKey,
	}
	return
}

func (signer rsaSha256Signer) Sign(message []byte) []byte {
	sha256Hash := sha256.Sum256(message)
	signature, err := rsa.SignPKCS1v15(nil, signer.privakeKey, crypto.SHA256, sha256Hash[:])
	if err != nil {
		panic(err)
	}
	return signature
}

func (signer rsaSha256Signer) Verify(message, signature []byte) bool {
	sha256Hash := sha256.Sum256(message)
	err := rsa.VerifyPKCS1v15(&signer.publicKey, crypto.SHA256, sha256Hash[:], signature)
	return err == nil
}
