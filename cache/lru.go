package cache

type LRUNode struct {
	key   int
	prev  *LRUNode
	next  *LRUNode
	Lines []*CacheLine
}

type LRU struct {
	capacity   int
	keyNodeMap map[int]*LRUNode
	head       *LRUNode
	tail       *LRUNode
}

func NewLRU(capacity int) *LRU {
	lru := &LRU{
		capacity:   capacity,
		keyNodeMap: make(map[int]*LRUNode),
		head:       &LRUNode{},
		tail:       &LRUNode{},
	}
	lru.head.next = lru.tail
	lru.tail.prev = lru.head
	return lru
}
