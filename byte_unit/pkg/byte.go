package pkg

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

// IEC Sizes.
// kib is of bits
const (
	Byte = 1 << (iota * 10)
	KiByte
	MiByte
	GiByte
	TiByte
	PiByte
	EiByte
)

// SI Sizes.
const (
	IByte = 1
	KByte = IByte * 1000
	MByte = KByte * 1000
	GByte = MByte * 1000
	TByte = GByte * 1000
	PByte = TByte * 1000
	EByte = PByte * 1000
)

var bytesSizeTable = map[string]uint64{
	"b":   Byte,
	"kib": KiByte,
	"kb":  KByte,
	"mib": MiByte,
	"mb":  MByte,
	"gib": GiByte,
	"gb":  GByte,
	"tib": TiByte,
	"tb":  TByte,
	"pib": PiByte,
	"pb":  PByte,
	"eib": EiByte,
	"eb":  EByte,
	// Without suffix
	"":   Byte,
	"ki": KiByte,
	"k":  KByte,
	"mi": MiByte,
	"m":  MByte,
	"gi": GiByte,
	"g":  GByte,
	"ti": TiByte,
	"t":  TByte,
	"pi": PiByte,
	"p":  PByte,
	"ei": EiByte,
	"e":  EByte,
}

var FPix = []string{
	"%.0f %s",
	"%.1f %s",
	"%.2f %s",
	"%.3f %s",
	"%.4f %s",
	"%.5f %s",
	"%.6f %s",
	"%.7f %s",
}

func logN(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

func humanizeBytes(s uint64, long int, base float64, sizes []string) string {
	if s < 10 {
		return fmt.Sprintf("%d B", s)
	}
	e := math.Floor(logN(float64(s), base))
	suffix := sizes[int(e)]
	val := math.Floor(float64(s)/math.Pow(base, e)*10+0.5) / 10
	f := "%.2f %s"
	if long < 8 {
		f = FPix[long]
	}

	return fmt.Sprintf(f, val, suffix)
}

func makeFloatBytes(s uint64, long int, base float64, sizes []string) string {
	if s < 10 {
		return fmt.Sprintf("%d B", s)
	}
	e := math.Floor(logN(float64(s), base))
	suffix := sizes[int(e)]
	val := float64(s) / math.Pow(base, e)
	f := "%.2f %s"
	if long < 8 {
		f = FPix[long]
	}

	return fmt.Sprintf(f, val, suffix)
}

// BytesHalfAdjust produces a humanizeBytes readable representation of an SI size.
//
// See also: ParseBytes.
//
// Bytes(82854982) -> 83 MB
func BytesHalfAdjust(s uint64, long int) string {
	sizes := []string{"B", "kB", "MB", "GB", "TB", "PB", "EB"}
	return humanizeBytes(s, long, 1000, sizes)
}

// IBytesHalfAdjust produces a humanizeBytes readable representation of an IEC size.
//
// See also: ParseBytes.
//
// IBytes(82854982) -> 79 MiB
func IBytesHalfAdjust(s uint64, long int) string {
	sizes := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}
	return humanizeBytes(s, long, 1024, sizes)
}

// Bytes produces a humanizeBytes readable representation of an SI size.
//
// See also: ParseBytes.
//
// Bytes(82854982) -> 83 MB
func Bytes(s uint64, long int) string {
	sizes := []string{"B", "kB", "MB", "GB", "TB", "PB", "EB"}
	return makeFloatBytes(s, long, 1000, sizes)
}

// IBytes produces a humanizeBytes readable representation of an IEC size.
//
// See also: ParseBytes.
//
// IBytes(82854982) -> 79 MiB
func IBytes(s uint64, long int) string {
	sizes := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}
	return makeFloatBytes(s, long, 1024, sizes)
}

// ParseBytes parses a string representation of bytes into the number
// of bytes it represents.
//
// See Also: Bytes, IBytes.
//
// ParseBytes("42 MB") -> 42000000, nil
// ParseBytes("42 mib") -> 44040192, nil
func ParseBytes(s string) (uint64, error) {
	lastDigit := 0
	hasComma := false
	for _, r := range s {
		if !(unicode.IsDigit(r) || r == '.' || r == ',') {
			break
		}
		if r == ',' {
			hasComma = true
		}
		lastDigit++
	}

	num := s[:lastDigit]
	if hasComma {
		num = strings.Replace(num, ",", "", -1)
	}

	f, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return 0, err
	}

	extra := strings.ToLower(strings.TrimSpace(s[lastDigit:]))
	if m, ok := bytesSizeTable[extra]; ok {
		f *= float64(m)
		if f >= math.MaxUint64 {
			return 0, fmt.Errorf("too large: %v", s)
		}
		return uint64(f), nil
	}

	return 0, fmt.Errorf("unhandled size name: %v", extra)
}
