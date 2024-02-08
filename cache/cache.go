package cache

import (
	"encoding/json"
	"math"
	"os"

	"github.com/nsengupta5/Cache-Simulator/utils"
)

const addressSize = 64

type CacheLine struct {
	Data  []byte `json:"data"`
	Valid bool   `json:"valid"`
	Tag   int    `json:"tag"`
	Size  int    `json:"line_size"`
	Index int    `json:"index"`
}

type CacheSet struct {
	Lines []CacheLine
	Size  int `json:"set_size"`
}

type Cache struct {
	Sets       []CacheSet `json:"sets"`
	Name       string     `json:"name"`
	Size       int        `json:"size"`
	Policy     string     `json:"replacement_policy"`
	Kind       string     `json:"kind"`
	LineSize   int        `json:"line_size"`
	TagSize    int        `json:"tag_size"`
	IndexSize  int        `json:"index_size"`
	OffsetSize int        `json:"offset_size"`
}

type CacheConfig struct {
	Caches []Cache `json:"caches"`
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

func (cache *Cache) SetDefaultPolicy() {
	if cache.Kind != "direct" && cache.Policy == "" {
		cache.Policy = "rr"
	}
}

func (cache *Cache) SetBitsSize() {
	cache.OffsetSize = cache.getOffsetBits()
	cache.IndexSize = cache.getIndexBits()
	cache.TagSize = cache.getTagBits()
}

func (cache *Cache) getOffsetBits() int {
	lineSize := cache.LineSize
	return int(math.Log2(float64(lineSize)))
}

func (cache *Cache) getIndexBits() int {
	var setSize int
	switch cache.Kind {
	case "direct":
		setSize = cache.Size / cache.LineSize
		return int(math.Log2(float64(setSize)))
	case "full":
		return 0
	case "2way", "4way", "6way", "8way":
		setSize = len(cache.Sets)
		return int(math.Log2(float64(setSize)))
	default:
		return 0
	}
}

func (cache *Cache) getTagBits() int {
	var tagBits int
	switch cache.Kind {
	case "full":
		tagBits = addressSize - cache.OffsetSize
	default:
		tagBits = addressSize - cache.OffsetSize - cache.IndexSize
	}
	return tagBits
}

func InitializeCaches(config *CacheConfig) {
	for i := range config.Caches {
		cache := &config.Caches[i]
		cache.SetSetsSize()
		cache.SetLinesSize()
		cache.SetBitsSize()
		cache.SetDefaultPolicy()
	}
}

func InitializeConfig(configFile string) CacheConfig {
	cacheData, err := os.ReadFile(configFile)
	utils.Check(err)

	var config CacheConfig
	err = json.Unmarshal(cacheData, &config)
	utils.Check(err)

	return config
}
