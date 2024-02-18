package cache

type LRU struct {
	lines    []CacheLine
	capacity int
}

func NewLRU(capacity int) *LRU {
	return &LRU{
		lines:    make([]CacheLine, capacity),
		capacity: capacity,
	}
}

// Insert a new cache line or update an existing one, resetting its age
func (lru *LRU) Insert(line *CacheLine) {
	// Increment age of all valid lines
	for i := range lru.lines {
		if lru.lines[i].Valid {
			lru.lines[i].Age++
		}
	}

	// Find a place to insert the new cache line
	oldestIndex := -1
	oldestAge := -1
	for i, l := range lru.lines {
		if !l.Valid { // Empty spot found
			lru.lines[i] = *line
			lru.lines[i].Age = 0
			lru.lines[i].Valid = true
			return
		} else if l.Age > oldestAge {
			oldestAge = l.Age
			oldestIndex = i
		}
	}

	// Replace the oldest (least recently used) cache line if no empty spot
	if oldestIndex != -1 {
		lru.lines[oldestIndex] = *line
		lru.lines[oldestIndex].Age = 0
	}
}

// Update the age of a cache line, resetting it since it's been accessed
func (lru *LRU) Update(line *CacheLine) {
	// Increment age of all valid lines
	for i := range lru.lines {
		if lru.lines[i].Valid {
			lru.lines[i].Age++
		}
	}

	// Reset age of the accessed line
	for i, l := range lru.lines {
		if l.Tag == line.Tag && l.Valid {
			lru.lines[i].Age = 0
			break
		}
	}
}

// Evict the least recently used cache line based on age
func (lru *LRU) Evict() int {
	oldestIndex := -1
	oldestAge := -1
	for i, l := range lru.lines {
		if l.Valid && l.Age > oldestAge {
			oldestAge = l.Age
			oldestIndex = i
		}
	}
	if oldestIndex != -1 {
		evictedTag := lru.lines[oldestIndex].Index
		lru.lines[oldestIndex].Valid = false // Mark as evicted
		return evictedTag
	}
	return -1 // Indicate no eviction occurred (should not happen if cache is full)
}
