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

func Start(deviceIndex uint32, gain int, frequency uint32, enableAGC bool, filename string, callbackMessage func(*Message), callbackAircraft func(*Aircraft)) error {
	var filenameCString *C.char

	eventMessage = make(chan *Message, eventSize)
	eventAircraft = make(chan *Aircraft, eventSize)

	go func() {
		for {
			msg := <-eventMessage
			callbackMessage(msg)
		}
	}()

	go func() {
		for {
			ac := <-eventAircraft
			callbackAircraft(ac)
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
