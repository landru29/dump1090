// Package file is the file display.
package file

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/landru29/dump1090/internal/model"
	"github.com/landru29/dump1090/internal/serialize"
)

// Transporter is the file transporter.
type Transporter struct {
	serializer serialize.Serializer
	fileDesc   *os.File
}

// Transport implements the transport.Transporter interface.
func (t Transporter) Transport(ac *model.Aircraft) error {
	data, err := t.serializer.Serialize(ac)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	_, err = t.fileDesc.WriteString(fmt.Sprintf("%s\n", string(data)))

	return err
}

// New creates a file transporter.
func New(ctx context.Context, filename string, serializer serialize.Serializer) (*Transporter, error) {
	if serializer == nil {
		return nil, fmt.Errorf("no valid formater")
	}

	file, err := os.OpenFile(filepath.Clean(filename), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o600) //nolint: gomnd
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()

		_ = file.Close()
	}()

	return &Transporter{
		serializer: serializer,
		fileDesc:   file,
	}, nil
}

// String implements the transport.Transporter interface.
func (t Transporter) String() string {
	return "file"
}
