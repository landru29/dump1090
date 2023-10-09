// Package udp is the udp display.
package udp

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/landru29/dump1090/internal/dump"
	"github.com/landru29/dump1090/internal/serialize"
)

// Transporter is the udp transporter.
type Transporter struct {
	serializer serialize.Serializer
	address    *net.UDPAddr
	udpServer  net.PacketConn
}

func New(ctx context.Context, serializer serialize.Serializer, addr string) (*Transporter, error) {
	splitter := strings.Split(addr, ":")

	udpServer, err := net.ListenPacket("udp", fmt.Sprintf(":%s", splitter[len(splitter)-1]))
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		_ = udpServer.Close()
	}()

	address, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		return nil, err
	}

	return &Transporter{
		serializer: serializer,
		address:    address,
		udpServer:  udpServer,
	}, nil
}

// Transport implements the transport.Transporter interface.
func (t *Transporter) Transport(ac *dump.Aircraft) error {
	data, err := t.serializer.Serialize(ac)
	if err != nil {
		return err
	}

	_, err = t.udpServer.WriteTo(data, t.address)

	return err
}
