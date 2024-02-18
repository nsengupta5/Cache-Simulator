package instruction

import (
	"bufio"
	"math"
	"os"
	"strings"

	"github.com/nsengupta5/Cache-Simulator/cache"
	"github.com/nsengupta5/Cache-Simulator/utils"
)

type CacheLine = cache.CacheLine

func ReadTraceFile(config *cache.CacheConfig, traceFile string) {
	file, err := os.Open(traceFile)
	utils.Check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		instructionString := scanner.Text()
		instructionArr := strings.Split(instructionString, " ")
		memAddress := utils.ConvertHexToBinary(instructionArr[1])
		size := utils.ConvertStringToInt(instructionArr[3])

		l1 := config.Caches[0]
		offset := utils.GetOffset(l1.TagSize, l1.IndexSize, memAddress)
		addresses := getAffectedAddresses(size, l1.LineSize, offset, memAddress)

		for i := 0; i < len(addresses); i++ {
			var dataFound bool = false
			var tag int
			var index int

			for j := 0; j < len(config.Caches); j++ {
				cache := config.Caches[j]
				index, tag, _ = utils.GetMemoryInfo(
					cache.TagSize,
					cache.IndexSize,
					cache.Kind,
					addresses[i],
				)

				hit, line, _ := cache.CheckHitOrMiss(tag, index)
				set := cache.Sets[index]
				if hit {
					config.Caches[j].Hits++
					set.Policy.Update(line)
					dataFound = true
					break
				} else {
					config.Caches[j].Misses++
					data := FetchFromMemory(tag)
					if cache.Kind == "direct" {
						set.Lines[0] = *data
					} else {
						set.Insert(data)
					}
				}
			}
			if !dataFound {
				config.MemoryAccesses++
			}
		}
	}

	config.PrintStats()
	err = scanner.Err()
	utils.Check(err)
}

func FetchFromMemory(tag int) *CacheLine {
	return &CacheLine{
		Tag:   tag,
		Valid: true,
		Index: -1,
		Freq:  1,
		Age:   0,
	}
}

func getAffectedAddresses(size int, lineSize int, offset int, memAddress string) []string {
	memAddressInt := utils.ConvertBinaryToInt(memAddress)
	addresses := []string{memAddress}

	initialBytes := lineSize - offset
	if size <= initialBytes {
		return addresses
	}

	remainingBytes := size - initialBytes
	remainingAddresses := int(math.Ceil(float64(remainingBytes) / float64(lineSize)))

	for i := 1; i <= remainingAddresses; i++ {
		address := memAddressInt + (i * lineSize)
		addresses = append(addresses, utils.ConvertIntToBinary(address))
	}

	return addresses
}
