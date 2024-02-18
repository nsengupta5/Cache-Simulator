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

type CacheSimulator struct {
	Config *cache.CacheConfig
}

func NewCacheSimulator(config *cache.CacheConfig) *CacheSimulator {
	return &CacheSimulator{
		Config: config,
	}
}

func (cs *CacheSimulator) ReadTraceFile(traceFile string) {
	file, err := os.Open(traceFile)
	utils.Check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		instruction := scanner.Text()
		cs.ExecuteInstruction(instruction)
	}

	cs.Config.PrintStats()
	err = scanner.Err()
	utils.Check(err)
}

func (cs *CacheSimulator) ExecuteInstruction(instruction string) {
	instructionArr := strings.Split(instruction, " ")
	memAddress := utils.ConvertHexToBinary(instructionArr[1])
	size := utils.ConvertStringToInt(instructionArr[3])

	l1 := cs.Config.Caches[0]
	offset := utils.GetOffset(l1.TagSize, l1.IndexSize, memAddress)
	addresses := getAffectedAddresses(size, l1.LineSize, offset, memAddress)

	for i := 0; i < len(addresses); i++ {
		if !cs.handleCacheOperations(addresses[i]) {
			cs.Config.MemoryAccesses++
		}
	}
}

func (cs *CacheSimulator) handleCacheOperations(address string) bool {
	var dataFound bool = false
	var tag int
	var index int

	for j := 0; j < len(cs.Config.Caches); j++ {
		cache := cs.Config.Caches[j]
		index, tag, _ = utils.GetMemoryInfo(
			cache.TagSize,
			cache.IndexSize,
			cache.Kind,
			address,
		)

		hit, line := cache.CheckHitOrMiss(tag, index)
		set := cache.Sets[index]
		if hit {
			cs.Config.Caches[j].Hits++
			set.Policy.Update(line)
			dataFound = true
			break
		} else {
			cs.Config.Caches[j].Misses++
			data := FetchFromMemory(tag)
			if cache.Kind == "direct" {
				set.Lines[0] = *data
			} else {
				set.Insert(data)
			}
		}
	}
	return dataFound
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
