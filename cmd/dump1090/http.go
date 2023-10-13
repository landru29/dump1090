package main

import (
	"fmt"
	"strings"
)

type httpConfig struct {
	addr    string
	apiPath string
}

const defaulthttpAPIpath = "/api"

func (h *httpConfig) String() string {
	return fmt.Sprintf("%s%s", h.addr, h.apiPath)
}

func (h *httpConfig) Set(str string) error {
	splitter := strings.Split(str, "/")
	if len(splitter) > 1 {
		apiPath := strings.Join(splitter[1:], "/")

		*h = httpConfig{
			addr:    splitter[0],
			apiPath: "/" + apiPath,
		}

		return nil
	}

	*h = httpConfig{
		apiPath: defaulthttpAPIpath,
		addr:    str,
	}

	return nil
}

func (h *httpConfig) Type() string {
	return "http configuration"
}
