package mocks

import "sync"

type IPStoreMock struct {
	IPs   []uint32
	guard sync.Mutex
}

func (sm *IPStoreMock) StoreIP(ip uint32) {
	sm.guard.Lock()
	defer sm.guard.Unlock()

	sm.IPs = append(sm.IPs, ip)
}

func (sm *IPStoreMock) StoreIPs(ips []uint32) {
	sm.guard.Lock()
	defer sm.guard.Unlock()

	sm.IPs = append(sm.IPs, ips...)
}

func (sm *IPStoreMock) UniqueCount() uint32 {
	sm.guard.Lock()
	defer sm.guard.Unlock()

	return uint32(len(sm.IPs))
}
