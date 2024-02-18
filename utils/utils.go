package utils

// This file contains utility functions that are used in the cache simulator
// These function are used for various purposes, such as converting between
// different representations of numbers (binary, hex, int), and extracting
// information from memory addresses

import (
	"fmt"
	"strconv"
)

// Check is a helper function to handle errors
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

// ConvertBinaryToHex converts a hex string to a binary string
func ConvertHexToBinary(hex string) string {
	val, err := strconv.ParseInt(hex, 16, 64)
	Check(err)
	binaryStr := strconv.FormatInt(val, 2)
	return fmt.Sprintf("%064s", binaryStr)
}

// ConvertHexToInt converts a hex string to an int
func ConvertHexToInt(hex string) int {
	val, err := strconv.ParseInt(hex, 16, 64)
	Check(err)
	return int(val)
}

// ConvertStringToRune converts a string to a rune
func ConvertStringToRune(str string) rune {
	return []rune(str)[0]
}

// ConvertStringToInt converts a string to an int
func ConvertStringToInt(str string) int {
	val, err := strconv.ParseInt(str, 10, 64)
	Check(err)
	return int(val)
}

// ConvertBinaryToInt converts a binary string to an int
func ConvertBinaryToInt(binary string) int {
	val, err := strconv.ParseInt(binary, 2, 64)
	Check(err)
	return int(val)
}

// ConvertIntToBinary converts an int to a binary string
func ConvertIntToBinary(val int) string {
	binaryStr := strconv.FormatInt(int64(val), 2)
	return fmt.Sprintf("%064s", binaryStr)
}

// ConvertIntToHex converts an int to a hex string
func ConvertIntToHex(val int) string {
	return fmt.Sprintf("%x", val)
}

// ConvertHexToBinary converts a hex string to a binary string
func ConvertBinaryToHex(binary string) string {
	val, err := strconv.ParseInt(binary, 2, 64)
	Check(err)
	return fmt.Sprintf("%x", val)
}

// GetMemoryInfo extracts the index, tag, and offset from a memory address
func GetMemoryInfo(tagSize int, indexSize int, kind string, address string) (int, int, int) {
	tagBin := address[:tagSize]
	var index int
	var tag int

	if kind == "full" {
		// If the cache is fully associative, the index is always 0
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

// GetIndex extracts the index from a memory address
func GetIndex(indexSize int, tagSize int, kind string, address string) int {
	if kind == "full" {
		return 0
	}

	indexBin := address[tagSize : tagSize+indexSize]
	return ConvertBinaryToInt(indexBin)
}

// GetTag extracts the tag from a memory address
func GetTag(tagSize int, address string) int {
	tagBin := address[:tagSize]
	return ConvertBinaryToInt(tagBin)
}

// GetOffset extracts the offset from a memory address
func GetOffset(tagSize int, indexSize int, address string) int {
	offsetBin := address[tagSize+indexSize:]
	return ConvertBinaryToInt(offsetBin)
}
