// Package net is the net display (TCP or UDP).
package net

import (
	"context"
	"fmt"
	nativenet "net"

	"github.com/landru29/dump1090/internal/dump"
	"github.com/landru29/dump1090/internal/serialize"
)

// Transporter is the udp transporter.
type Transporter struct {
	server   nativenet.Conn
	formater serialize.Serializer
}

func New(ctx context.Context, formater serialize.Serializer, addr string, tcpudp string) (*Transporter, error) {
	if formater == nil {
		return nil, fmt.Errorf("no valid formater")
	}

	server, err := nativenet.Dial(tcpudp, addr)
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		_ = server.Close()
	}()

	return &Transporter{
		//address:   address,
		formater: formater,
		server:   server,
	}, nil
}

// Transport implements the transport.Transporter interface.
func (t *Transporter) Transport(ac *dump.Aircraft) error {
	data, err := t.formater.Serialize(ac)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	data = append(data, '\n')

	_, err = t.server.Write(data)

	return err
}
