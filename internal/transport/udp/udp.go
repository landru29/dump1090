// Package udp is the udp display.
package udp

import (
	"fmt"
	"net"
	"sync"

	"github.com/landru29/dump1090/internal/dump"
	"github.com/landru29/dump1090/internal/serialize"
)

// Transporter is the udp transporter.
type Transporter struct {
	serializer serialize.Serializer
	mutex      sync.Mutex

	udpServer net.PacketConn

	openConn []net.Addr
}

func New(serializer serialize.Serializer, port int) (*Transporter, error) {
	udpServer, err := net.ListenPacket("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	output := Transporter{
		serializer: serializer,
		udpServer:  udpServer,
	}

	go func() {
		for {
			buf := make([]byte, 1024)
			_, addr, err := udpServer.ReadFrom(buf)
			if err != nil {
				continue
			}

			output.mutex.Lock()
			output.openConn = append(output.openConn, addr)
			output.mutex.Unlock()
		}
	}()

	return &Transporter{
		serializer: serializer,
		udpServer:  udpServer,
	}, nil
}

// Transport implements the transport.Transporter interface.
func (t *Transporter) Transport(ac *dump.Aircraft) error {
	data, err := t.serializer.Serialize(ac)
	if err != nil {
		return err
	}

	t.mutex.Lock()
	for _, addr := range t.openConn {
		_, _ = t.udpServer.WriteTo(data, addr)
	}
	defer t.mutex.Unlock()

	return nil
}
