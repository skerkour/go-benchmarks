package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"testing"

	"github.com/skerkour/go-benchmarks/utils"
	"golang.org/x/crypto/chacha20"
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

	xChaCha20Key := utils.RandBytes(b, chacha20.KeySize)
	xChaCha20Nonce := utils.RandBytes(b, chacha20.NonceSizeX)

	chaCha20Key := utils.RandBytes(b, chacha20.KeySize)
	chaCha20Nonce := utils.RandBytes(b, chacha20.NonceSize)

	aes256GcmKey := utils.RandBytes(b, 32)
	aes256GcmNonce := utils.RandBytes(b, 12)

	aes128GcmKey := utils.RandBytes(b, 16)
	aes128GcmNonce := utils.RandBytes(b, 12)

	aes256CbcKey := utils.RandBytes(b, 32)
	aes256CbcIv := utils.RandBytes(b, 16)

	aes256CfbKey := utils.RandBytes(b, 32)
	aes256CfbIv := utils.RandBytes(b, 16)

	for _, size := range BENCHMARKS {
		benchmarkEncrypt(b, size, "XChaCha20_Poly1305", newXChaCha20Poly1305Cipher(b, xChaCha20Key), xChaCha20Nonce, additionalData)
		benchmarkEncrypt(b, size, "XChaCha20", newxChaCha20Cipher(b, xChaCha20Key, xChaCha20Nonce), xChaCha20Nonce, additionalData)
		benchmarkEncrypt(b, size, "ChaCha20_Poly1305", newChaCha20Poly1305Cipher(b, chaCha20Key), chaCha20Nonce, additionalData)
		benchmarkEncrypt(b, size, "ChaCha20", newChaCha20Cipher(b, chaCha20Key, chaCha20Nonce), chaCha20Nonce, additionalData)
		benchmarkEncrypt(b, size, "AES_128_GCM", newAesGcmCipher(b, aes128GcmKey), aes128GcmNonce, additionalData)
		benchmarkEncrypt(b, size, "AES_256_GCM", newAesGcmCipher(b, aes256GcmKey), aes256GcmNonce, additionalData)
		benchmarkEncrypt(b, size, "AES_256_CBC", newAesCbcCipher(b, aes256CbcKey), aes256CbcIv, additionalData)
		benchmarkEncrypt(b, size, "AES_256_CFB", newAesCfbCipher(b, aes256CfbKey), aes256CfbIv, additionalData)
	}
}

func BenchmarkDecrypt(b *testing.B) {
	additionalData := utils.RandBytes(b, 100)

	xChaCha20Key := utils.RandBytes(b, chacha20.KeySize)
	xChaCha20Nonce := utils.RandBytes(b, chacha20.NonceSizeX)

	chaCha20Key := utils.RandBytes(b, chacha20poly1305.KeySize)
	chaCha20Nonce := utils.RandBytes(b, chacha20poly1305.NonceSize)

	aes256GcmKey := utils.RandBytes(b, 32)
	aes256GcmNonce := utils.RandBytes(b, 12)

	aes128GcmKey := utils.RandBytes(b, 16)
	aes128GcmNonce := utils.RandBytes(b, 12)

	aes256CbcKey := utils.RandBytes(b, 32)
	aes256CbcIv := utils.RandBytes(b, 16)

	aes256CfbKey := utils.RandBytes(b, 32)
	aes256CfbIv := utils.RandBytes(b, 16)

	for _, size := range BENCHMARKS {
		benchmarkDecrypt(b, size, "XChaCha20_Poly1305", newXChaCha20Poly1305Cipher(b, xChaCha20Key), xChaCha20Nonce, additionalData)
		benchmarkDecrypt(b, size, "XChaCha20", newxChaCha20Cipher(b, xChaCha20Key, xChaCha20Nonce), xChaCha20Nonce, additionalData)
		benchmarkDecrypt(b, size, "ChaCha20_Poly1305", newChaCha20Poly1305Cipher(b, chaCha20Key), chaCha20Nonce, additionalData)
		benchmarkDecrypt(b, size, "ChaCha20", newChaCha20Cipher(b, chaCha20Key, chaCha20Nonce), chaCha20Nonce, additionalData)
		benchmarkDecrypt(b, size, "AES_128_GCM", newAesGcmCipher(b, aes128GcmKey), aes128GcmNonce, additionalData)
		benchmarkDecrypt(b, size, "AES_256_GCM", newAesGcmCipher(b, aes256GcmKey), aes256GcmNonce, additionalData)
		benchmarkDecrypt(b, size, "AES_256_CBC", newAesCbcCipher(b, aes256CbcKey), aes256CbcIv, additionalData)
		benchmarkDecrypt(b, size, "AES_256_CFB", newAesCfbCipher(b, aes256CfbKey), aes256CfbIv, additionalData)
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

type xChaCha20Cipher struct {
	cipher *chacha20.Cipher
}

func newxChaCha20Cipher(b *testing.B, key, nonce []byte) xChaCha20Cipher {
	cipher, err := chacha20.NewUnauthenticatedCipher(key, nonce)
	if err != nil {
		b.Error(err)
	}

	return xChaCha20Cipher{
		cipher: cipher,
	}
}

func (cipher xChaCha20Cipher) Encrypt(dst, _nonce, plaintext, _additionalData []byte) []byte {
	cipher.cipher.XORKeyStream(dst, plaintext)
	return dst
}

func (cipher xChaCha20Cipher) Decrypt(dst, _nonce, ciphertext, _additionalData []byte) {
	cipher.cipher.XORKeyStream(dst, ciphertext)
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

type chaCha20Cipher struct {
	cipher *chacha20.Cipher
}

func newChaCha20Cipher(b *testing.B, key, nonce []byte) chaCha20Cipher {
	cipher, err := chacha20.NewUnauthenticatedCipher(key, nonce)
	if err != nil {
		b.Error(err)
	}

	return chaCha20Cipher{
		cipher: cipher,
	}
}

func (cipher chaCha20Cipher) Encrypt(dst, _nonce, plaintext, _additionalData []byte) []byte {
	cipher.cipher.XORKeyStream(dst, plaintext)
	return dst
}

func (cipher chaCha20Cipher) Decrypt(dst, _nonce, ciphertext, _additionalData []byte) {
	cipher.cipher.XORKeyStream(dst, ciphertext)
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

type aesCbcCipher struct {
	cipher cipher.Block
}

func newAesCbcCipher(b *testing.B, key []byte) aesCbcCipher {
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		b.Error(err)
	}

	// decrypter := cipher.NewCBCDecrypter(aesCipher, iv)

	return aesCbcCipher{
		cipher: aesCipher,
	}
}

func (aesCbcCipher aesCbcCipher) Encrypt(dst, nonce, plaintext, _additionalData []byte) []byte {
	encrypter := cipher.NewCBCEncrypter(aesCbcCipher.cipher, nonce)
	paddedPlaintext, _ := pkcs7Pad(plaintext, encrypter.BlockSize())
	encrypter.CryptBlocks(dst, paddedPlaintext)
	return dst
}

func (aesCbcCipher aesCbcCipher) Decrypt(dst, nonce, ciphertext, additionalData []byte) {
	decrypter := cipher.NewCBCDecrypter(aesCbcCipher.cipher, nonce)
	unpaddedCiphertext, err := pkcs7Unpad(ciphertext, decrypter.BlockSize())
	if err != nil {
		panic(err)
	}
	decrypter.CryptBlocks(dst, unpaddedCiphertext)
}

func PKCS5Padding(ciphertext []byte, blockSize int, after int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7 errors.
var (
	// ErrInvalidBlockSize indicates hash blocksize <= 0.
	ErrInvalidBlockSize = errors.New("invalid blocksize")

	// ErrInvalidPKCS7Data indicates bad input to PKCS7 pad or unpad.
	ErrInvalidPKCS7Data = errors.New("invalid PKCS7 data (empty or not padded)")

	// ErrInvalidPKCS7Padding indicates PKCS7 unpad fails to bad input.
	ErrInvalidPKCS7Padding = errors.New("invalid padding on input")
)

// pkcs7Pad right-pads the given byte slice with 1 to n bytes, where
// n is the block size. The size of the result is x times n, where x
// is at least 1.
func pkcs7Pad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if len(b) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	n := blocksize - (len(b) % blocksize)
	pb := make([]byte, len(b)+n)
	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))
	return pb, nil
}

// pkcs7Unpad validates and unpads data from the given bytes slice.
// The returned value will be 1 to n bytes smaller depending on the
// amount of padding, where n is the block size.
func pkcs7Unpad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if len(b) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	if len(b)%blocksize != 0 {
		return nil, ErrInvalidPKCS7Padding
	}
	c := b[len(b)-1]
	n := int(c)
	if n == 0 || n > len(b) {
		return nil, ErrInvalidPKCS7Padding
	}
	for i := 0; i < n; i++ {
		if b[len(b)-n+i] != c {
			return nil, ErrInvalidPKCS7Padding
		}
	}
	return b[:len(b)-n], nil
}

type aesCfbCipher struct {
	aesCipher cipher.Block
}

func newAesCfbCipher(b *testing.B, key []byte) aesCfbCipher {
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		b.Error(err)
	}

	return aesCfbCipher{
		aesCipher: aesCipher,
	}
}

func (aesCfbCipher aesCfbCipher) Encrypt(dst, nonce, plaintext, _additionalData []byte) []byte {
	encrypter := cipher.NewCFBEncrypter(aesCfbCipher.aesCipher, nonce)
	encrypter.XORKeyStream(dst, plaintext)
	return dst
}

func (aesCfbCipher aesCfbCipher) Decrypt(dst, nonce, ciphertext, additionalData []byte) {
	decrypter := cipher.NewCFBDecrypter(aesCfbCipher.aesCipher, nonce)
	decrypter.XORKeyStream(dst, ciphertext)
	// TODO
	// cipher.decrypter.CryptBlocks(dst, ciphertext)
	// pkcs7Unpad(dst, cipher.decrypter.BlockSize())
}
