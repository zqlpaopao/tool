package pkg


import (
	"fmt"
	"math"
	"math/big"
	"strconv"
)

// var bytesSizeBps = []string{"BPS", "KBPS", "MBPS", "GBPS", "TBPS", "PBPS", "EBPS"}
var bytesSizeBps = []string{"B", "K", "M", "G", "T", "P", "E"}

func humanizeBytesBps(s float64, long int, base uint, sizes []string) string {
	if uint(s) < base {
		return fmt.Sprintf("%f B", s)
	}

	num, _ := new(big.Float).SetPrec(base).SetString(strconv.FormatFloat(s, 'f', long+2, 64))

	e := math.Floor(logN(s, float64(base)))
	denominator := big.NewFloat(math.Pow(float64(base), e))

	suffix := sizes[int(e)]

	f := "%.2f%s"
	if long < 7 {
		f = FPix[long]
	}

	denominator = num.Quo(num, denominator)

	return fmt.Sprintf(f, denominator, suffix)
}


// Bytes produces a humanizeBytes readable representation of an SI size.
//
// See also: ParseBytes.
//
// Bytes(82854982) 3 -> 82.855 MB
// Bytes(82854982) 2 -> 82.85 MB
// Bytes(82854982) 0 -> 83 MB
func BytesBps(s float64, long int) string {
	return humanizeBytesBps(s, long, 1000, bytesSizeBps)
}
