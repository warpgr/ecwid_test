package app

import "context"

type IPParserApp interface {
	Init() error
	Run(ctx context.Context) <-chan struct{}
	Stop() error
}
