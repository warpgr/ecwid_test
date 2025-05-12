package reader

import "context"

type IPReaderEngine interface {
	Run(ctx context.Context)
}
