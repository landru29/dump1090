package nmea_test

import (
	"bytes"
	"testing"

	"github.com/landru29/dump1090/internal/nmea"
	"github.com/stretchr/testify/assert"
)

func TestCheckSum(t *testing.T) {
	sentences := []string{
		"!AIVDM,1,1,,B,177KQJ5000G?tO`K>RA1wUbN0TKH,0*5C",
		"!AIVDM,1,1,,A,15RTgt0PAso;90TKcjM8h6g208CQ,0*4A",
		"!AIVDM,1,1,,A,16SteH0P00Jt63hHaa6SagvJ087r,0*42",
	}

	for _, sentence := range sentences {
		value := sentence

		t.Run(value, func(t *testing.T) {
			fields := nmea.Fields(bytes.Split([]byte(value), []byte{','}))

			assert.Equal(t, value[len(value)-2:], fields.CheckSum())
		})
	}

}
