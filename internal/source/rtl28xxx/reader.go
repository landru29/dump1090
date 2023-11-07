package rtl28xxx

import (
	"context"
	"errors"
	"io"

	"github.com/landru29/dump1090/internal/source"
)

type Reader struct {
	reader io.Reader

	processor source.Processer
}

func NewReader(rd io.Reader, processor source.Processer) *Reader {
	return &Reader{
		reader:    rd,
		processor: processor,
	}
}

// Start implements the source.Starter interface.
func (r *Reader) Start(ctx context.Context) error {
	for {
		data := make([]byte, 1024)
		cnt, err := r.reader.Read(data)
		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return err
		}

		processRaw(ctx, data[:cnt], r.processor)
	}
}
