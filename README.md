# Cache-Simulator
A multi-hierarchy cache simulator that supports LRU, LFU and Round Robin eviction policies

## Instructions

To run using the executable:
```bash
./cache_simulator ./sample-inputs/<input-file> /cs/studres/CS4202/Coursework/P1-CacheSim/trace-files/<trace-file>
```

To compile and run:
```bash
go run main.go ./sample-inputs/<input-file> /cs/studres/CS4202/Coursework/P1-CacheSim/trace-files/<trace-file>
```

To build the executable:
```bash
go build -o cache_simulator main.go
```
