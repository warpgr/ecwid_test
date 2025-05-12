package reader

import (
	"context"
	"io"
	"os"
	"sync"

	"github.com/warpgr/ecwid_test/parser"

	"github.com/sirupsen/logrus"
)

func NewSectionedIPReader(f *os.File, workers, chunkSize int, pr parser.IPParser) (IPReaderEngine, error) {
	// Computing file ranges to distribute between workers.
	offs, err := computeOffsets(f, workers)
	if err != nil {
		return nil, err
	}

	return &sectionedIPReader{
		file:      f,
		offsets:   offs,
		chunkSize: chunkSize,
		pr:        pr,
	}, nil
}

type sectionedIPReader struct {
	file      *os.File
	offsets   []int64
	chunkSize int

	pr parser.IPParser
}

func (r *sectionedIPReader) Run(ctx context.Context) {
	var wg sync.WaitGroup
	n := len(r.offsets) - 1
	wg.Add(n)

	const maxLineLength = 32

	for i := 0; i < n; i++ {
		start := r.offsets[i]
		end := r.offsets[i+1]
		go func(start, end int64) {
			defer wg.Done()
			sr := io.NewSectionReader(r.file, start, end-start)

			var leftover []byte
			buf := make([]byte, r.chunkSize+maxLineLength)

			for {
				select {
				case <-ctx.Done():
					return

				default:
					n, err := sr.Read(buf[len(leftover):])
					if n > 0 {
						leftover = r.pr.ExtractIPs(buf[:n+len(leftover)])
						copy(buf[0:len(leftover)], leftover)
					}

					if err != nil {
						if err == io.EOF {
							if len(leftover) != 0 {
								logrus.WithFields(logrus.Fields{
									"leftover": string(leftover),
								}).Warn("Can't handle the end of the file.")
							}
							return
						} else {
							logrus.WithError(err).WithFields(logrus.Fields{
								"start": start,
								"end":   end,
							}).Error("Error occurs during reading the file section.")
							return
						}
					}
				}
			}
		}(start, end)
	}

	wg.Wait()
}

func computeOffsets(f *os.File, workers int) ([]int64, error) {
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	size := info.Size()
	base := size / int64(workers)
	offs := make([]int64, workers+1)
	offs[0] = 0
	buf := make([]byte, 1)
	for i := 1; i < workers; i++ {
		pos := int64(i) * base
		for ; pos < size; pos++ {
			_, err := f.ReadAt(buf, pos)
			if err != nil {
				return nil, err
			}

			if buf[0] == '\n' {
				offs[i] = pos + 1
				break
			}
		}
		if offs[i] == 0 {
			offs[i] = size
		}
	}
	offs[workers] = size
	return offs, nil
}
