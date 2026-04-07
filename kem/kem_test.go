package kem

import (
	"crypto/mlkem"
	"testing"

	"github.com/skerkour/go-benchmarks/crypto/xwing"
)

type EncapulationKey interface {
	Encapsulate() (ciphertext, sharedKey []byte)
}

type DecapsulationKey interface {
	Decapsulate(ciphertext []byte) (sharedKey []byte, err error)
}

func BenchmarkEncapsulate(b *testing.B) {
	mlKem768, err := mlkem.GenerateKey768()
	if err != nil {
		b.Fatal(err)
	}
	mlKem768EncapsulationKey := mlKem768.EncapsulationKey()

	mlKem1024, err := mlkem.GenerateKey1024()
	if err != nil {
		b.Fatal(err)
	}
	mlKem1024EncapsulationKey := mlKem1024.EncapsulationKey()

	xwing, err := xwing.GenerateKey()
	if err != nil {
		b.Fatal(err)
	}
	xwingEncapsulationKey := xwing.EncapsulationKey()

	benchmarkEncapsulate("ML-KEM-768", mlKem768EncapsulationKey, b)
	benchmarkEncapsulate("ML-KEM-1024", mlKem1024EncapsulationKey, b)
	benchmarkEncapsulate("X-Wing", xwingEncapsulationKey, b)
}

func BenchmarkDecapsulate(b *testing.B) {
	mlKem768, err := mlkem.GenerateKey768()
	if err != nil {
		b.Fatal(err)
	}
	mlKem768EncapsulationKey := mlKem768.EncapsulationKey()
	mlKem768Ciphertext, _ := mlKem768EncapsulationKey.Encapsulate()
	// TODO: the mlkem API is likely to change due to a mismatch between the Go API and the spec
	// https://github.com/golang/go/issues/70950
	if len(mlKem768Ciphertext) != mlkem.CiphertextSize768 {
		panic("wrong ciphertext 768 size")
	}

	mlKem1024, err := mlkem.GenerateKey1024()
	if err != nil {
		b.Fatal(err)
	}
	mlKem1024EncapsulationKey := mlKem1024.EncapsulationKey()
	mlKem1024Ciphertext, _ := mlKem1024EncapsulationKey.Encapsulate()
	// TODO: the mlkem API is likely to change due to a mismatch between the Go API and the spec
	// https://github.com/golang/go/issues/70950
	if len(mlKem1024Ciphertext) != mlkem.CiphertextSize1024 {
		panic("wrong ciphertext 1024 size")
	}

	xwingKem, err := xwing.GenerateKey()
	if err != nil {
		b.Fatal(err)
	}
	xwingKemEncapsulationKey := xwingKem.EncapsulationKey()
	xwingCiphertext, _ := xwingKemEncapsulationKey.Encapsulate()

	benchmarkDecapsulate("ML-KEM-768", mlKem768, mlKem768Ciphertext, b)
	benchmarkDecapsulate("ML-KEM-1024", mlKem1024, mlKem1024Ciphertext, b)
	benchmarkDecapsulate("X-Wing", xwingKem, xwingCiphertext, b)
}

func benchmarkEncapsulate[K EncapulationKey](algorithm string, kem K, b *testing.B) {
	b.Run(algorithm, func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = kem.Encapsulate()
		}
	})
}

func benchmarkDecapsulate[K DecapsulationKey](algorithm string, kem K, ciphertext []byte, b *testing.B) {
	b.Run(algorithm, func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := kem.Decapsulate(ciphertext)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
