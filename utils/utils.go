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

func ConvertStringToUint(str string) uint {
	val, err := strconv.ParseUint(str, 10, 64)
	Check(err)
	return uint(val)
}

func ConvertBinaryToUint(binary string) uint {
	val, err := strconv.ParseUint(binary, 2, 64)
	Check(err)
	return uint(val)
}

func GetMemoryInfo(tagSize uint, indexSize uint, kind string, address string) (uint, uint, uint) {
	tagBin := address[:tagSize]
	var index uint
	var tag uint

	if kind == "full" {
		index = 0
	} else {
		indexBin := address[tagSize : tagSize+indexSize]
		index = ConvertBinaryToUint(indexBin)
	}

	tag = ConvertBinaryToUint(tagBin)
	offsetBin := address[tagSize+indexSize:]
	offset := ConvertBinaryToUint(offsetBin)
	return index, tag, offset
}
