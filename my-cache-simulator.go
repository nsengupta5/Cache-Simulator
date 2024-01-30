package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type CacheLine struct {
	Data  []byte `json:"data"`
	Valid bool   `json:"valid"`
	Tag   int    `json:"tag"`
	Size  int    `json:"line_size"`
}

type CacheSet struct {
	Lines []CacheLine
	Size  int `json:"set_size"`
}

type Cache struct {
	Sets     []CacheSet `json:"sets"`
	Name     string     `json:"name"`
	Size     int        `json:"size"`
	Policy   string     `json:"replacement_policy"`
	Kind     string     `json:"kind"`
	LineSize int        `json:"line_size"`
}

type CacheConfig struct {
	Caches []Cache `json:"caches"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func setDefaultPolicy(config *CacheConfig) {
	for _, cache := range config.Caches {
		if cache.Kind != "direct" && cache.Policy == "" {
			cache.Policy = "rr"
		}
	}
}

func setSetsAndLinesSize(config *CacheConfig) {
	for i := range config.Caches {
		cache := &config.Caches[i]
		cache.setSetsSize()
		cache.setLinesSize()
	}
}

func (cache *Cache) setLinesSize() {
	cacheLines := cache.Size / cache.LineSize
	linesPerSet := cacheLines / len(cache.Sets)
	for i := range cache.Sets {
		cache.Sets[i].Lines = make([]CacheLine, linesPerSet)
	}
}

func (cache *Cache) setSetsSize() {
	switch cache.Kind {
	case "direct", "full":
		cache.Sets = make([]CacheSet, 1)
	case "2way":
		cache.Sets = make([]CacheSet, 2)
	case "4way":
		cache.Sets = make([]CacheSet, 4)
	default:
		cache.Sets = make([]CacheSet, 8)
	}
}

func main() {
	configFile := os.Args[1]
	// traceFile := os.Args[2]

	cacheData, err := os.ReadFile(configFile)
	check(err)
	fmt.Println(string(cacheData))

	var config CacheConfig
	err = json.Unmarshal(cacheData, &config)
	check(err)

	setDefaultPolicy(&config)
	setSetsAndLinesSize(&config)

	for _, cache := range config.Caches {
		fmt.Println("Cache Set Size: ", len(cache.Sets))
		fmt.Println("Cache Lines Size: ", len(cache.Sets[0].Lines))
		fmt.Println("")
	}
}
