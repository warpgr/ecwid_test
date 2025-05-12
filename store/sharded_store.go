package store

import (
	"sync/atomic"
)

const (
	numShards     = 64
	bitsPerShard  = 1 << 32 / numShards // = 2^32 / 64 = 2^26 bits
	wordsPerShard = bitsPerShard / 64   // = 2^26 / 64 = 2^20 words
)

type atomicShard struct {
	bitset []uint64      // length = wordsPerShard
	count  atomic.Uint32 // unique count in this shard
}

func NewShardedAtomicStore() IPStore {
	s := &shardedAtomicStore{}
	for i := range s.shards {
		s.shards[i] = &atomicShard{
			bitset: make([]uint64, wordsPerShard),
		}
	}
	return s
}

type shardedAtomicStore struct {
	shards [numShards]*atomicShard
}

func (s *shardedAtomicStore) StoreIP(ip uint32) {
	// Getting high 6 bits as shard index.
	shardIdx := ip >> (32 - 6)
	shard := s.shards[shardIdx]

	// Getting low 26 bits as offset in the shard.
	offset := ip & ((1 << (32 - 6)) - 1)

	// Getting mask by low 8 bits.
	bit := offset & 63
	mask := uint64(1) << bit

	// Getting index in the bitset of shard.
	word := offset >> 6
	ptr := &shard.bitset[word]

	for {
		old := atomic.LoadUint64(ptr)
		// Checking checking is current bit set.
		if old&mask != 0 {
			return
		}

		if atomic.CompareAndSwapUint64(ptr, old, old|mask) {
			shard.count.Add(1)
			return
		}
	}
}

func (s *shardedAtomicStore) StoreIPs(ips []uint32) {
	for _, ip := range ips {
		s.StoreIP(ip)
	}
}

func (s *shardedAtomicStore) UniqueCount() uint32 {
	var total uint32
	for _, shard := range s.shards {
		total += shard.count.Load()
	}
	return total
}
