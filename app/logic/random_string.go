package logic

import (
	"math/rand"
	"strings"
)

var randLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func generateRandomString(length int, letters []rune) string {
	b := strings.Builder{}
	b.Grow(length)
	lc := len(letters)

	for i := 0; i < length; i++ {
		b.WriteRune(letters[rand.Intn(lc)])
	}
	return b.String()
}

func GenerateRandomString(length int) string {
	return generateRandomString(length, randLetters)
}
