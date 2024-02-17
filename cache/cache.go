package cache

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/nsengupta5/Cache-Simulator/utils"
)

const addressSize uint = 64

type CacheLine struct {
	Valid    bool `json:"valid"`
	Tag      uint `json:"tag"`
	Freq     uint `json:"frequency"`
	LFUIndex int  `json:"lru_index"`
	Index    int  `json:"index"`
}

type CacheSet struct {
	Lines  []CacheLine       `json:"lines"`
	Size   uint              `json:"set_size"`
	Policy ReplacementPolicy `json:"replacement_policy"`
}

type Cache struct {
	Sets       []CacheSet `json:"sets"`
	Name       string     `json:"name"`
	Size       uint       `json:"size"`
	PolicyName string     `json:"replacement_policy"`
	Kind       string     `json:"kind"`
	LineSize   uint       `json:"line_size"`
	TagSize    uint       `json:"tag_size"`
	IndexSize  uint       `json:"index_size"`
	OffsetSize uint       `json:"offset_size"`
	Hits       uint       `json:"hits"`
	Misses     uint       `json:"misses"`
}

type CacheConfig struct {
	Caches         []Cache `json:"caches"`
	MemoryAccesses int     `json:"memory_accesses"`
}

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

func (cache *Cache) SetSetsSize() {
	cacheLines := cache.Size / cache.LineSize
	var setSize uint
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

func (cache *Cache) SetDefaultPolicy() {
	if cache.Kind != "direct" && cache.PolicyName == "" {
		cache.PolicyName = "rr"
	}

	//
	for s := range cache.Sets {
		set := &cache.Sets[s]
		capacity := len(set.Lines)
		switch cache.PolicyName {
		// case "lru":
		// 	set.Policy = NewLRU()
		case "lfu":
			set.Policy = NewLFU(capacity)
		default:
			set.Policy = NewRR(capacity)
		}
	}
}

func (cache *Cache) SetBitsSize() {
	cache.OffsetSize = cache.getOffsetBits()
	cache.IndexSize = cache.getIndexBits()
	cache.TagSize = cache.getTagBits()
}

func (cache *Cache) getOffsetBits() uint {
	lineSize := cache.LineSize
	return uint(math.Log2(float64(lineSize)))
}

func (cache *Cache) getIndexBits() uint {
	var setSize uint
	switch cache.Kind {
	case "direct":
		setSize = cache.Size / cache.LineSize
		return uint(math.Log2(float64(setSize)))
	case "full":
		return 0
	case "2way", "4way", "8way":
		setSize = uint(len(cache.Sets))
		return uint(math.Log2(float64(setSize)))
	default:
		return 0
	}
}

func (cache *Cache) getTagBits() uint {
	var tagBits uint
	switch cache.Kind {
	case "full":
		tagBits = addressSize - cache.OffsetSize
	default:
		tagBits = addressSize - cache.OffsetSize - cache.IndexSize
	}
	return tagBits
}

func (cache *Cache) CheckHitOrMiss(tag uint, index uint) (bool, *CacheLine, int) {
	set := cache.Sets[index]

	for i := range set.Lines {
		line := &set.Lines[i]
		if line.Tag == tag && line.Valid {
			return true, line, i
		}
	}

	return false, nil, -1
}

func (set *CacheSet) Insert(newLine *CacheLine) {
	if len(set.Lines) == set.Policy.GetCapacity() {
		evictedLineIdx := set.Policy.Evict()
		set.Lines[evictedLineIdx] = *newLine
	} else {
		set.Lines = append(set.Lines, *newLine)
		set.Policy.Insert(newLine)
		set.Policy.Update(newLine)
	}
}

func (cache *Cache) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"hits":   cache.Hits,
		"misses": cache.Misses,
		"name":   cache.Name,
	}
}

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
