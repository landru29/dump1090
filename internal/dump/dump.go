package dump

/*
  #cgo LDFLAGS: -lrtlsdr -lm
  #include "dump1090.h"
*/
import "C"
import (
	"context"
	"errors"
	"unsafe"
)

const eventSize = 100

var (
	eventMessage  chan *Message
	eventAircraft chan *Aircraft
)

func EventMessage() chan *Message {
	return eventMessage
}

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

	C.free(unsafe.Pointer(ac))

	eventAircraft <- &aircraft
}

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

			//SetCountry(ac)

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
		C.uint8_t(map[bool]C.uint8_t{
			true:  1,
			false: 0,
		}[enableAGC]),
		filenameCString,
		loopInt,
	); ret != 0 {
		return errors.New("there were errors")
	}

	<-ctx.Done()

	return nil
}
