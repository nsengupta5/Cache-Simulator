package main

// This file contains the main function to run the cache simulator.
// It reads in the command line arguements and initializes the cache configuration
// and the cache simulator.

import (
	"os"

	"github.com/nsengupta5/Cache-Simulator/cache"
	"github.com/nsengupta5/Cache-Simulator/instruction"
)

func main() {
	// Read in the command line arguements
	configFile := os.Args[1]
	traceFile := os.Args[2]

	config := cache.InitializeConfig(configFile)
	cache.InitializeCaches(&config)
	simulator := instruction.NewCacheSimulator(&config)
	simulator.ReadTraceFile(traceFile)
}
