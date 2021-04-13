package internal

import (
	"math"
)

func Pow(a, n int) int32 {
	return int32(math.Pow(float64(a), float64(n)))
}
