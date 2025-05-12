package store

type IPStore interface {
	StoreIP(ip uint32)
	StoreIPs(ips []uint32)
	UniqueCount() uint32
}
