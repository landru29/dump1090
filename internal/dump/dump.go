// Package dump is the embed c-code.
package dump

/*
  #cgo LDFLAGS: -lrtlsdr -lm
  #include "dump1090.h"
*/
import "C"

import (
	"context"
	"errors"
)

const eventSize = 100

var (
	eventMessage  chan *Message  //nolint: gochecknoglobals
	eventAircraft chan *Aircraft //nolint: gochecknoglobals
)

// EventMessage is the readonly variable.
func EventMessage() chan *Message {
	return eventMessage
}

// EventAircraft is the readonly variable.
func EventAircraft() chan *Aircraft {
	return eventAircraft
}

//export goSendMessage
func goSendMessage(msg *C.modesMessage) {
	if msg == nil {
		return
	}

	message := newMessage(msg)

	eventMessage <- &message
}

//export goSendAircraft
func goSendAircraft(msg *C.modesMessage, ac *C.aircraft) {
	if ac == nil {
		return
	}

	aircraft := newAircraft(ac, msg)

	eventAircraft <- &aircraft
}

// Start starts listening to the device stream.
func Start(
	ctx context.Context,
	deviceIndex uint32,
	gain int,
	frequency uint32,
	enableAGC bool,
	filename string,
	evtMessage chan *Message,
	evtAircraft chan *Aircraft,
	loop bool,
) error {
	var filenameCString *C.char

	eventMessage = make(chan *Message, eventSize)
	eventAircraft = make(chan *Aircraft, eventSize)

	go func() {
		for {
			msg := <-eventMessage

			if evtMessage != nil {
				evtMessage <- msg
			}
		}
	}()

	go func() {
		for {
			ac := <-eventAircraft

			if evtAircraft != nil {
				evtAircraft <- ac
			}
		}
	}()

	if filename != "" {
		filenameCString = C.CString(filename)
	}

	var loopInt C.int

	if loop {
		loopInt = C.int(1)
	}

	if ret := C.startProcess(
		C.uint32_t(deviceIndex),
		C.int(gain),
		C.uint32_t(frequency),
		map[bool]C.uint8_t{
			true:  1,
			false: 0,
		}[enableAGC],
		filenameCString,
		loopInt,
	); ret != 0 {
		return errors.New("there were errors")
	}

	<-ctx.Done()

	return nil
}
