package utils

import (
	crand "crypto/rand"
	"math/big"
)

func Random(n int) int {
	val, err := crand.Int(crand.Reader, big.NewInt(int64(n)))
	if err != nil {
		panic(err)
	}
	return int(val.Int64())
}
