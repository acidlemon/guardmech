package logic

import (
	"crypto/rand"
	"strings"
)

// GenerateRandomString returns a cryptographically secure random string of the given length.
// The result consists of characters from the crypto/rand.Text alphabet (A-Z and 2-7).
func GenerateRandomString(length int) string {
	b := strings.Builder{}
	b.Grow(length)
	for b.Len() < length {
		b.WriteString(rand.Text())
	}
	return b.String()[:length]
}
