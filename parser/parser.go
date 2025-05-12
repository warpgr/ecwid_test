package parser

import "github.com/warpgr/ecwid_test/store"

func NewIPParser(st store.IPStore) IPParser {
	return &ipParser{st: st}
}

type ipParser struct {
	st store.IPStore
}

func (p *ipParser) ExtractIPs(buff []byte) []byte {
	start := 0
	for i, b := range buff {
		if b == '\n' {
			p.st.StoreIP(p.parseIPv4(buff[start:i]))
			start = i + 1
		}
	}

	return buff[start:]
}

func (p *ipParser) parseIPv4(line []byte) uint32 {
	var ip, acc uint32

	for _, b := range line {
		switch {
		case b >= '0' && b <= '9':
			acc = acc*10 + uint32(b-'0')
		case b == '.':
			ip = (ip << 8) | acc
			acc = 0
		}
	}

	ip = (ip << 8) | acc

	return ip
}
