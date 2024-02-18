package cache

type LFU struct {
	lines    []CacheLine
	capacity int
}

func NewLFU(capacity int) *LFU {
	return &LFU{
		lines:    make([]CacheLine, capacity),
		capacity: capacity,
	}
}

// Insert a new cache line
func (lfu *LFU) Insert(line *CacheLine) {
	// Find an empty spot or the least frequently used cache line
	minFreqIndex := -1
	minFreq := int(^uint(0) >> 1) // Initialize to max int value
	for i, l := range lfu.lines {
		if !l.Valid { // Empty spot found
			minFreqIndex = i
			break
		} else if l.Freq < minFreq {
			minFreq = l.Freq
			minFreqIndex = i
		}
	}
	if minFreqIndex != -1 {
		lfu.lines[minFreqIndex] = *line
	}
}

// Update the frequency of a cache line
func (lfu *LFU) Update(line *CacheLine) {
	for i, l := range lfu.lines {
		if l.Tag == line.Tag && l.Valid {
			lfu.lines[i].Freq++
			break
		}
	}
}

// Evict the least frequently used cache line
func (lfu *LFU) Evict() int {
	minFreqIndex := -1
	minFreq := int(^uint(0) >> 1) // Initialize to max int value
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

// Get the capacity of the LFU cache
func (lfu *LFU) GetCapacity() int {
	return lfu.capacity
}
