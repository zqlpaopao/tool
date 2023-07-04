package pkg

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

var timeSizes = []string{"ns", "us", "ms", "s"}

//

// humanizeTimes
//
// 82 854 982 -> 82.855 ms
func humanizeTimes(s uint64, long int, base uint) string {
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

// halfAdjustTimes
//
// 82 854 982 -> 82.855 ms
func halfAdjustTimes(s uint64, long int, base uint) string {
	if s < uint64(base) {
		return fmt.Sprintf("%d ns", s)
	}
	e := math.Floor(logN(float64(s), float64(base)))
	suffix := timeSizes[int(e)]
	return Float2String(float64(s)/math.Pow(float64(base), e), long, suffix)

}

func Float2String(float642 float64, long int, suffix string) string {

	v := strconv.FormatFloat(float642, 'f', -1, 64)

	s := strings.Split(v, ".")

	f := "%.0f %s"
	if long < 4 {
		f = FPix[long]
	}
	if len(s) < 2 {
		return fmt.Sprintf(f, float642, suffix)
	}
	if len(s[1]) == long {
		return fmt.Sprintf(f, float642, suffix)
	}
	for i := 0; i < len(s[1]); i++ {
		if i == long && s[1][i] == byte('5') {
			return fmt.Sprintf(f, math.Floor(float642*math.Pow(10, float64(long))+5/math.Pow(10, float64(long-1)))/math.Pow(10, float64(long)), suffix)
		} else if i == long && s[1][i] != byte('5') {
			return fmt.Sprintf(f, float642, suffix)
		}
	}
	return fmt.Sprintf(f, float642, suffix)
}

// Times produces a humanizeBytes readable representation of an SI size.
//
// Duration(82854982) -> 82.855 ms
// Duration(82814982) -> 82.81 ms
func Times(s uint64, long int) string {
	return humanizeTimes(s, long, 1000)
}

func TimesHalfAdjust(s uint64, long int) string {
	return halfAdjustTimes(s, long, 1000)
}
