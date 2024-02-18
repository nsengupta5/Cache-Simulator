package cache

// This file contains the implementation of the cache and its components
// The cache design is dynamically intialized based on the cache's kind
// which allows for a flexible and scalable approach to simulating various
// cache configurations by adjusting the number of sets and lines within
// each cache according to its architecture. As such, it makes it especially
// easy to extend to higher associativity caches and different replacement
// policies.

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/nsengupta5/Cache-Simulator/utils"
)

// The addresss size for ths practical is 64 bits
const addressSize int = 64

// The ReplacementPolicy interface is a contract for implementing
// different replacement policies for the cache
type ReplacementPolicy interface {
	Insert(line *CacheLine)
	Update(line *CacheLine)
	Evict() int
}

// The CacheLine struct represents a line in the cache
type CacheLine struct {
	Valid bool `json:"valid"`
	Tag   int  `json:"tag"`
	Freq  int  `json:"frequency"`
	Index int  `json:"index"`
	Age   int  `json:"last_access"`
}

// The CacheSet struct represents a set in the cache
type CacheSet struct {
	Lines  []CacheLine       `json:"lines"`
	Size   int               `json:"set_size"`
	Policy ReplacementPolicy `json:"replacement_policy"`
}

// The Cache struct represents a cache
type Cache struct {
	Sets       []CacheSet `json:"sets"`
	Name       string     `json:"name"`
	Size       int        `json:"size"`
	PolicyName string     `json:"replacement_policy"`
	Kind       string     `json:"kind"`
	LineSize   int        `json:"line_size"`
	TagSize    int        `json:"tag_size"`
	IndexSize  int        `json:"index_size"`
	OffsetSize int        `json:"offset_size"`
	Hits       int        `json:"hits"`
	Misses     int        `json:"misses"`
}

// The CacheConfig struct represents the configuration of
// the cache
type CacheConfig struct {
	Caches         []Cache `json:"caches"`
	MemoryAccesses int     `json:"memory_accesses"`
}

/* ------------------- Cache Function ------------------- */

// SetLiesSize sets the number of lines in each set
// It checks the kind of cache and sets the number of lines accordingly
func (cache *Cache) SetLinesSize() {
	cacheLines := cache.Size / cache.LineSize
	for i := range cache.Sets {
		set := &cache.Sets[i]
		switch cache.Kind {
		case "direct":
			set.Lines = make([]CacheLine, 1)
		case "full":
			set.Lines = make([]CacheLine, cacheLines)
		case "2way":
			set.Lines = make([]CacheLine, 2)
		case "4way":
			set.Lines = make([]CacheLine, 4)
		default:
			set.Lines = make([]CacheLine, 8)
		}
	}
}

// SetSetsSize sets the number of sets in the cache
// It checks the kind of cache and sets the number of sets accordingly
func (cache *Cache) SetSetsSize() {
	cacheLines := cache.Size / cache.LineSize
	var setSize int
	switch cache.Kind {
	case "direct":
		setSize = cacheLines
		cache.Sets = make([]CacheSet, setSize)
	case "full":
		cache.Sets = make([]CacheSet, 1)
	case "2way":
		setSize = cacheLines / 2
		cache.Sets = make([]CacheSet, setSize)
	case "4way":
		setSize = cacheLines / 4
		cache.Sets = make([]CacheSet, setSize)
	default:
		setSize = cacheLines / 8
		cache.Sets = make([]CacheSet, setSize)
	}
}

// SetDefaultPolicy sets the default replacement policy for the cache
func (cache *Cache) SetDefaultPolicy() {
	if cache.Kind != "direct" && cache.PolicyName == "" {
		// Default policy for set associative caches is round robin
		cache.PolicyName = "rr"
	}

	for s := range cache.Sets {
		set := &cache.Sets[s]
		capacity := len(set.Lines)
		switch cache.PolicyName {
		case "lru":
			set.Policy = NewLRU(capacity)
		case "lfu":
			set.Policy = NewLFU(capacity)
		default:
			set.Policy = NewRR(capacity)
		}
	}
}

// SetBitsSize sets the number of bits for offset, index and tag
func (cache *Cache) SetBitsSize() {
	cache.OffsetSize = cache.getOffsetBits()
	cache.IndexSize = cache.getIndexBits()
	cache.TagSize = cache.getTagBits()
}

// getOffsetBits returns the number of bits for the offset
func (cache *Cache) getOffsetBits() int {
	lineSize := cache.LineSize
	return int(math.Log2(float64(lineSize)))
}

// getIndexBits returns the number of bits for the index
// The number of bits is dependent on the kind of cache
func (cache *Cache) getIndexBits() int {
	var setSize int
	switch cache.Kind {
	case "direct":
		setSize = cache.Size / cache.LineSize
		return int(math.Log2(float64(setSize)))
	case "full":
		return 0
	case "2way", "4way", "8way":
		setSize = int(len(cache.Sets))
		return int(math.Log2(float64(setSize)))
	default:
		return 0
	}
}

// getTagBits returns the number of bits for the tag
// The number of bits is dependent on the kind of cache
func (cache *Cache) getTagBits() int {
	var tagBits int
	switch cache.Kind {
	// The tag is the remaining bits before the offset and index bits
	case "full":
		tagBits = addressSize - cache.OffsetSize
	default:
		tagBits = addressSize - cache.OffsetSize - cache.IndexSize
	}
	return tagBits
}

// CheckHitOrMiss checks if the address is in the cache
// It identifies the set from the index and looks for the tag in the set

// It returns a boolean indicating if the address is in the cache and the cache line
func (cache *Cache) CheckHitOrMiss(tag int, index int) (bool, *CacheLine) {
	set := cache.Sets[index]

	for i := range set.Lines {
		line := &set.Lines[i]
		// If the tag is the same and the line is valid, it's a hit
		if line.Tag == tag && line.Valid {
			return true, line
		}
	}

	return false, nil
}

// GetStats returns the cache statistics
func (cache *Cache) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"hits":   cache.Hits,
		"misses": cache.Misses,
		"name":   cache.Name,
	}
}

/* ------------------- Cache Set Function ------------------- */

// Insert adds a new line to the set
// This function is exclusive to the set associative caches
func (set *CacheSet) Insert(newLine *CacheLine) {
	for i := range set.Lines {
		line := &set.Lines[i]
		if !line.Valid {
			// Set the index of the new line to the index of the invalid line
			newLine.Index = i

			// Update the policy with the new line
			set.Policy.Insert(newLine)
			set.Lines[i] = *newLine
			return
		}
	}

	// If the set is full, evict a line and insert the new line
	evictIndex := set.Policy.Evict()
	newLine.Index = evictIndex
	set.Policy.Insert(newLine)
	set.Lines[evictIndex] = *newLine
}

/* ------------------- Cache Config Function ------------------- */

// PrintStats prints the cache statistics
func (config *CacheConfig) PrintStats() {
	cacheStats := []map[string]interface{}{}

	for _, cache := range config.Caches {
		stats := cache.GetStats()
		cacheStats = append(cacheStats, stats)
	}

	stats := map[string]interface{}{
		"caches":               cacheStats,
		"main_memory_accesses": config.MemoryAccesses,
	}

	output, err := json.MarshalIndent(stats, "", "  ")
	utils.Check(err)
	fmt.Println(string(output))
}
