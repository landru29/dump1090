// Package net is the net display (TCP or UDP).
package net

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strings"
	"sync"

	localerrors "github.com/landru29/dump1090/internal/errors"
	"github.com/landru29/dump1090/internal/logger"
	"github.com/landru29/dump1090/internal/model"
	"github.com/landru29/dump1090/internal/serialize"
	"github.com/pkg/errors"
)

const (
	errNoValidFormater localerrors.Error = "no valid formater"
)

// Transporter is the udp transporter.
type Transporter struct {
	clients  []io.WriteCloser
	formater serialize.Serializer
	mutex    sync.Mutex
	log      *slog.Logger
}

// New creates a net transporter.
func New(
	ctx context.Context,
	formater map[string]serialize.Serializer,
	conf ProtocolConfig,
	log *slog.Logger,
) (*Transporter, error) {
	if formater == nil {
		return nil, errNoValidFormater
	}

	if log == nil {
		return nil, logger.ErrMissingLogger
	}

	serial, found := formater[conf.format]
	if !found {
		return nil, fmt.Errorf("serializer %s not found", conf.format)
	}

	output := &Transporter{
		formater: serial,
		log:      log,
	}

	switch conf.direction {
	case protocolBind:
		return output, output.Bind(ctx, conf.protocolType, conf.addr)
	case protocolDial:
		return output, output.Dial(ctx, conf.protocolType, conf.addr)
	}

	return nil, fmt.Errorf("unknown %s: specify 'dial' or 'bind'", conf.direction)
}

// Bind is the net binder.
func (t *Transporter) Bind(ctx context.Context, pType protocolType, addr string) error {
	splitter := strings.Split(addr, ":")
	port := splitter[len(splitter)-1]

	tcpServer, err := net.Listen(string(pType), fmt.Sprintf(":%s", port))
	if err != nil {
		return err
	}

	log := t.log.With("type", string(pType), "port", port)

	log.Info("net listening")

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

				log.Info("accept connection", "from", conn.RemoteAddr().String())

				t.mutex.Lock()
				t.clients = append(t.clients, conn)
				t.mutex.Unlock()
			}
		}
	}()

	return nil
}

// Dial is the net dialer.
func (t *Transporter) Dial(ctx context.Context, pType protocolType, addr string) error {
	log := t.log.With("type", string(pType), "to", addr)

	log.Info("dialing")

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

	log.Info("connection accepted")

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
			t.log.Error("client error", "msg", err)

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

// String implements the transport.Transporter interface.
func (t *Transporter) String() string {
	return "net"
}
