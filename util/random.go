package util

import (
	"math/rand"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates random integer in rage of [min, max).
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min)
}

// randomString returns a string of the given length n.
func randomString(n int) string {
	l := len(alphabet)
	bs := make([]byte, n)

	for i := 0; i < n; i++ {
		bs[i] = alphabet[rand.Intn(l)]
	}

	return string(bs)
}

// RandomOwner returns random owner.
func RandomOwner() string {
	return randomString(8)
}

// RandomBalance returns random balance.
func RandomBalance() int64 {
	return RandomInt(0, 100000)
}

// RandomCurrency returns random currency code.
func RandomCurrency() string {
	currencies := []string{"USD", "GBP", "CAD", "EUR"}
	l := len(currencies)
	return currencies[rand.Intn(l)]
}

// RandomAmount returns random amount of money.
func RandomAmount() int64 {
	return RandomInt(1, 1000)
}
