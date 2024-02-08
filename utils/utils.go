package utils

import (
	"fmt"
	"strconv"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func ConvertHexToBinary(hex string) string {
	val, err := strconv.ParseUint(hex, 16, 64)
	Check(err)
	binaryStr := strconv.FormatUint(val, 2)
	return fmt.Sprintf("%064s", binaryStr)
}

func ConvertStringToRune(str string) rune {
	return []rune(str)[0]
}

func ConvertStringToUint(str string) uint64 {
	val, err := strconv.ParseUint(str, 10, 64)
	Check(err)
	return val
}
