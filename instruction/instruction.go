package instruction

// This file contains the implementation of the cache simulator. The
// cache simulator reads the trace file and executes the instructions.
// It also handles the cache operations and memory accesses.

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

// NewCacheSimulator creates a new cache simulator
func NewCacheSimulator(config *cache.CacheConfig) *CacheSimulator {
	return &CacheSimulator{
		Config: config,
	}
}

// ReadTraceFile reads the trace file and executes the instructions
// After all instructions are executed, it prints the cache statistics
func (cs *CacheSimulator) ReadTraceFile(traceFile string) {
	file, err := os.Open(traceFile)
	utils.Check(err)
	defer file.Close()

	// Use a buffer to read the file line by line, to avoid reading
	// the entire file into memory
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		instruction := scanner.Text()
		cs.executeInstruction(instruction)
	}

	cs.Config.PrintStats()
	err = scanner.Err()
	utils.Check(err)
}

// ExecuteInstruction first obtains the necessary information to handle
// the cache operation. It then checks if the data is present in the cache
// and if not, fetches it from memory. It also updates the cache statistics,
// partcularly, the number of memory accesses, hits and misses.
func (cs *CacheSimulator) executeInstruction(instruction string) {
	instructionArr := strings.Split(instruction, " ")
	memAddress := utils.ConvertHexToBinary(instructionArr[1])
	size := utils.ConvertStringToInt(instructionArr[3])

	// Use the first cache to calculate the offset and affected addresses
	l1 := cs.Config.Caches[0]
	offset := utils.GetOffset(l1.TagSize, l1.IndexSize, memAddress)
	addresses := getAffectedAddresses(size, l1.LineSize, offset, memAddress)

	// Here we are looping over all affected addresses - this is where if the
	// size of the operation is larger than the line size, we will have to
	// handle multiple cache operations. The length of the addresses array
	// indicates the number of cache operations we need to handle.
	for i := 0; i < len(addresses); i++ {
		// If not hits in any cache, then it's a memory access
		if !cs.handleCacheOperations(addresses[i]) {
			cs.Config.MemoryAccesses++
		}
	}
}

// handleCacheOperations checks if the data is present in the cache
// If not, it fetches it from memory and updates the cache statistics
// It returns a boolean indicating if the data was found in the cache
func (cs *CacheSimulator) handleCacheOperations(address string) bool {
	var dataFound bool = false
	var tag int
	var index int

	// For each address, we loop over all caches to check if the data
	// is present
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

		// If the data is found in the cache, we update the cache statistics
		// If the cache has a policy, we also update the policy statistics
		// i.e the frequency and age of the line for LFU and LRU policies
		// respectively. We break out of the loop as we don't need to check
		// the other caches if a hit is found.
		if hit {
			cs.Config.Caches[j].Hits++
			set.Policy.Update(line)
			dataFound = true
			break
		} else {
			// If the data is not found in the cache, we update the cache
			// miss statistics and assign a new cache line to the data.
			// Depending on the cache kind, we either insert the data directly
			// or use the cache policy to insert the data.
			cs.Config.Caches[j].Misses++

			// A new cache line will have 1 frequency and 0 age,
			// where Freq represents the number of times the line
			// has been accessed and Age represents the number of
			// instructions since the line was last accessed.
			data := &CacheLine{
				Tag:   tag,
				Valid: true,
				Index: -1,
				Freq:  1,
				Prev:  nil,
				Next:  nil,
			}
			if cache.Kind == "direct" {
				set.Lines[0] = *data
			} else {
				set.Insert(data)
			}
		}
	}
	return dataFound
}

// getAffectedAddresses returns the addresses affected by the operation
// If the size of the operation is larger than the line size, we will
// have to handle multiple cache operations. The length of the addresses
// array indicates the number of cache operations we need to handle.
func getAffectedAddresses(size int, lineSize int, offset int, memAddress string) []string {
	memAddressInt := utils.ConvertBinaryToInt(memAddress)
	addresses := []string{memAddress}

	// We first calculate the initial address. If the size of the operation
	// is larger than the line size, we will have to handle multiple cache
	// operations. Otherwise, we only need to handle one cache operation.
	initialBytes := lineSize - offset
	if size <= initialBytes {
		return addresses
	}

	remainingBytes := size - initialBytes
	// We calculate the number of remaining addresses we need to handle
	// based on the remaining bytes and the line size
	// The ceil function is used to round up the number of addresses
	// as we can't have a fraction of an address
	remainingAddresses := int(math.Ceil(float64(remainingBytes) / float64(lineSize)))

	// We then calculate the remaining addresses and append them to the
	// addresses array
	for i := 1; i <= remainingAddresses; i++ {
		address := memAddressInt + (i * lineSize)
		addresses = append(addresses, utils.ConvertIntToBinary(address))
	}

	return addresses
}
