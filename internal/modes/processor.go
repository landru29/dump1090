package modes

import (
	"encoding/hex"
	"fmt"

	"github.com/landru29/dump1090/internal/transport"
)

// Process is the data processor.
type Process struct {
	transporters []transport.Transporter
}

// New creates a data processor.
func New(transporters []transport.Transporter) *Process {
	return &Process{
		transporters: transporters,
	}
}

// Process implements source.Processor the interface.
func (p Process) Process(data []byte) {
	fmt.Println(hex.EncodeToString(data))
}
