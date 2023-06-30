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

	f1 := float64(s) / math.Pow(float64(base), e)
	fmt.Println(f1)
	return fmt.Sprintf("%s %s", Round(f1, long), suffix)
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

func Round(s float64, long int) (str string) {
	st := strconv.FormatFloat(s, 'f', 5, 64)
	vl := strings.Split(st, ".")
	if len(vl) < 2 {
		return st
	}
	str = vl[0] + "."
	if len(vl[1]) < long+1 {
		return st
	}
	for i := 1; i < len(vl[1]); i++ {
		if i == long {
			if vl[1][i] >=  byte('5') {
				str += string(vl[1][i-1] + 1)
				return
			}
			str += string(vl[1][i-1])
			return
		}
		str += string(vl[1][i-1])

	}
	return st

}
