// Package net is the net display (TCP or UDP).
package net

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"

	localerrors "github.com/landru29/dump1090/internal/errors"
	"github.com/landru29/dump1090/internal/model"
	"github.com/landru29/dump1090/internal/serialize"
	"github.com/pkg/errors"
)

const (
	errNoValidFormater localerrors.Error = "no valid formater"
)

// Transporter is the udp transporter.
type Transporter struct {
	clients   []io.WriteCloser
	formater  serialize.Serializer
	mutex     sync.Mutex
	outWriter io.Writer
}

// New creates a net transporter.
func New(
	ctx context.Context,
	formater map[string]serialize.Serializer,
	conf ProtocolConfig,
	outputLog io.Writer,
) (*Transporter, error) {
	if formater == nil {
		return nil, errNoValidFormater
	}

	if outputLog == nil {
		outputLog = io.Discard
	}

	serial, found := formater[conf.format]
	if !found {
		return nil, fmt.Errorf("serializer %s not found", conf.format)
	}

	output := &Transporter{
		formater:  serial,
		outWriter: outputLog,
	}

	switch conf.direction {
	case protocolBind:
		return output, output.Bind(ctx, conf.protocolType, conf.addr, outputLog)
	case protocolDial:
		return output, output.Dial(ctx, conf.protocolType, conf.addr, outputLog)
	}

	return nil, fmt.Errorf("unknown %s: specify 'dial' or 'bind'", conf.direction)
}

// Bind is the net binder.
func (t *Transporter) Bind(ctx context.Context, pType protocolType, addr string, outputLog io.Writer) error {
	splitter := strings.Split(addr, ":")
	port := splitter[len(splitter)-1]

	tcpServer, err := net.Listen(string(pType), fmt.Sprintf(":%s", port))
	if err != nil {
		return err
	}

	fmt.Fprintf(outputLog, "Listening %s on port %s\n", strings.ToUpper(string(pType)), port)

	go func() {
		for {
			select {
			case <-ctx.Done():
				t.close()

				_ = tcpServer.Close()

				return
			default:
				var conn net.Conn
				// Wait for a connection.
				conn, err = tcpServer.Accept()
				if err != nil {
					fmt.Printf("ERROR: %s\n", err) //nolint: forbidigo

					return
				}

				fmt.Fprintf(outputLog, "Accept connection from %s\n", conn.RemoteAddr())

				t.mutex.Lock()
				t.clients = append(t.clients, conn)
				t.mutex.Unlock()
			}
		}
	}()

	return nil
}

// Dial is the net dialer.
func (t *Transporter) Dial(ctx context.Context, pType protocolType, addr string, outputLog io.Writer) error {
	fmt.Fprintf(outputLog, "Dialing %s to %s\n", strings.ToUpper(string(pType)), addr)

	server, err := net.Dial(string(pType), addr)
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

	fmt.Fprintf(outputLog, "Connection accepted\n")

	return nil
}

// Transport implements the transport.Transporter interface.
func (t *Transporter) Transport(ac *model.Aircraft) error {
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
			fmt.Fprintf(t.outWriter, "ERR: %s\n", err)

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
