package encryption_aead

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"testing"

	"github.com/bloom42/stdx/crypto/bchacha20blake3"
	"github.com/bloom42/stdx/crypto/chacha20"
	"github.com/bloom42/stdx/crypto/chacha20blake3"
	"github.com/bloom42/stdx/crypto/schacha20blake3"
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
		100 * 1024 * 1024,
		1024 * 1024 * 1024,
	}
)

type AEADCipher interface {
	Encrypt(dst, nonce, plaintext, additionalData []byte) []byte
	Decrypt(dst, nonce, ciphertext, additionalData []byte)
}

func BenchmarkEncryptAEAD(b *testing.B) {
	additionalData := utils.RandBytes(b, 100)

	xChaCha20Key := utils.RandBytes(b, chacha20.KeySize)
	xChaCha20Nonce := utils.RandBytes(b, chacha20.NonceSizeX)

	chaCha20Key := utils.RandBytes(b, chacha20poly1305.KeySize)
	chaCha20Nonce := utils.RandBytes(b, chacha20poly1305.NonceSize)

	bChaCha20Nonce := utils.RandBytes(b, bchacha20blake3.NonceSize)

	aes256GcmKey := utils.RandBytes(b, 32)
	aes256GcmNonce := utils.RandBytes(b, 12)

	aes128GcmKey := utils.RandBytes(b, 16)
	aes128GcmNonce := utils.RandBytes(b, 12)

	for _, size := range BENCHMARKS {
		benchmarkEncrypt(b, size, "XChaCha20_BLAKE3", newXChaCha20Blake3Cipher(b, xChaCha20Key), xChaCha20Nonce, additionalData)
		benchmarkEncrypt(b, size, "XChaCha20_Poly1305", newXChaCha20Poly1305Cipher(b, xChaCha20Key), xChaCha20Nonce, additionalData)
		benchmarkEncrypt(b, size, "ChaCha20_Poly1305", newChaCha20Poly1305Cipher(b, chaCha20Key), chaCha20Nonce, additionalData)
		benchmarkEncrypt(b, size, "AES_256_GCM", newAesGcmCipher(b, aes256GcmKey), aes256GcmNonce, additionalData)
		benchmarkEncrypt(b, size, "AES_128_GCM", newAesGcmCipher(b, aes128GcmKey), aes128GcmNonce, additionalData)
		benchmarkEncrypt(b, size, "BChaCha20_BLAKE3", newBChaCha20Blake3Cipher(b, xChaCha20Key), bChaCha20Nonce, additionalData)
		benchmarkEncrypt(b, size, "SChaCha20_BLAKE3", newSChaCha20Blake3Cipher(b, xChaCha20Key), bChaCha20Nonce, additionalData)
	}
}

func BenchmarkDecryptAEAD(b *testing.B) {
	additionalData := utils.RandBytes(b, 100)

	xChaCha20Key := utils.RandBytes(b, chacha20.KeySize)
	xChaCha20Nonce := utils.RandBytes(b, chacha20.NonceSizeX)

	chaCha20Key := utils.RandBytes(b, chacha20poly1305.KeySize)
	chaCha20Nonce := utils.RandBytes(b, chacha20poly1305.NonceSize)

	aes256GcmKey := utils.RandBytes(b, 32)
	aes256GcmNonce := utils.RandBytes(b, 12)

	aes128GcmKey := utils.RandBytes(b, 16)
	aes128GcmNonce := utils.RandBytes(b, 12)

	bChaCha20Nonce := utils.RandBytes(b, bchacha20blake3.NonceSize)

	for _, size := range BENCHMARKS {
		benchmarkDecrypt(b, size, "XChaCha20_BLAKE3", newXChaCha20Blake3Cipher(b, xChaCha20Key), xChaCha20Nonce, additionalData)
		benchmarkDecrypt(b, size, "XChaCha20_Poly1305", newXChaCha20Poly1305Cipher(b, xChaCha20Key), xChaCha20Nonce, additionalData)
		benchmarkDecrypt(b, size, "ChaCha20_Poly1305", newChaCha20Poly1305Cipher(b, chaCha20Key), chaCha20Nonce, additionalData)
		benchmarkDecrypt(b, size, "AES_256_GCM", newAesGcmCipher(b, aes256GcmKey), aes256GcmNonce, additionalData)
		benchmarkDecrypt(b, size, "AES_128_GCM", newAesGcmCipher(b, aes128GcmKey), aes128GcmNonce, additionalData)
		benchmarkDecrypt(b, size, "BChaCha20_BLAKE3", newBChaCha20Blake3Cipher(b, xChaCha20Key), bChaCha20Nonce, additionalData)
		benchmarkDecrypt(b, size, "SChaCha20_BLAKE3", newSChaCha20Blake3Cipher(b, xChaCha20Key), bChaCha20Nonce, additionalData)
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
		cipherText := make([]byte, len(plaintext)+512)
		cipherText = cipher.Encrypt(cipherText, nonce, plaintext, additionalData)
		dst := make([]byte, len(cipherText))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cipher.Decrypt(dst, nonce, cipherText, additionalData)
		}
	})
}

type xChaCha20Blake3Cipher struct {
	cipher cipher.AEAD
}

func newXChaCha20Blake3Cipher(b *testing.B, key []byte) xChaCha20Blake3Cipher {
	cipher, err := chacha20blake3.NewX(key)
	if err != nil {
		b.Error(err)
	}

	return xChaCha20Blake3Cipher{
		cipher: cipher,
	}
}

func (cipher xChaCha20Blake3Cipher) Encrypt(dst, nonce, plaintext, additionalData []byte) []byte {
	return cipher.cipher.Seal(dst, nonce, plaintext, additionalData)
}

func (cipher xChaCha20Blake3Cipher) Decrypt(dst, nonce, ciphertext, additionalData []byte) {
	_, _ = cipher.cipher.Open(dst, nonce, ciphertext, additionalData)
}

type bChaCha20Blake3Cipher struct {
	cipher cipher.AEAD
}

func newBChaCha20Blake3Cipher(b *testing.B, key []byte) bChaCha20Blake3Cipher {
	cipher, err := bchacha20blake3.New(key)
	if err != nil {
		b.Error(err)
	}

	return bChaCha20Blake3Cipher{
		cipher: cipher,
	}
}

func (cipher bChaCha20Blake3Cipher) Encrypt(dst, nonce, plaintext, additionalData []byte) []byte {
	return cipher.cipher.Seal(dst, nonce, plaintext, additionalData)
}

func (cipher bChaCha20Blake3Cipher) Decrypt(dst, nonce, ciphertext, additionalData []byte) {
	_, _ = cipher.cipher.Open(dst, nonce, ciphertext, additionalData)
}

type sChaCha20Blake3Cipher struct {
	cipher cipher.AEAD
}

func newSChaCha20Blake3Cipher(b *testing.B, key []byte) sChaCha20Blake3Cipher {
	cipher, err := schacha20blake3.New(key)
	if err != nil {
		b.Error(err)
	}

	return sChaCha20Blake3Cipher{
		cipher: cipher,
	}
}

func (cipher sChaCha20Blake3Cipher) Encrypt(dst, nonce, plaintext, additionalData []byte) []byte {
	return cipher.cipher.Seal(dst, nonce, plaintext, additionalData)
}

func (cipher sChaCha20Blake3Cipher) Decrypt(dst, nonce, ciphertext, additionalData []byte) {
	_, _ = cipher.cipher.Open(dst, nonce, ciphertext, additionalData)
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

func (cipher xChaCha20Poly1305Cipher) Decrypt(dst, nonce, ciphertext, additionalData []byte) {
	_, _ = cipher.cipher.Open(dst, nonce, ciphertext, additionalData)
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

func (cipher chaCha20Poly1305Cipher) Decrypt(dst, nonce, ciphertext, additionalData []byte) {
	_, _ = cipher.cipher.Open(dst, nonce, ciphertext, additionalData)
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
