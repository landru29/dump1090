package main

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/landru29/dump1090/internal/database"
	"github.com/landru29/dump1090/internal/model"
	"github.com/landru29/dump1090/internal/serialize"
	"github.com/landru29/dump1090/internal/transport"
	"github.com/landru29/dump1090/internal/transport/file"
	"github.com/landru29/dump1090/internal/transport/http"
	"github.com/landru29/dump1090/internal/transport/net"
	"github.com/landru29/dump1090/internal/transport/screen"
)

func provideTransporters(
	ctx context.Context,
	log *slog.Logger,
	availableSerializers []serialize.Serializer,
	serializers map[string]serialize.Serializer,
	aircraftDB *database.ElementStorage[model.ICAOAddr, model.Aircraft],
	httpConf httpConfig,
	udpConf net.ProtocolConfig,
	tcpConf net.ProtocolConfig,
	transportScreen string,
	transportFile string,
) ([]transport.Transporter, error) {
	transporters := []transport.Transporter{}

	if httpConf.addr != "" {
		httpTransport, err := http.New(ctx, httpConf.addr, httpConf.apiPath, aircraftDB, availableSerializers)
		if err != nil {
			return nil, err
		}

		transporters = append(transporters, httpTransport)

		log.Info("API", "addr", fmt.Sprintf("http://%s%s\n", httpConf.addr, httpConf.apiPath))
	}

	if udpConf.IsValid() {
		udpTransport, err := net.New(ctx, serializers, udpConf, log)
		if err != nil {
			return nil, err
		}

		transporters = append(transporters, udpTransport)
	}

	if tcpConf.IsValid() {
		tcpTransport, err := net.New(ctx, serializers, tcpConf, log)
		if err != nil {
			return nil, err
		}

		transporters = append(transporters, tcpTransport)
	}

	if transportScreen != "" {
		screenTransport, err := screen.New(serializers[transportScreen])
		if err != nil {
			return nil, err
		}

		transporters = append(transporters, screenTransport)
	}

	if transportFile != "" {
		splitter := strings.Split(transportFile, "@")
		if len(splitter) > 1 {
			fileTransport, err := file.New(ctx, strings.Join(splitter[1:], "@"), serializers[splitter[0]])
			if err != nil {
				return nil, err
			}

			transporters = append(transporters, fileTransport)
		}
	}

	if len(transporters) == 0 {
		transporters = append(transporters, screen.Transporter{})
	}

	return transporters, nil
}
