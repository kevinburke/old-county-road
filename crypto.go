package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"

	"golang.org/x/crypto/nacl/secretbox"
)

// NewRandomKey returns a random key or panics if one cannot be provided.
func NewRandomKey() *[32]byte {
	key := new([32]byte)
	if _, err := io.ReadFull(rand.Reader, key[:]); err != nil {
		panic(err)
	}
	return key
}

// getSecretKey produces a valid [32]byte secret key or returns an error. If
// hexKey is the empty string, a valid 32 byte key will be randomly generated
// and returned. If hexKey is invalid, an error is returned.
func getSecretKey(hexKey string) (*[32]byte, error) {
	if hexKey == "" {
		return NewRandomKey(), nil
	}

	if len(hexKey) != 64 {
		return nil, errWrongLength
	}
	secretKeyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, err
	}
	secretKey := new([32]byte)
	copy(secretKey[:], secretKeyBytes)
	return secretKey, nil
}

func newNonce() *[24]byte {
	nonce := new([24]byte)
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		panic(err)
	}
	return nonce
}

func opaqueByte(b []byte, secretKey *[32]byte) string {
	nonce := newNonce()
	encrypted := secretbox.Seal(nonce[:], b, nonce, secretKey)
	return base64.URLEncoding.EncodeToString(encrypted)
}

var errTooShort = errors.New("Encrypted string is too short")
var errInvalidInput = errors.New("Could not decrypt invalid input")

func unopaqueByte(compressed string, secretKey *[32]byte) ([]byte, error) {
	encrypted, err := base64.URLEncoding.DecodeString(compressed)
	if err != nil {
		return nil, err
	}
	if len(encrypted) < 24 {
		return nil, errTooShort
	}
	decryptNonce := new([24]byte)
	copy(decryptNonce[:], encrypted[:24])
	decrypted, ok := secretbox.Open([]byte{}, encrypted[24:], decryptNonce, secretKey)
	if !ok {
		return nil, errInvalidInput
	}
	return decrypted, nil
}

// Opaque encrypts s with secretKey and returns the encrypted string encoded
// with base64, or an error.
func opaque(s string, secretKey *[32]byte) string {
	return opaqueByte([]byte(s), secretKey)
}

// Unopaque decodes compressed using base64, then decrypts the decoded byte
// array using the secretKey.
func unopaque(compressed string, secretKey *[32]byte) (string, error) {
	b, err := unopaqueByte(compressed, secretKey)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
