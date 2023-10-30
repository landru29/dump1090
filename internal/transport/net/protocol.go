package net

import (
	"fmt"
	"strings"
)

type protocolType string

type protocolDirection string

const (
	protocolDial protocolDirection = "dial"
	protocolBind protocolDirection = "bind"

	protocolTypeTCP protocolType = "tcp"
	protocolTypeUDP protocolType = "udp"

	defaultProtocolFormat = "nmea"
)

type ProtocolConfig struct {
	addr         string
	format       string
	direction    protocolDirection
	protocolType protocolType
}

func NewProtocol(pType string) ProtocolConfig {
	return ProtocolConfig{
		protocolType: protocolType(pType),
	}
}

func (p *ProtocolConfig) String() string {
	return fmt.Sprintf(
		"%s/%s:%s@%s",
		p.direction,
		p.protocolType,
		p.format,
		p.addr,
	)
}

func (p *ProtocolConfig) Set(str string) error {
	actionSplitter := strings.Split(str, ">")
	switch len(actionSplitter) {
	case 1:
		format, addr := parseData(actionSplitter[0])

		p.format = format
		p.direction = protocolDial
		p.addr = addr

		return nil
	case 2:
		format, addr := parseData(actionSplitter[1])

		p.format = format
		p.direction = protocolDirection(actionSplitter[0])
		p.addr = addr

		return nil
	}

	return fmt.Errorf("wrong format %s (should be like dial>text@0.0.0.0:30003)", str)
}

func (p *ProtocolConfig) Type() string {
	return "protocol configuration"
}

func (p ProtocolConfig) IsValid() bool {
	return p.addr != ""
}

func parseData(str string) (string, string) {
	addr := "0.0.0.0:30003"
	format := defaultProtocolFormat

	if str != "" {
		addr = str
	}

	splitter := strings.Split(str, "@")
	if len(splitter) > 1 {
		format = splitter[0]
		addr = strings.Join(splitter[1:], "@")
	}

	return format, addr

}
