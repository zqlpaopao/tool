package pkg

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

//GetRandInt rand num
func GetRandInt(n int) int {
	return rand.Intn(n)
}

const (
	RandSourceLetter                   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	RandSourceNumber                   = "0123456789"
	RandSourceLetterAndNumber          = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	RandSourceUppercaseLetterAndNumber = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)


//RandGenString src string,len string
func RandGenString(source string, randLen int64) string {
	letterRunes := []rune(source)
	randRe := make([]rune, randLen)
	for i := range randRe {
		randRe[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(randRe)
}

//RandGenStringForUUID src string,len string
func RandGenStringForUUID(source string, randLen int64) string {
	letterRunes := []rune(source)
	randRe := make([]rune, randLen)
	date := time.Now().Format("20060102")
	start := 0
	offset := 1
	for i := range randRe {
		if len(date) > start && i%2 != 0 {
			randRe[i] = []rune(date[start:offset])[0]
			start++
			offset++
			continue
		}
		randRe[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(randRe)
}

//RandParseStringForUUID parse uuid timeData
func RandParseStringForUUID(uuid string) (date string) {
	var res []rune
	if len(uuid) == 34 {
		start := 0
		offset := 1
		temp := []rune(uuid)
		for i := range temp {
			if len(res) < 8 && i%2 != 0 {
				res = append(res, temp[i])
				start++
				offset++
			}
		}
	}
	return string(res)
}

