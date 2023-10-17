package main

import (
	"fmt"
	"strings"
)

type protocolConfig struct {
	addr   string
	format string
}

const defaultUDPformat = "nmea"

func (p *protocolConfig) String() string {
	return fmt.Sprintf("%s@%s", p.format, p.addr)
}

func (p *protocolConfig) Set(str string) error {
	splitter := strings.Split(str, "@")
	if len(splitter) > 1 {
		*p = protocolConfig{
			format: splitter[0],
			addr:   strings.Join(splitter[1:], "@"),
		}

		return nil
	}

	*p = protocolConfig{
		format: defaultUDPformat,
		addr:   str,
	}

	return nil
}

func (p *protocolConfig) Type() string {
	return "protocol configuration"
}
