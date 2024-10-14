package pkg

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
)

var FPix = []string{
	"%.0f%s",
	"%.1f%s",
	"%.2f%s",
	"%.3f%s",
	"%.4f%s",
	"%.5f%s",
	"%.6f%s",
}

var bytesSize = []string{"B", "kB", "MB", "GB", "TB", "PB", "EB"}
var iBytesSize = []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}

func logN(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

func humanizeBytes(s uint64, long int, base uint, sizes []string) string {
	if s < uint64(base) {
		return fmt.Sprintf("%d B", s)
	}

	num, _ := new(big.Float).SetPrec(base).SetString(strconv.FormatUint(s, 10))

	e := math.Floor(logN(float64(s), float64(base)))
	denominator := big.NewFloat(math.Pow(float64(base), e))

	suffix := sizes[int(e)]

	f := "%.2f %s"
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
func Bytes(s uint64, long int) string {
	return humanizeBytes(s, long, 1000, bytesSize)
}

// IBytes produces a humanizeBytes readable representation of an IEC size.
//
// See also: ParseBytes.
//
// IBytes(82854982) -> 79 MiB
// IBytes(82854982)  3 -> 79.017 MiB
// IBytes(82854982)  2 ->  78.98 MiB
// IBytes(82854982)  0 ->  79 MiB
func IBytes(s uint64, long int) string {
	return humanizeBytes(s, long, 1024, iBytesSize)
}
