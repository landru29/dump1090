// Package net is the net display (TCP or UDP).
package net

import (
	"context"
	"fmt"
	"io"
	"net"
	nativenet "net"
	"strings"
	"sync"

	"github.com/landru29/dump1090/internal/dump"
	"github.com/landru29/dump1090/internal/serialize"
	"github.com/pkg/errors"
)

// Transporter is the udp transporter.
type Transporter struct {
	clients  []io.WriteCloser
	formater serialize.Serializer
	mutex    sync.Mutex
}

func New(ctx context.Context, formater map[string]serialize.Serializer, conf ProtocolConfig) (*Transporter, error) {
	if formater == nil {
		return nil, fmt.Errorf("no valid formater")
	}

	serial, found := formater[conf.format]
	if !found {
		return nil, fmt.Errorf("serializer %s not found", conf.format)
	}

	output := &Transporter{
		formater: serial,
	}

	switch conf.direction {
	case protocolBind:
		return output, output.Bind(ctx, conf.protocolType, conf.addr)
	case protocolDial:
		return output, output.Dial(ctx, conf.protocolType, conf.addr)
	}

	return nil, fmt.Errorf("unknown %s: specify 'dial' or 'bind'", conf.direction)
}

func (t *Transporter) Bind(ctx context.Context, pType protocolType, addr string) error {
	splitter := strings.Split(addr, ":")

	tcpServer, err := net.Listen("tcp", fmt.Sprintf(":%s", splitter[len(splitter)-1]))
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				t.close()
				_ = tcpServer.Close()
				return
			default:
				var conn nativenet.Conn
				// Wait for a connection.
				conn, err = tcpServer.Accept()
				if err != nil {
					fmt.Printf("ERROR: %s\n", err)

					return
				}

				t.mutex.Lock()
				t.clients = append(t.clients, conn)
				t.mutex.Unlock()
			}
		}
	}()

	return nil
}

func (t *Transporter) Dial(ctx context.Context, pType protocolType, addr string) error {
	server, err := nativenet.Dial(string(pType), addr)
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		t.close()
	}()

	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.clients = append(t.clients, server)

	return nil
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

	t.mutex.Lock()
	defer t.mutex.Unlock()

	var (
		globalErr error
		hasError  bool
	)
	for _, client := range t.clients {
		_, err = client.Write(data)

		if err != nil {
			hasError = true
			globalErr = errors.Wrap(globalErr, err.Error())
		}
	}

	if hasError {
		return globalErr
	}

	return nil
}

func (t *Transporter) close() {
	t.mutex.Lock()
	for _, client := range t.clients {
		_ = client.Close()
	}
	t.mutex.Unlock()
}
