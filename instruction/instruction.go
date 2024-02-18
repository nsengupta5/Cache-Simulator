package instruction

// This file contains the implementation of the cache simulator. The
// cache simulator reads the trace file and executes the instructions.
// It also handles the cache operations and memory accesses.

import (
	"bufio"
	"math"
	"os"
	"strings"
	"sync"

	"github.com/nsengupta5/Cache-Simulator/cache"
	"github.com/nsengupta5/Cache-Simulator/utils"
)

// bufferSize is the size of the buffer used to read the trace file
const bufferSize int = 8000

type CacheLine = cache.CacheLine

type CacheInstruction struct {
	Addresses []string
}

type CacheSimulator struct {
	Config *cache.CacheConfig
}

// NewCacheSimulator creates a new cache simulator
func NewCacheSimulator(config *cache.CacheConfig) *CacheSimulator {
	return &CacheSimulator{
		Config: config,
	}
}

// Execute serves as the entry point for the cache simulator
// It reads the trace file and executes the instructions and
// prints the cache statistics at the end. One major optimization
// is that it uses goroutines to read the trace file and execute
// the instructions concurrently. Goroutines are lightweight
// threads of execution in Go. Since the trace file can be quite
// large, reading it can be slow. By using goroutines, we can
// read the file concurrently and execute the instructions at the
// same time, which significantly reduces the time taken to execute
// the instructions, reducing the time taken by up to around 50%.
// The buffer size is used to read the trace file concurrently, and
// has been experimentally determined to be the optimal size.
func (cs *CacheSimulator) Execute(traceFile string) {

	// WaitGroups are used to wait for the goroutines to finish
	// This is necessary because we are using goroutines to read
	// the trace file and execute the instructions concurrently,
	// so we must ensure all processing is complete before printing
	// the cache statistics
	var wg sync.WaitGroup

	// The instructions channel is used to send the cache instructions
	// to the goroutine that executes the instructions.
	instructions := make(chan CacheInstruction, bufferSize)

	// The wait group is incremented to wait for the goroutine that
	// reads the trace file to finish and execute the instructions
	wg.Add(1)

	// The goroutine reads the trace file and sends the instructions
	go func() {
		// Defer the Done call to ensure the wait group is decremented,
		// so that the main goroutine can continue
		defer wg.Done()
		file, err := os.Open(traceFile)
		utils.Check(err)
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			instruction := scanner.Text()
			instructionArr := strings.Split(instruction, " ")
			memAddress := utils.ConvertHexToBinary(instructionArr[1])
			size := utils.ConvertStringToInt(instructionArr[3])

			// Use the first cache to calculate the offset and affected addresses
			l1 := cs.Config.Caches[0]
			offset := utils.GetOffset(l1.TagSize, l1.IndexSize, memAddress)
			addresses := getAffectedAddresses(size, l1.LineSize, offset, memAddress)
			instructions <- CacheInstruction{Addresses: addresses}
		}
		// Close the instructions channel to signal that all instructions
		// have been sent
		close(instructions)
	}()

	// The wait group is incremented to wait for the goroutine that
	// executes the instructions to finish
	wg.Add(1)
	go func() {
		defer wg.Done()
		for instruction := range instructions {
			cs.executeInstruction(instruction)
		}
	}()

	// Wait for the goroutines to finish before printing the cache
	// statistics
	wg.Wait()
	cs.Config.PrintStats()
}

// executeInstruction executes the given cache instruction
// It loops over all the addresses in the instruction and
// calls handleCacheOperations to check if the data is present
// in the cache and updates the cache statistics accordingly
func (cs *CacheSimulator) executeInstruction(instruction CacheInstruction) {
	for i := 0; i < len(instruction.Addresses); i++ {
		if !cs.handleCacheOperations(instruction.Addresses[i]) {
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
