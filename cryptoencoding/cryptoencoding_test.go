package cryptoencoding

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"testing"

	"github.com/skerkour/go-benchmarks/utils"
)

type Encoder interface {
	// Encode(data []byte)
	Decode(data []byte)
}

func BenchmarkDecode(b *testing.B) {
	p256PrivateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		b.Fatal(err)
	}

	p256PublicKeyPkix, err := x509.MarshalPKIXPublicKey(&p256PrivateKey.PublicKey)
	if err != nil {
		b.Fatal(err)
	}

	p256PublicKeyPemBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: p256PublicKeyPkix,
	}
	p256PublicKeyPem := pem.EncodeToMemory(p256PublicKeyPemBlock)

	publicKeyBinarryCompressed := elliptic.MarshalCompressed(elliptic.P256(), p256PrivateKey.X, p256PrivateKey.Y)

	benchmarkDecode("pkix", p256PublicKeyPkix, pkix{}, b)
	benchmarkDecode("pkix+PEM", p256PublicKeyPem, pkixPem{}, b)
	benchmarkDecode("binary", publicKeyBinarryCompressed, binaryCompressed{}, b)
}

func benchmarkDecode[E Encoder](algorithm string, data []byte, decoder E, b *testing.B) {
	b.Run(fmt.Sprintf("%s-%s", utils.BytesCount(int64(len(data))), algorithm), func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(int64(len(data)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			decoder.Decode(data)
		}
	})
}

type pkix struct{}

func (pkix) Decode(data []byte) {
	parsedPublicKey, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		panic(err)
	}
	ecdsaPublicKey, ecdsaPublicKeyOk := parsedPublicKey.(*ecdsa.PublicKey)
	if !ecdsaPublicKeyOk {
		panic("parsedPublicKey is not an *ecdsa.PublicKey")
	}
	_ = ecdsaPublicKey
}

type pkixPem struct{}

func (pkixPem) Decode(data []byte) {
	pemBlock, _ := pem.Decode(data)

	parsedPublicKey, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
	if err != nil {
		panic(err)
	}
	ecdsaPublicKey, ecdsaPublicKeyOk := parsedPublicKey.(*ecdsa.PublicKey)
	if !ecdsaPublicKeyOk {
		panic("parsedPublicKey is not an *ecdsa.PublicKey")
	}
	_ = ecdsaPublicKey
}

type binaryCompressed struct{}

func (binaryCompressed) Decode(data []byte) {
	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), data)
	if x == nil || y == nil {
		panic("error unmarshalling public key")
	}
}
