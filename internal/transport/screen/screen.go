// Package screen is the screen display.
package screen

import (
	"fmt"

	"github.com/landru29/dump1090/internal/model"
	"github.com/landru29/dump1090/internal/serialize"
)

// Transporter is the screen transporter.
type Transporter struct {
	serializer serialize.Serializer
}

// Transport implements the transport.Transporter interface.
func (t Transporter) Transport(ac *model.Aircraft) error {
	if ac == nil {
		return nil
	}

	data, err := t.serializer.Serialize(ac)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	fmt.Printf("%s\n", string(data)) //nolint: forbidigo

	return nil
}

// New creates  a screen serializer.
func New(serializer serialize.Serializer) (*Transporter, error) {
	if serializer == nil {
		return nil, fmt.Errorf("no valid formater")
	}

	return &Transporter{
		serializer: serializer,
	}, nil
}

// String implements the transport.Transporter interface.
func (t Transporter) String() string {
	return "screen"
}
