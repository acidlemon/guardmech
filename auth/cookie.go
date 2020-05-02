package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// type Config struct {
// 	CryptoKey    string
// 	SignatureKey string
// }

type PayloadData interface {
	String() string
	Restore(data string) error
}

type SessionPayload struct {
	ExpireAt time.Time
	Data     PayloadData
	// cfg      Config
}

func RestoreSessionPayload(fromCookie string, box PayloadData) (*SessionPayload, error) {
	ss := strings.Split(fromCookie, "('-'*)")
	data := ss[0]
	signature := ss[1]

	// validate signature
	signKey := os.Getenv("GUARDMECH_SIGNATURE_KEY")
	sig := &signator{SignKey: []byte(signKey)}
	result, err := sig.Verify(data, signature)
	if err != nil {
		return nil, errors.Wrap(err, "Session Signature Verification failed")
	}
	if !result {
		return nil, fmt.Errorf("Session Sigunature does not matched")
	}

	// parse data
	ss = strings.Split(data, "(#'-')")
	encVal := ss[0]
	expireUnix, err := strconv.ParseInt(ss[1], 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse ExpireAt epoch")
	}
	expireAt := time.Unix(expireUnix, 0)

	// decrypt data
	cryptoKey := os.Getenv("GUARDMECH_CRYPTO_KEY")
	c := cryptor{
		CipherKey: []byte(cryptoKey),
	}
	sessionData, err := c.Decrypt(encVal)
	if err != nil {
		return nil, errors.Wrap(err, "Decrypt Error")
	}

	err = box.Restore(sessionData)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to restore")
	}

	return &SessionPayload{
		ExpireAt: expireAt,
		Data:     box,
	}, nil
}

func (s *SessionPayload) sign(data string) string {
	key := os.Getenv("GUARDMECH_SIGNATURE_KEY")
	sig := &signator{SignKey: []byte(key)}
	return sig.Sign(data)
}

func (s *SessionPayload) String() string {
	key := os.Getenv("GUARDMECH_CRYPTO_KEY")
	c := &cryptor{
		CipherKey: []byte(key),
	}

	encVal, err := c.Encrypt(s.Data.String())
	if err != nil {
		log.Println("encrypt error: ", err)
		return ""
	}

	data := fmt.Sprintf("%s(#'-')%d", encVal, s.ExpireAt.Unix())
	signature := s.sign(data)

	return fmt.Sprintf("%s('-'*)%s", data, signature)
}

func (s *SessionPayload) MakeCookie(req *http.Request, key string, extend time.Duration) *http.Cookie {
	domain := req.URL.Host
	value := s.String()
	return makeCookie(domain, key, value, s.ExpireAt.Add(extend))
}
func (s *SessionPayload) RevokeCookie(domain, key string) *http.Cookie {
	return revokeCookie(domain, key)
}

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
