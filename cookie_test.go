package guardmech

import (
	"encoding/base64"
	"encoding/hex"
	"testing"
)

func TestCryptor(t *testing.T) {
	c := &Cryptor{[]byte("12345678901234567890123456789012")}

	enc, err := c.Encrypt("hoge")
	if err != nil {
		t.Error("encrypt error", err)
	}

	dec, err := c.Decrypt(enc)
	if err != nil {
		t.Error("decrypt error", err)
	}

	if dec != "hoge" {
		t.Error("Decrypt(enc) != original text")
	}
}

func TestSignator(t *testing.T) {

	// check HMAC-SHA256 using https://www.freeformatter.com/hmac-generator.html
	text := `(o'-') < secret text`
	precomputed := `50d82b3a68204da83f89877eff0208fc5ee0f607bead90306384c06a53bc71a9`
	secret := `09876543210w0`

	s := &Signator{[]byte(secret)}

	signed := s.Sign(text)
	buf, err := base64.StdEncoding.DecodeString(signed)
	if err != nil {
		t.Errorf(`signed string is not base64`)
	}
	if hex.EncodeToString(buf) != precomputed {
		t.Errorf(`signed = "%s" != precomputed "%s"  \n`, signed, precomputed)
	}

	buf, err = hex.DecodeString(precomputed)
	if err != nil {
		t.Fatal(`invalid hex string`)
	}
	res, err := s.Verify(text, base64.StdEncoding.EncodeToString(buf))
	if err != nil {
		t.Error("verification error", err)
	}
	if !res {
		t.Error("verification failed")
	}

}
