package utils

import (
	"math/rand"
	"time"
)

const (
	digitsCharset = "0123456789"
	charset       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" + digitsCharset
)

var (
	seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// RandStringWithCharset generate a random string with a specific charset
func RandStringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// RandString generate a random string with the default charset ([A-Za-z])
func RandString(length int) string {
	return RandStringWithCharset(length, charset)
}

// RandDigitString generate a random string with the default digits charset ([0-9])
func RandDigitString(length int) string {
	return RandStringWithCharset(length, digitsCharset)
}
