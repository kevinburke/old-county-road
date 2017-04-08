package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"time"

	"golang.org/x/crypto/nacl/secretbox"
)

// FlashSuccess encrypts msg and sets it as a cookie on w. Only one success
// message can be set on w; the last call to FlashSuccess will be set on the
// response.
func FlashSuccess(w http.ResponseWriter, msg string, key *[32]byte) {
	setCookie(w, msg, "flash-success", key)
}

// FlashError encrypts msg and sets it as a cookie on w. Only one error can be
// set on w; the last call to FlashError will be set on the response.
func FlashError(w http.ResponseWriter, msg string, key *[32]byte) {
	setCookie(w, msg, "flash-error", key)
}

func setCookie(w http.ResponseWriter, msg string, name string, key *[32]byte) {
	c := &http.Cookie{
		Name:     name,
		Path:     "/",
		Value:    opaque(msg, key),
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}

// GetFlashSuccess finds a flash success message in the request (if one exists).
// If one exists then it's unset and returned.
func GetFlashSuccess(w http.ResponseWriter, r *http.Request, key *[32]byte) string {
	return getCookie(w, r, "flash-success", key)
}

// GetFlashError finds a flash error in the request (if one exists). If one
// exists then it's unset and returned.
func GetFlashError(w http.ResponseWriter, r *http.Request, key *[32]byte) string {
	return getCookie(w, r, "flash-error", key)
}

func getCookie(w http.ResponseWriter, r *http.Request, name string, key *[32]byte) string {
	cookie, err := r.Cookie(name)
	if err == http.ErrNoCookie {
		return ""
	}
	clearCookie(w, name)
	msg, err := unopaque(cookie.Value, key)
	if err != nil {
		return ""
	}
	return msg
}

func clearCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:    name,
		Path:    "/",
		MaxAge:  -1,
		Expires: time.Unix(1, 0),
	})
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

var errTooShort = errors.New("google_oauth_handler: Encrypted string is too short")
var errInvalidInput = errors.New("google_oauth_handler: Could not decrypt invalid input")

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
