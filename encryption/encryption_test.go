package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"testing"

	"github.com/skerkour/go-benchmarks/utils"
	"golang.org/x/crypto/chacha20poly1305"
)

var (
	BENCHMARKS = []int64{
		64,
		1024,
		16 * 1024,
		64 * 1024,
		1024 * 1024,
		10 * 1024 * 1024,
		1024 * 1024 * 1024,
	}
)

type AEADCipher interface {
	Encrypt(dst, nonce, plaintext, additionalData []byte) []byte
	Decrypt(dst, nonce, ciphertext, additionalData []byte)
}

func BenchmarkEncrypt(b *testing.B) {
	additionalData := utils.RandBytes(b, 100)

	xChaCha20Poly1305Key := utils.RandBytes(b, chacha20poly1305.KeySize)
	xChaCha20Poly1305Nonce := utils.RandBytes(b, chacha20poly1305.NonceSizeX)

	chaCha20Poly1305Key := utils.RandBytes(b, chacha20poly1305.KeySize)
	chaCha20Poly1305Nonce := utils.RandBytes(b, chacha20poly1305.NonceSize)

	aes256GcmKey := utils.RandBytes(b, 32)
	aes256GcmNonce := utils.RandBytes(b, 12)

	aes128GcmKey := utils.RandBytes(b, 16)
	aes128GcmNonce := utils.RandBytes(b, 12)

	for _, size := range BENCHMARKS {
		benchmarkEncrypt(b, size, "XChaCha20_Poly1305", newXChaCha20Poly1305Cipher(b, xChaCha20Poly1305Key), xChaCha20Poly1305Nonce, additionalData)
		benchmarkEncrypt(b, size, "ChaCha20_Poly1305", newChaCha20Poly1305Cipher(b, chaCha20Poly1305Key), chaCha20Poly1305Nonce, additionalData)
		benchmarkEncrypt(b, size, "AES_128_GCM", newAesGcmCipher(b, aes128GcmKey), aes128GcmNonce, additionalData)
		benchmarkEncrypt(b, size, "AES_256_GCM", newAesGcmCipher(b, aes256GcmKey), aes256GcmNonce, additionalData)
	}
}

func BenchmarkDecrypt(b *testing.B) {
	additionalData := utils.RandBytes(b, 100)

	xChaCha20Poly1305Key := utils.RandBytes(b, chacha20poly1305.KeySize)
	xChaCha20Poly1305Nonce := utils.RandBytes(b, chacha20poly1305.NonceSizeX)

	chaCha20Poly1305Key := utils.RandBytes(b, chacha20poly1305.KeySize)
	chaCha20Poly1305Nonce := utils.RandBytes(b, chacha20poly1305.NonceSize)

	aesGcmKey := utils.RandBytes(b, 32)
	aesGcmNonce := utils.RandBytes(b, 12)

	aes128GcmKey := utils.RandBytes(b, 16)
	aes128GcmNonce := utils.RandBytes(b, 12)

	for _, size := range BENCHMARKS {
		benchmarkDecrypt(b, size, "XChaCha20_Poly1305", newXChaCha20Poly1305Cipher(b, xChaCha20Poly1305Key), xChaCha20Poly1305Nonce, additionalData)
		benchmarkDecrypt(b, size, "ChaCha20_Poly1305", newChaCha20Poly1305Cipher(b, chaCha20Poly1305Key), chaCha20Poly1305Nonce, additionalData)
		benchmarkDecrypt(b, size, "AES_128_GCM", newAesGcmCipher(b, aes128GcmKey), aes128GcmNonce, additionalData)
		benchmarkDecrypt(b, size, "AES_256_GCM", newAesGcmCipher(b, aesGcmKey), aesGcmNonce, additionalData)
	}
}

func benchmarkEncrypt[C AEADCipher](b *testing.B, size int64, algorithm string, cipher C, nonce, additionalData []byte) {
	b.Run(fmt.Sprintf("%s-%s", utils.BytesCount(size), algorithm), func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(size)
		plaintext := utils.RandBytes(b, size)
		dst := make([]byte, len(plaintext)+512)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cipher.Encrypt(dst, nonce, plaintext, additionalData)
		}
	})
}

func benchmarkDecrypt[C AEADCipher](b *testing.B, size int64, algorithm string, cipher C, nonce, additionalData []byte) {
	b.Run(fmt.Sprintf("%s-%s", utils.BytesCount(size), algorithm), func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(size)
		plaintext := utils.RandBytes(b, size)
		cipherText := cipher.Encrypt(nil, nonce, plaintext, additionalData)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cipher.Decrypt(plaintext, nonce, cipherText, additionalData)
		}
	})
}

type xChaCha20Poly1305Cipher struct {
	cipher cipher.AEAD
}

func newXChaCha20Poly1305Cipher(b *testing.B, key []byte) xChaCha20Poly1305Cipher {
	cipher, err := chacha20poly1305.NewX(key)
	if err != nil {
		b.Error(err)
	}

	return xChaCha20Poly1305Cipher{
		cipher: cipher,
	}
}

func (cipher xChaCha20Poly1305Cipher) Encrypt(dst, nonce, plaintext, additionalData []byte) []byte {
	return cipher.cipher.Seal(dst, nonce, plaintext, additionalData)
}

func (cipher xChaCha20Poly1305Cipher) Decrypt(dst, nonce, plaintext, additionalData []byte) {
	_, _ = cipher.cipher.Open(dst, nonce, plaintext, additionalData)
}

type chaCha20Poly1305Cipher struct {
	cipher cipher.AEAD
}

func newChaCha20Poly1305Cipher(b *testing.B, key []byte) chaCha20Poly1305Cipher {
	cipher, err := chacha20poly1305.New(key)
	if err != nil {
		b.Error(err)
	}

	return chaCha20Poly1305Cipher{
		cipher: cipher,
	}
}

func (cipher chaCha20Poly1305Cipher) Encrypt(dst, nonce, plaintext, additionalData []byte) []byte {
	return cipher.cipher.Seal(dst, nonce, plaintext, additionalData)
}

func (cipher chaCha20Poly1305Cipher) Decrypt(dst, nonce, plaintext, additionalData []byte) {
	_, _ = cipher.cipher.Open(dst, nonce, plaintext, additionalData)
}

type aesGcmCipher struct {
	cipher cipher.AEAD
}

func newAesGcmCipher(b *testing.B, key []byte) aesGcmCipher {
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		b.Error(err)
	}

	cipher, err := cipher.NewGCM(aesCipher)
	if err != nil {
		b.Error(err)
	}

	return aesGcmCipher{
		cipher: cipher,
	}
}

func (cipher aesGcmCipher) Encrypt(dst, nonce, plaintext, additionalData []byte) []byte {
	return cipher.cipher.Seal(dst, nonce, plaintext, additionalData)
}

func (cipher aesGcmCipher) Decrypt(dst, nonce, plaintext, additionalData []byte) {
	_, _ = cipher.cipher.Open(dst, nonce, plaintext, additionalData)
}
