package pkg

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
)

var timeSizes = []string{"ns", "us", "ms", "s"}

// humanizeTimes
//
// 82 854 982 -> 82.855 ms
func humanizeTimes(s uint64, long int, base uint, sizes []string) string {
	if s < uint64(base) {
		return fmt.Sprintf("%d ns", s)
	}

	num, _ := new(big.Float).SetPrec(base).SetString(strconv.FormatUint(s, 10))

	e := math.Floor(logN(float64(s), float64(base)))
	denominator := big.NewFloat(math.Pow(float64(base), e))

	suffix := timeSizes[int(e)]

	f := "%.2f %s"
	if long < 4 {
		f = FPix[long]
	}

	denominator = num.Quo(num, denominator)

	return fmt.Sprintf(f, denominator, suffix)
}

// Times produces a humanizeBytes readable representation of an SI size.
//
// Duration(82854982) -> 82.855 ms
// Duration(82814982) -> 82.81 ms
func Times(s uint64, long int) string {
	return humanizeTimes(s, long, 1000, timeSizes)
}
