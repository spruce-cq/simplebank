package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt32 generate a random Integer32 between min and max
func RandomInt32(min, max int32) int32 {
	return min + rand.Int31n(max-min+1)
}

// RandomString generate a random string with the given length
func RandomString(n int) string {
	var sb strings.Builder
	length := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(length)]
		sb.WriteByte(c)
	}
	return sb.String()
}

// RandomOwner generate a random owner string for testing
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generate a random integer(32) less than 100
func RandomMoney() int32 {
	return rand.Int31n(100)
}

// RandomCurrency generate a random currency from ["EUR", "USD", "CAD"]
func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}
	return currencies[rand.Intn(len(currencies))]
}
