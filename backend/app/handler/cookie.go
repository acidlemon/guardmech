package handler

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"time"
)

type signator struct {
	SignKey []byte
}

func (s *signator) calc(in string) []byte {
	mac := hmac.New(sha256.New, s.SignKey)
	mac.Write([]byte(in))
	return mac.Sum(nil)
}

func (s *signator) Sign(in string) string {
	return base64.StdEncoding.EncodeToString(s.calc(in))
}

func (s *signator) Verify(in, b64str string) (bool, error) {
	data, err := base64.StdEncoding.DecodeString(b64str)
	if err != nil {
		return false, err
	}

	hash := s.calc(in)
	return hmac.Equal(data, hash), nil
}

type cryptor struct {
	CipherKey []byte
}

// ref. https://gist.github.com/kkirsche/e28da6754c39d5e7ea10
func (c *cryptor) Encrypt(s string) (string, error) {
	// using AES-GCM Encryption
	block, err := aes.NewCipher(c.CipherKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nil, nonce, []byte(s), nil)
	cipherText = append(nonce, cipherText...)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func (c *cryptor) Decrypt(b64str string) (string, error) {
	buf, err := base64.StdEncoding.DecodeString(b64str)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(c.CipherKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := buf[:gcm.NonceSize()]
	plainByte, err := gcm.Open(nil, nonce, buf[gcm.NonceSize():], nil)
	if err != nil {
		return "", err
	}

	return string(plainByte), nil
}

func revokeCookie(domain, key string) *http.Cookie {
	return makeCookie(domain, key, "", time.Now().Add(-1*time.Hour))
}

func makeCookie(domain, key, value string, expireAt time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     key,
		Value:    value,
		Path:     "/",
		Domain:   domain,
		HttpOnly: true,
		Secure:   false, // TODO
		Expires:  expireAt,
	}
}
