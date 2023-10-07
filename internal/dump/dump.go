package dump

/*
  #cgo LDFLAGS: -lrtlsdr -lm
  #include "dump1090.h"
*/
import "C"
import (
	"errors"
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
func goSendAircraft(ac *C.aircraft) {
	if ac == nil {
		return
	}

	aircraft := newAircraft(ac)

	eventAircraft <- &aircraft
}

func Start(
	deviceIndex uint32,
	gain int,
	frequency uint32,
	enableAGC bool,
	filename string,
	evtMessage chan *Message,
	evtAircraft chan *Aircraft,
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

	if ret := C.startProcess(
		C.uint32_t(deviceIndex),
		C.int(gain),
		C.uint32_t(frequency),
		C.uint8_t(map[bool]C.uint8_t{
			true:  1,
			false: 0,
		}[enableAGC]),
		filenameCString,
	); ret != 0 {
		return errors.New("there were errors")
	}

	close(eventMessage)
	close(eventAircraft)

	return nil
}
