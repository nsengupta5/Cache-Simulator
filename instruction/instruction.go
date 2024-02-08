package instruction

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nsengupta5/Cache-Simulator/cache"
	"github.com/nsengupta5/Cache-Simulator/utils"
)

type Instruction struct {
	PCAddress     string `json:"data"`
	MemoryAddress string `json:"memory_address"`
	Kind          rune   `json:"kind"`
	Size          uint64 `json:"size"`
}

func ReadTraceFile(config *cache.CacheConfig, traceFile string) {
	file, err := os.Open(traceFile)
	utils.Check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		instructionString := scanner.Text()
		instructionArr := strings.Split(instructionString, " ")
		instruction := Instruction{
			PCAddress:     utils.ConvertHexToBinary(instructionArr[0]),
			MemoryAddress: utils.ConvertHexToBinary(instructionArr[1]),
			Kind:          utils.ConvertStringToRune(instructionArr[2]),
			Size:          utils.ConvertStringToUint(instructionArr[3]),
		}

		for _, cache := range config.Caches {
			fmt.Println("Cache: ", cache.Name)
			extractInstructionBits(cache, instruction.MemoryAddress)
		}
	}
	err = scanner.Err()
	utils.Check(err)
}

func extractInstructionBits(cache cache.Cache, instruction string) {
	indexEnd := cache.TagSize + cache.IndexSize
	tag := instruction[:cache.TagSize]
	index := instruction[cache.TagSize:indexEnd]
	offset := instruction[indexEnd:]

	fmt.Println("Tag: ", tag)
	fmt.Println("Index: ", index)
	fmt.Println("Offset: ", offset)

	decIndex, err := strconv.ParseInt(index, 2, 64)
	utils.Check(err)
	fmt.Println("Decimal Index: ", decIndex)
	setIndex := int(decIndex) % len(cache.Sets)
	cacheLine := cache.Sets[setIndex].Lines
	fmt.Println("Cache Line: ", cacheLine)
}
