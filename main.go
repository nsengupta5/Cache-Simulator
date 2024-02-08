package main

import (
	"os"

	"github.com/nsengupta5/Cache-Simulator/cache"
	"github.com/nsengupta5/Cache-Simulator/instruction"
)

func main() {
	configFile := os.Args[1]
	traceFile := os.Args[2]

	config := cache.InitializeConfig(configFile)
	cache.InitializeCaches(&config)
	instruction.ReadTraceFile(&config, traceFile)
}
