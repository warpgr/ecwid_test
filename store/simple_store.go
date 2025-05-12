package store

import "sync"

func NewIPStore() IPStore {
	return newIpStore()
}

func newIpStore() *ipStore {
	const bitsetSize = 1 << 29

	return &ipStore{
		bitset: make([]byte, bitsetSize),
	}
}

type ipStore struct {
	bitset      []byte
	uniqueCount uint32

	guard sync.Mutex
}

func (st *ipStore) StoreIP(ip uint32) {
	st.store(ip)
}

func (st *ipStore) StoreIPs(ips []uint32) {
	for _, ip := range ips {
		st.store(ip)
	}
}

func (st *ipStore) UniqueCount() uint32 {
	return st.uniqueCount
}

func (st *ipStore) store(ip uint32) {
	st.guard.Lock()
	defer st.guard.Unlock()

	idx := ip >> 3
	mask := byte(1 << (ip & 7))

	if st.bitset[idx]&mask == 0 {
		st.bitset[idx] |= mask
		st.uniqueCount++
	}
}
