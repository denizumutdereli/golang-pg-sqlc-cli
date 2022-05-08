package util

import (
	"math/rand"
	"strings"
	"time"
)

const alpabet = "abcdefghijklmnoprstuvyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Random int min-max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1) //min->max
}

//Random String length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alpabet)

	for i := 0; i < n; i++ {
		c := alpabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()
}

//Random owner name
func RandomOwner() string {
	owner := RandomString(rand.Intn(20) + 1)

	// if owner[0:1] == "ğ" {
	// 	return RandomOwner()
	// } else {
	// 	return owner
	// }

	return owner

}

//Random Money
func RandomBalance() int64 {
	return RandomInt(0, 1000)
}

//Random Currency
func RandomCurrency() string {
	currencies := []string{"BTC", "USDT"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
