package utils

import (
	"crypto/rand"
	"math/big"
	"net/url"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateSlug(length int) string {
	b := make([]byte, length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

func IsValidURL(toTest string) bool {
	u, err := url.ParseRequestURI(toTest)
	return err == nil && u.Scheme != "" && u.Host != ""
}
