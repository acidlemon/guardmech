package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type PayloadData interface {
	String() string
	Restore(data string) error
}

type SessionPayload struct {
	ExpireAt time.Time
	Data     PayloadData
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
