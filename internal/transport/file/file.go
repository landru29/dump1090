// Package file is the file display.
package file

import (
	"context"
	"fmt"
	"os"

	"github.com/landru29/dump1090/internal/dump"
	"github.com/landru29/dump1090/internal/serialize"
)

// Transporter is the file transporter.
type Transporter struct {
	serializer serialize.Serializer
	fileDesc   *os.File
}

// Transport implements the transport.Transporter interface.
func (t Transporter) Transport(ac *dump.Aircraft) error {
	data, err := t.serializer.Serialize(ac)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	t.fileDesc.WriteString(fmt.Sprintf("%s\n", string(data)))

	return nil
}

func New(ctx context.Context, filename string, serializer serialize.Serializer) (*Transporter, error) {
	if serializer == nil {
		return nil, fmt.Errorf("no valid formater")
	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		_ = f.Close()
	}()

	return &Transporter{
		serializer: serializer,
		fileDesc:   f,
	}, nil
}
