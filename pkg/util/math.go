package util

import (
	"fmt"
	"math"
	"strconv"
)

// MinInt returns the smaller of x or y.
func MinInt(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0
	}

	return value
}
