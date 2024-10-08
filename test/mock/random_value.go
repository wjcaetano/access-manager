package mock

import (
	"crypto/rand"
	"math/big"
	"strconv"
	"time"
)

func RandomInt64() int64 {
	n, _ := rand.Int(rand.Reader, big.NewInt(99999))
	return n.Int64()
}

func RandomInt() int {
	return int(RandomInt64())
}

func RandomBool() bool {
	return RandomInt() > 0
}

func RandomUint() uint {
	return uint(RandomInt())
}

func RandomUint64() uint64 {
	return uint64(RandomInt64())
}

func RandomString() string {
	return strconv.Itoa(RandomInt())
}

func RandomFloat64() float64 {
	return float64(RandomInt())
}

func RandomDate() time.Time {
	return time.Now().UTC().AddDate(0, 0, RandomInt())
}
