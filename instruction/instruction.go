package instruction

import (
	"bufio"
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
		// size := utils.ConvertStringToUint(instructionArr[3])

		var dataLine *CacheLine
		var tag uint
		var index uint
		var breakIdx int
		var dataFound bool

		for i := 0; i < len(config.Caches); i++ {
			cache := config.Caches[i]
			index, tag, _ = utils.GetMemoryInfo(
				cache.TagSize,
				cache.IndexSize,
				cache.Kind,
				memAddress,
			)
			hit, line, _ := cache.CheckHitOrMiss(tag, index)
			if hit {
				config.Caches[i].Hits++
				dataLine = line
				dataFound = true
				breakIdx = i
				break
			} else {
				config.Caches[i].Misses++
				breakIdx = i
			}

		}

		if !dataFound {
			config.MemoryAccesses++
			dataLine = FetchFromMemory(tag)
		}
		UpdateCaches(config, dataLine, memAddress, breakIdx)
	}

	config.PrintStats()
	err = scanner.Err()
	utils.Check(err)
}

func FetchFromMemory(tag uint) *CacheLine {
	return &CacheLine{
		Tag:      tag,
		Valid:    true,
		Index:    -1,
		LFUIndex: 0,
		Freq:     0,
	}
}

func UpdateCaches(config *cache.CacheConfig, data *CacheLine, memAddress string, breakIdx int) {
	for i := 0; i <= breakIdx; i++ {
		cache := config.Caches[i]
		index, _, _ := utils.GetMemoryInfo(
			cache.TagSize,
			cache.IndexSize,
			cache.Kind,
			memAddress,
		)

		set := cache.Sets[index]
		if cache.Kind == "direct" {
			set.Lines[0] = *data
		} else {
			set.Insert(data)
		}
	}
}
