// Package xwing implements the hybrid quantum-resistant key encapsulation
// method X-Wing, which combines X25519, ML-KEM-768, and SHA3-256 as specified
// in [draft-connolly-cfrg-xwing-kem].
//
// adapted from https://github.com/FiloSottile/mlkem768/blob/main/xwing/xwing.go
//
// [draft-connolly-cfrg-xwing-kem]: https://www.ietf.org/archive/id/draft-connolly-cfrg-xwing-kem-05.html
package xwing

import (
	"bytes"
	"crypto/ecdh"
	"crypto/mlkem"
	"crypto/rand"
	"errors"

	"crypto/sha3"
)

const (
	CiphertextSize       = mlkem.CiphertextSize768 + 32
	EncapsulationKeySize = mlkem.EncapsulationKeySize768 + 32
	SharedKeySize        = 32
	SeedSize             = 32
)

// A DecapsulationKey is the secret key used to decapsulate a shared key from a
// ciphertext. It includes various precomputed values.
type DecapsulationKey struct {
	sk  [SeedSize]byte
	skM *mlkem.DecapsulationKey768
	skX *ecdh.PrivateKey
	pk  [EncapsulationKeySize]byte
}

type EncapsulationKey struct {
	mlKemEncapsulationKey *mlkem.EncapsulationKey768
	x25519Key             *ecdh.PublicKey
}

// Bytes returns the decapsulation key as a 32-byte seed.
func (dk *DecapsulationKey) Bytes() []byte {
	return bytes.Clone(dk.sk[:])
}

// EncapsulationKey returns the public encapsulation key necessary to produce
// ciphertexts.
func (dk *DecapsulationKey) EncapsulationKey() *EncapsulationKey {
	pkMBytes := dk.pk[:mlkem.EncapsulationKeySize768]
	pkX := dk.pk[mlkem.EncapsulationKeySize768:]

	peerKey, err := ecdh.X25519().NewPublicKey(pkX)
	if err != nil {
		panic(err)
	}

	pkM, err := mlkem.NewEncapsulationKey768(pkMBytes)
	if err != nil {
		panic(err)
	}

	return &EncapsulationKey{
		mlKemEncapsulationKey: pkM,
		x25519Key:             peerKey,
	}
}

// GenerateKey generates a new decapsulation key, drawing random bytes from
// crypto/rand. The decapsulation key must be kept secret.
func GenerateKey() (*DecapsulationKey, error) {
	sk := make([]byte, SeedSize)
	if _, err := rand.Read(sk); err != nil {
		return nil, err
	}
	return NewKeyFromSeed(sk)
}

// NewKeyFromSeed deterministically generates a decapsulation key from a 32-byte
// seed. The seed must be uniformly random.
func NewKeyFromSeed(sk []byte) (*DecapsulationKey, error) {
	if len(sk) != SeedSize {
		return nil, errors.New("xwing: invalid seed length")
	}

	s := sha3.NewSHAKE256()
	s.Write(sk)
	expanded := make([]byte, mlkem.SeedSize+32)
	if _, err := s.Read(expanded); err != nil {
		return nil, err
	}

	skM, err := mlkem.NewDecapsulationKey768(expanded[:mlkem.SeedSize])
	if err != nil {
		return nil, err
	}
	pkM := skM.EncapsulationKey()

	skX := expanded[mlkem.SeedSize:]
	x, err := ecdh.X25519().NewPrivateKey(skX)
	if err != nil {
		return nil, err
	}
	pkX := x.PublicKey().Bytes()

	dk := &DecapsulationKey{}
	copy(dk.sk[:], sk)
	dk.skM = skM
	dk.skX = x
	copy(dk.pk[:], append(pkM.Bytes(), pkX...))
	return dk, nil
}

const xwingLabel = (`` +
	`\./` +
	`/^\`)

func combiner(ssM, ssX, ctX, pkX []byte) []byte {
	h := sha3.New256()
	h.Write(ssM)
	h.Write(ssX)
	h.Write(ctX)
	h.Write(pkX)
	h.Write([]byte(xwingLabel))
	return h.Sum(nil)
}

// Encapsulate generates a shared key and an associated ciphertext from an
// encapsulation key, drawing random bytes from crypto/rand.
// If the encapsulation key is not valid, Encapsulate returns an error.
//
// The shared key must be kept secret.
func (encapsulationKey *EncapsulationKey) Encapsulate() (ciphertext, sharedKey []byte) {
	ephemeralKey, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	ctX := ephemeralKey.PublicKey().Bytes()
	ssX, err := ephemeralKey.ECDH(encapsulationKey.x25519Key)
	if err != nil {
		panic(err)
	}

	ctM, ssM := encapsulationKey.mlKemEncapsulationKey.Encapsulate()

	ss := combiner(ssM, ssX, ctX, encapsulationKey.x25519Key.Bytes())
	ct := append(ctM, ctX...)
	return ct, ss
}

// Decapsulate generates a shared key from a ciphertext and a decapsulation key.
// If the ciphertext is not valid, Decapsulate returns an error.
//
// The shared key must be kept secret.
func (dk *DecapsulationKey) Decapsulate(ciphertext []byte) (sharedKey []byte, err error) {
	if len(ciphertext) != CiphertextSize {
		return nil, errors.New("xwing: invalid ciphertext length")
	}

	ctM := ciphertext[:mlkem.CiphertextSize768]
	ctX := ciphertext[mlkem.CiphertextSize768:]
	pkX := dk.pk[mlkem.EncapsulationKeySize768:]

	ssM, err := dk.skM.Decapsulate(ctM)
	if err != nil {
		return nil, err
	}

	peerKey, err := ecdh.X25519().NewPublicKey(ctX)
	if err != nil {
		return nil, err
	}
	ssX, err := dk.skX.ECDH(peerKey)
	if err != nil {
		return nil, err
	}

	ss := combiner(ssM, ssX, ctX, pkX)
	return ss, nil
}
