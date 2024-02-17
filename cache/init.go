package cache

import (
	"encoding/json"
	"os"

	"github.com/nsengupta5/Cache-Simulator/utils"
)

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
