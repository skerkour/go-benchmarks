package encryption_unauthenticated

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"testing"

	"github.com/bloom42/stdx/crypto/chacha20"
	chacha12 "github.com/bloom42/stdx/crypto/xchacha12"
	"github.com/bloom42/stdx/crypto/xchacha20"
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

type Cipher interface {
	Encrypt(dst, nonce, plaintext []byte) []byte
	Decrypt(dst, nonce, ciphertext []byte)
}

func BenchmarkEncryptUnauthenticated(b *testing.B) {
	xChaCha20Key := utils.RandBytes(b, xchacha20.KeySize)
	xChaCha20Nonce := utils.RandBytes(b, xchacha20.NonceSize)

	chaCha20Key := utils.RandBytes(b, chacha20.KeySize)
	chaCha20Nonce := utils.RandBytes(b, chacha20.NonceSize)
	aes256CbcKey := utils.RandBytes(b, 32)
	aes256CbcIv := utils.RandBytes(b, 16)

	aes256CfbKey := utils.RandBytes(b, 32)
	aes256CfbIv := utils.RandBytes(b, 16)

	for _, size := range BENCHMARKS {
		benchmarkEncrypt(b, size, "XChaCha20", newXChaCha20Cipher(b, xChaCha20Key, xChaCha20Nonce), xChaCha20Nonce)
		benchmarkEncrypt(b, size, "XChaCha12", newXChaCha12Cipher(b, xChaCha20Key, xChaCha20Nonce), xChaCha20Nonce)
		benchmarkEncrypt(b, size, "ChaCha20", newChaCha20Cipher(b, chaCha20Key, chaCha20Nonce), chaCha20Nonce)
		benchmarkEncrypt(b, size, "AES_256_CBC", newAesCbcCipher(b, aes256CbcKey), aes256CbcIv)
		benchmarkEncrypt(b, size, "AES_256_CFB", newAesCfbCipher(b, aes256CfbKey), aes256CfbIv)
	}
}

func BenchmarkDecryptUnauthenticated(b *testing.B) {
	xChaCha20Key := utils.RandBytes(b, xchacha20.KeySize)
	xChaCha20Nonce := utils.RandBytes(b, xchacha20.NonceSize)

	chaCha20Key := utils.RandBytes(b, chacha20poly1305.KeySize)
	chaCha20Nonce := utils.RandBytes(b, chacha20poly1305.NonceSize)

	aes256CbcKey := utils.RandBytes(b, 32)
	aes256CbcIv := utils.RandBytes(b, 16)

	aes256CfbKey := utils.RandBytes(b, 32)
	aes256CfbIv := utils.RandBytes(b, 16)

	for _, size := range BENCHMARKS {
		benchmarkEncrypt(b, size, "XChaCha20", newXChaCha20Cipher(b, xChaCha20Key, xChaCha20Nonce), xChaCha20Nonce)
		benchmarkEncrypt(b, size, "XChaCha12", newXChaCha12Cipher(b, xChaCha20Key, xChaCha20Nonce), xChaCha20Nonce)
		benchmarkDecrypt(b, size, "ChaCha20", newChaCha20Cipher(b, chaCha20Key, chaCha20Nonce), chaCha20Nonce)
		benchmarkDecrypt(b, size, "AES_256_CBC", newAesCbcCipher(b, aes256CbcKey), aes256CbcIv)
		benchmarkDecrypt(b, size, "AES_256_CFB", newAesCfbCipher(b, aes256CfbKey), aes256CfbIv)
	}
}

func benchmarkEncrypt[C Cipher](b *testing.B, size int64, algorithm string, cipher C, nonce []byte) {
	b.Run(fmt.Sprintf("%s-%s", utils.BytesCount(size), algorithm), func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(size)
		plaintext := utils.RandBytes(b, size)
		dst := make([]byte, len(plaintext)+64)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cipher.Encrypt(dst, nonce, plaintext)
		}
	})
}

func benchmarkDecrypt[C Cipher](b *testing.B, size int64, algorithm string, cipher C, nonce []byte) {
	b.Run(fmt.Sprintf("%s-%s", utils.BytesCount(size), algorithm), func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(size)
		plaintext := utils.RandBytes(b, size)
		cipherText := make([]byte, len(plaintext)+512)
		cipherText = cipher.Encrypt(cipherText, nonce, plaintext)
		dst := make([]byte, len(cipherText))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cipher.Decrypt(dst, nonce, cipherText)
		}
	})
}

type xChaCha20Cipher struct {
	cipher cipher.Stream
}

func newXChaCha20Cipher(b *testing.B, key, nonce []byte) xChaCha20Cipher {
	cipher, err := xchacha20.New(key, nonce)
	if err != nil {
		b.Error(err)
	}

	return xChaCha20Cipher{
		cipher: cipher,
	}
}

func (cipher xChaCha20Cipher) Encrypt(dst, _nonce, plaintext []byte) []byte {
	cipher.cipher.XORKeyStream(dst, plaintext)
	return dst
}

func (cipher xChaCha20Cipher) Decrypt(dst, _nonce, ciphertext []byte) {
	cipher.cipher.XORKeyStream(dst, ciphertext)
}

type xChaCha12Cipher struct {
	cipher cipher.Stream
}

func newXChaCha12Cipher(b *testing.B, key, nonce []byte) xChaCha12Cipher {
	cipher, err := chacha12.NewCipher(key, nonce)
	if err != nil {
		b.Error(err)
	}

	return xChaCha12Cipher{
		cipher: cipher,
	}
}

func (cipher xChaCha12Cipher) Encrypt(dst, _nonce, plaintext []byte) []byte {
	cipher.cipher.XORKeyStream(dst, plaintext)
	return dst
}

func (cipher xChaCha12Cipher) Decrypt(dst, _nonce, ciphertext []byte) {
	cipher.cipher.XORKeyStream(dst, ciphertext)
}

type chaCha20Cipher struct {
	cipher cipher.Stream
}

func newChaCha20Cipher(b *testing.B, key, nonce []byte) chaCha20Cipher {
	cipher, err := chacha20.NewCipher(key, nonce)
	if err != nil {
		b.Error(err)
	}

	return chaCha20Cipher{
		cipher: cipher,
	}
}

func (cipher chaCha20Cipher) Encrypt(dst, _nonce, plaintext []byte) []byte {
	cipher.cipher.XORKeyStream(dst, plaintext)
	return dst
}

func (cipher chaCha20Cipher) Decrypt(dst, _nonce, ciphertext []byte) {
	cipher.cipher.XORKeyStream(dst, ciphertext)
}

type aesCbcCipher struct {
	cipher cipher.Block
}

func newAesCbcCipher(b *testing.B, key []byte) aesCbcCipher {
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		b.Error(err)
	}

	return aesCbcCipher{
		cipher: aesCipher,
	}
}

func (aesCbcCipher aesCbcCipher) Encrypt(dst, nonce, plaintext []byte) []byte {
	encrypter := cipher.NewCBCEncrypter(aesCbcCipher.cipher, nonce)
	paddedPlaintext, _ := pkcs7Pad(plaintext, encrypter.BlockSize())
	encrypter.CryptBlocks(dst, paddedPlaintext)
	return dst
}

func (aesCbcCipher aesCbcCipher) Decrypt(dst, nonce, ciphertext []byte) {
	decrypter := cipher.NewCBCDecrypter(aesCbcCipher.cipher, nonce)
	decrypter.CryptBlocks(dst, ciphertext)
	// unpaddedCiphertext, err := pkcs7Unpad(dst, decrypter.BlockSize())
	// if err != nil {
	// 	panic(err)
	// }
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

func (aesCfbCipher aesCfbCipher) Encrypt(dst, nonce, plaintext []byte) []byte {
	encrypter := cipher.NewCFBEncrypter(aesCfbCipher.aesCipher, nonce)
	encrypter.XORKeyStream(dst, plaintext)
	return dst
}

func (aesCfbCipher aesCfbCipher) Decrypt(dst, nonce, ciphertext []byte) {
	decrypter := cipher.NewCFBDecrypter(aesCfbCipher.aesCipher, nonce)
	decrypter.XORKeyStream(dst, ciphertext)
}
