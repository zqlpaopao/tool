package pkg

import "strconv"

func StringFromAssertionFloat(number interface{}) string {
	var numberString string
	switch floatOriginal := number.(type) {
	case float64:
		numberString = strconv.FormatInt(int64(floatOriginal), 10)
	case float32:
		numberString = strconv.FormatInt(int64(floatOriginal), 10)
	case int:
		numberString = strconv.FormatInt(int64(floatOriginal), 10)
	case int8:
		numberString = strconv.FormatInt(int64(floatOriginal), 10)
	case int32:
		numberString = strconv.FormatInt(int64(floatOriginal), 10)
	case int16:
		numberString = strconv.FormatInt(int64(floatOriginal), 10)
	case int64:
		numberString = strconv.FormatInt(floatOriginal, 10)
	case []uint8:
		numberString = string(floatOriginal)
		break
	case string:
		numberString = floatOriginal
	}
	return numberString
}

