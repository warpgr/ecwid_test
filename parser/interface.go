package parser

type IPParser interface {
	ExtractIPs(buff []byte) []byte
}
