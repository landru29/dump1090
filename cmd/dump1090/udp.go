package main

import (
	"fmt"
	"strings"
)

type udpConfig struct {
	addr   string
	format string
}

const defaultUDPformat = "nmea"

func (u *udpConfig) String() string {
	return fmt.Sprintf("%s@%s", u.format, u.addr)
}

func (u *udpConfig) Set(str string) error {
	splitter := strings.Split(str, "@")
	if len(splitter) > 1 {
		*u = udpConfig{
			format: splitter[0],
			addr:   strings.Join(splitter[1:], "@"),
		}

		return nil
	}

	*u = udpConfig{
		format: defaultUDPformat,
		addr:   str,
	}

	return nil
}

func (u *udpConfig) Type() string {
	return "udp configuration"
}
