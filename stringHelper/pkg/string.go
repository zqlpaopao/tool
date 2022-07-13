package pkg

import (
	"reflect"
	"strconv"
	"time"
)

//StringFromAssertionFloat Convert to string type
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
	case uint:
		numberString = strconv.FormatInt(int64(floatOriginal), 10)
	case uint8:
		numberString = strconv.FormatInt(int64(floatOriginal), 10)
	case uint32:
		numberString = strconv.FormatInt(int64(floatOriginal), 10)
	case uint16:
		numberString = strconv.FormatInt(int64(floatOriginal), 10)
	case uint64:
		numberString = strconv.FormatInt(int64(floatOriginal), 10)
	case []uint8:
		numberString = string(floatOriginal)
		break
	case string:
		numberString = floatOriginal
	case time.Time:
		numberString = floatOriginal.String()
	}
	return numberString
}

//StringFromAssertionFloatNotZero 判断结构体的类型是不是指针类型
func StringFromAssertionFloatNotZero(ty string, vl reflect.Value) (b bool, numberString string) {

	switch ty {
	case "float", "float8", "float16", "float32", "float64":
		//numberString = strconv.FormatInt(int64(vl.Float()), 10)
		if vl.Interface() == 0 {
			return
		}
	case "int", "int8", "int16", "int32", "int64":
		if vl.Interface() == 0 {
			return
		}
		numberString = strconv.FormatInt(int64(vl.Int()), 10)
	case "string":
		if vl.String() == "" {
			return
		}
		b = true
	case "*string", "*int", "*int8", "*int16", "*int32", "*int64", "*float", "*float8", "*float16", "*float32", "*float64":
		//numberString = strconv.FormatInt(int64(vl.Float()), 10)
		//if vl.Elem().String()
		b = true
	}
	return
}
