package cache

import "container/heap"

type LFUHeap []*CacheLine

func (lfu LFUHeap) Len() int {
	return len(lfu)
}

func (lfu LFUHeap) Less(i, j int) bool {
	return lfu[i].Freq < lfu[j].Freq
}

func (lfu LFUHeap) Swap(i, j int) {
	lfu[i], lfu[j] = lfu[j], lfu[i]
	lfu[i].LFUIndex = i
	lfu[j].LFUIndex = j
}

func (lfu *LFUHeap) Push(x interface{}) {
	lfuLen := len(*lfu)
	item := x.(*CacheLine)
	item.LFUIndex = lfuLen
	*lfu = append(*lfu, item)
}

func (lfu *LFUHeap) Pop() interface{} {
	old := *lfu
	n := len(old)
	item := old[n-1]
	item.LFUIndex = -1
	*lfu = old[0 : n-1]
	return item
}

func (lfu *LFUHeap) Update(line *CacheLine) {
	line.Freq++
	heap.Fix(lfu, line.LFUIndex)
}

type LFU struct {
	Capacity int
	Heap     *LFUHeap
}

func NewLFU(capacity int) *LFU {
	lfuCache := &LFU{
		Capacity: 0,
		Heap:     &LFUHeap{},
	}
	heap.Init(lfuCache.Heap)
	return lfuCache
}

func (c *LFU) Init() {
	heap.Init(c.Heap)
}

func (c *LFU) Insert(line *CacheLine) {
	if c.Capacity == 0 {
		return
	}
	if c.Capacity == c.Heap.Len() {
		c.Evict()
	}
	heap.Push(c.Heap, line)
	// TODO fix
}

func (c *LFU) Evict() (index int) {
	if c.Heap.Len() == 0 {
		return -1
	}
	return heap.Pop(c.Heap).(*CacheLine).Index
}

func (c *LFU) Update(line *CacheLine) {
	(*c.Heap).Update(line)
}

func (c *LFU) GetCapacity() int {
	return c.Capacity
}
