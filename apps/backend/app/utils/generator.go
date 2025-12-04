package utils

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Daftar karakter untuk ID pendek
const shortIDChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func GenerateShortID(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = shortIDChars[seededRand.Intn(len(shortIDChars))]
	}
	return string(b)
}

const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateReferralCode(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func GenerateRandomNumber(length int) int {
	max := int(math.Pow10(length))
	return seededRand.Intn(max) + 1
}

func GenerateInvoiceCode(orderID string) string {
	now := time.Now()
	return fmt.Sprintf("INV-%02d%02d%04d%02d%02d%03d-%s",
		now.Day(), now.Month(), now.Year(),
		now.Hour(), now.Minute(), now.Second(), orderID)
}
