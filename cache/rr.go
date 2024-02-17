package cache

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

func (c *RR) Insert(line *CacheLine) {
	// do nothing
}

func (c *RR) Evict() (index int) {
	evictIndex := c.Index
	c.Index = (c.Index + 1) % c.Capacity
	return evictIndex
}

func (c *RR) Update(line *CacheLine) {
	// do nothing
}

func (c *RR) GetCapacity() int {
	return c.Capacity
}
