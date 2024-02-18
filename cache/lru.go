package cache

type LRU struct {
	capacity   int
	cache      map[int]*CacheLine // Maps tags to pointers to cache lines
	head, tail *CacheLine         // Pointers to head and tail of the doubly-linked list
}

func NewLRU(capacity int) *LRU {
	return &LRU{
		capacity: capacity,
		cache:    make(map[int]*CacheLine),
	}
}

// Insert inserts a new cache line
func (lru *LRU) Insert(line *CacheLine) {
	if existingLine, exists := lru.cache[line.Tag]; exists {
		// Move the accessed line to the front of the list
		lru.remove(existingLine)
		lru.addToFront(existingLine)
		return
	}

	// Insert the new cache line at the front of the list
	lru.cache[line.Tag] = line
	lru.addToFront(line)
}

// Update moves the accessed cache line to the front of the list,
// marking it as most recently used
func (lru *LRU) Update(line *CacheLine) {
	if existingLine, exists := lru.cache[line.Tag]; exists {
		lru.remove(existingLine)
		lru.addToFront(existingLine)
	}
}

// Evict returns the index of the least recently used cache line
func (lru *LRU) Evict() int {
	if lru.tail == nil {
		return -1 // Cache is empty
	}

	// Remove the least recently used cache line from the list
	// which is the tail of the list
	evictedIndex := lru.tail.Index
	delete(lru.cache, lru.tail.Tag)
	lru.remove(lru.tail)
	return evictedIndex
}

// remove removes a cache line from the doubly-linked list
func (lru *LRU) remove(line *CacheLine) {
	if line.Prev != nil {
		line.Prev.Next = line.Next
	} else {
		lru.head = line.Next
	}
	if line.Next != nil {
		line.Next.Prev = line.Prev
	} else {
		lru.tail = line.Prev
	}
}

// addToFront adds a cache line to the front of the doubly-linked list
func (lru *LRU) addToFront(line *CacheLine) {
	line.Next = lru.head
	line.Prev = nil
	if lru.head != nil {
		lru.head.Prev = line
	}
	lru.head = line
	if lru.tail == nil {
		lru.tail = line
	}
}
