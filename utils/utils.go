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
	val, err := strconv.ParseInt(hex, 16, 64)
	Check(err)
	binaryStr := strconv.FormatInt(val, 2)
	return fmt.Sprintf("%064s", binaryStr)
}

func ConvertHexToInt(hex string) int {
	val, err := strconv.ParseInt(hex, 16, 64)
	Check(err)
	return int(val)
}

func ConvertStringToRune(str string) rune {
	return []rune(str)[0]
}

func ConvertStringToInt(str string) int {
	val, err := strconv.ParseInt(str, 10, 64)
	Check(err)
	return int(val)
}

func ConvertBinaryToInt(binary string) int {
	val, err := strconv.ParseInt(binary, 2, 64)
	Check(err)
	return int(val)
}

func ConvertIntToBinary(val int) string {
	binaryStr := strconv.FormatInt(int64(val), 2)
	return fmt.Sprintf("%064s", binaryStr)
}

func ConvertIntToHex(val int) string {
	return fmt.Sprintf("%x", val)
}

func ConvertBinaryToHex(binary string) string {
	val, err := strconv.ParseInt(binary, 2, 64)
	Check(err)
	return fmt.Sprintf("%x", val)
}

func GetMemoryInfo(tagSize int, indexSize int, kind string, address string) (int, int, int) {
	tagBin := address[:tagSize]
	var index int
	var tag int

	if kind == "full" {
		index = 0
	} else {
		indexBin := address[tagSize : tagSize+indexSize]
		index = ConvertBinaryToInt(indexBin)
	}

	tag = ConvertBinaryToInt(tagBin)
	offsetBin := address[tagSize+indexSize:]
	offset := ConvertBinaryToInt(offsetBin)
	return index, tag, offset
}

func GetIndex(indexSize int, tagSize int, kind string, address string) int {
	if kind == "full" {
		return 0
	}

	indexBin := address[tagSize : tagSize+indexSize]
	return ConvertBinaryToInt(indexBin)
}

func GetTag(tagSize int, address string) int {
	tagBin := address[:tagSize]
	return ConvertBinaryToInt(tagBin)
}

func GetOffset(tagSize int, indexSize int, address string) int {
	offsetBin := address[tagSize+indexSize:]
	return ConvertBinaryToInt(offsetBin)
}
