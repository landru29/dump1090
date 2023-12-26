package rtl28xxx

import (
	"context"
	"errors"
	"io"

	"github.com/landru29/dump1090/internal/processor"
	localcontext "github.com/landru29/dump1090/internal/source/context"
)

// Reader is device reader.
type Reader struct {
	reader io.Reader

	processors []processor.Processer
}

// NewReader creates a new device reader.
func NewReader(rd io.Reader, processors []processor.Processer) *Reader {
	return &Reader{
		reader:     rd,
		processors: processors,
	}
}

// Start implements the source.Starter interface.
func (r *Reader) Start(ctx context.Context) error {
	cContext := localcontext.New(ctx, r.processors)

	defer func() {
		localcontext.DisposeContext(cContext.Key)
	}()

	for {
		data := make([]byte, 1024) //nolint: gomnd

		cnt, err := r.reader.Read(data)
		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return err
		}

		processRaw(data[:cnt], cContext.Ccontext)
	}
}
