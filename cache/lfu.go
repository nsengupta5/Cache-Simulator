package cache

// This file contains the implementation of the Least Frequently Used (LFU)
// cache replacement policy. The LFU policy evicts the least frequently used
// cache line when the cache is full and a new cache line needs to be inserted.
// To implement the LFU policy, we keep track of the frequency of each cache
// line, and evict the cache line with the lowest frequency when necessary.
// The LFU policy is implemented using the LFU struct, which contains an array
// of CacheLine structs, and a capacity field to keep track of the maximum
// number of cache lines that can be stored in the cache.

type LFU struct {
	lines    []CacheLine
	capacity int
}

func NewLFU(capacity int) *LFU {
	return &LFU{
		lines:    make([]CacheLine, capacity), // Holds the cache lines
		capacity: capacity,                    // Maximum number of cache lines
	}
}

// Insert a cache line into the cache
func (lfu *LFU) Insert(line *CacheLine) {
	// Keep track of the index of the cache line with the
	// lowest frequency, and the lowest frequency itself
	minFreqIndex := -1
	minFreq := int(^uint(0) >> 1)
	for i, l := range lfu.lines {
		// If the line is not valid, update the index and break
		if !l.Valid {
			minFreqIndex = i
			break
			// If the line is valid and its frequency is less than the
			// current minimum frequency, update the minimum frequency
		} else if l.Freq < minFreq {
			minFreq = l.Freq
			minFreqIndex = i
		}
	}

	// Insert the line at the index of the cache line
	// with the lowest frequency
	if minFreqIndex != -1 {
		lfu.lines[minFreqIndex] = *line
	}
}

// Update the frequency of a cache line when it is accessed
func (lfu *LFU) Update(line *CacheLine) {
	for i, l := range lfu.lines {
		if l.Tag == line.Tag && l.Valid {
			lfu.lines[i].Freq++
			break
		}
	}
}

// Evict identifies the cache line with the lowest frequency
// It returns the index of the cache line to be evicted
func (lfu *LFU) Evict() int {
	minFreqIndex := -1
	minFreq := int(^uint(0) >> 1)
	for i, l := range lfu.lines {
		if l.Valid && l.Freq < minFreq {
			minFreq = l.Freq
			minFreqIndex = i
		}
	}
	if minFreqIndex != -1 {
		evictedIndex := lfu.lines[minFreqIndex].Index
		lfu.lines[minFreqIndex].Valid = false // Mark as evicted
		return evictedIndex
	}
	return -1 // Indicate no eviction occurred (should not happen if cache is full)
}
