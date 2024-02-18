package cache

type ReplacementPolicy interface {
	Insert(line *CacheLine)
	Update(line *CacheLine)
	Evict() int
	GetCapacity() int
}

func GenerateReplacementPolicy(capacity int, policyName string) ReplacementPolicy {
	switch policyName {
	case "lru":
		return NewLRU(capacity)
	case "lfu":
		return NewLFU(capacity)
	case "rr":
		return NewRR(capacity)
	default:
		return nil
	}
}
