package cache

// This file contains the functions to initialize the cache configuration
// and caches.

import (
	"encoding/json"
	"os"

	"github.com/nsengupta5/Cache-Simulator/utils"
)

// InitializeCaches initializes the caches with the given configuration.
// It sets the size of the sets, lines, bits and the default policy of
// the caches.
func InitializeCaches(config *CacheConfig) {
	for i := range config.Caches {
		cache := &config.Caches[i]
		cache.SetSetsSize()
		cache.SetLinesSize()
		cache.SetBitsSize()
		cache.SetDefaultPolicy()
	}
}

// InitializeConfig initializes the cache configuration with the given
// configuration file. It reads the JSON config file and unmarshals
// the data into the CacheConfig struct.
func InitializeConfig(configFile string) CacheConfig {
	cacheData, err := os.ReadFile(configFile)
	utils.Check(err)

	var config CacheConfig
	err = json.Unmarshal(cacheData, &config)
	utils.Check(err)

	return config
}
