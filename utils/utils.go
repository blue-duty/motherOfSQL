package utils

import (
	"github.com/google/uuid"
	"math/rand"
	"time"
)

// GenIntValue gen int value randomly
func GenIntValue(len int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(len + 1)
}

// GenStringValue gen string value randomly
func GenStringValue(len int) string {
	l := GenIntValue(len)
	var str string
	for i := 0; i < l; i++ {
		str += string(rune(rand.Intn(26) + 97))
	}
	return str
}

// GenDatetimeValue gen datetime value randomly
func GenDatetimeValue() string {
	return "2020-01-01 00:00:00"
}

// GenFloatValue gen float value randomly
func GenFloatValue() float64 {
	return rand.Float64()
}

// GenChineseValue gen simple chinese characters randomly
func GenChineseValue(len int) string {
	l := GenIntValue(len + 1)
	var str string
	for i := 0; i < l; i++ {
		str += string(rune(rand.Intn(20902) + 19968))
	}
	return str
}

// UUID generate uuid
func UUID() string {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	return newUUID.String()
}

// TakeParentheses Take the contents in parentheses from the string
func TakeParentheses(str string) string {
	start := 0
	end := 0
	for i, v := range str {
		if v == '(' {
			start = i
		}
		if v == ')' {
			end = i
		}
	}
	return str[start+1 : end]
}
