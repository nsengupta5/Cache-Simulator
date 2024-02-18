package cache

// This file implement the Round Robin replacement policy
// The Round Robin policy evicts the next cache line in the set
// in a round robin fashion based on the round robin index. It
// is used as the default replacement policy in the cache simulator

type RR struct {
	Capacity int
	Index    int
}

func NewRR(capacity int) *RR {
	return &RR{
		Capacity: capacity,
		Index:    0,
	}
}

// No action required to insert a line
func (c *RR) Insert(line *CacheLine) {
}

// Evict returns the index of the next line to evict
// based on the round robin index
func (c *RR) Evict() (index int) {
	evictIndex := c.Index
	// Wrap around the index if it exceeds the capacity
	c.Index = (c.Index + 1) % c.Capacity
	return evictIndex
}

// No action required to update a line
func (c *RR) Update(line *CacheLine) {
}
