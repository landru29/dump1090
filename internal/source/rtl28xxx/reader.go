package rtl28xxx

import (
	"context"
	"errors"
	"io"
	"unsafe"

	"github.com/landru29/dump1090/internal/source"
)

/*
  #cgo LDFLAGS: -lrtlsdr -lm
  #include "rtlsdr.h"

*/
import "C"

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
	newContext := context.WithValue(ctx, deviceInContext{}, r.processor)

	rtlContext := unsafe.Pointer(C.newContext(unsafe.Pointer(&newContext)))

	for {
		data := make([]byte, 1024)
		cnt, err := r.reader.Read(data)
		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return err
		}

		cstr := (*C.uchar)(unsafe.Pointer(C.CString(string(data[:cnt]))))

		C.rtlsdrCallback(cstr, C.uint(cnt), rtlContext)

		C.free(unsafe.Pointer(cstr))
	}
}
