// Package screen is the screen display.
package screen

import (
	"fmt"

	"github.com/landru29/dump1090/internal/dump"
	"github.com/landru29/dump1090/internal/serialize"
)

// Transporter is the screen transporter.
type Transporter struct {
	serializer serialize.Serializer
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

	fmt.Printf("%s\n", string(data))

	return nil
}

func New(serializer serialize.Serializer) *Transporter {
	return &Transporter{
		serializer: serializer,
	}
}
