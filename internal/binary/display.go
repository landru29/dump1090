package binary

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

// DisplayBits displays a set of bits on a writer.
func DisplayBits(output io.Writer, dataByte any, leadingSpaces int, size int) {
	line := bytes.NewBuffer(nil)

	switch data := dataByte.(type) {
	case []byte:
		for _, elt := range data {
			line.WriteString(fmt.Sprintf("%08b", elt))
		}
	case uint64:
		line.WriteString(trimStringZero(fmt.Sprintf("%064b", data)))
	case uint32:
		line.WriteString(trimStringZero(fmt.Sprintf("%032b", data)))
	case uint16:
		line.WriteString(trimStringZero(fmt.Sprintf("%016b", data)))
	case uint8:
		line.WriteString(trimStringZero(fmt.Sprintf("%08b", data)))
	}

	fmt.Fprintf(output, "%s\n", segmentString(strings.Repeat(" ", leadingSpaces)+timStringSize(line.String(), size)))
}

func segmentString(str string) string {
	output := bytes.NewBuffer(nil)

	for idx, char := range str {
		if idx != 0 && (idx%8) == 0 { //nolint: gomnd
			output.WriteString(" ")
		}

		output.WriteRune(char)
	}

	return output.String()
}

func timStringSize(str string, size int) string {
	if size < 0 {
		return str
	}

	if size > len(str) {
		return strings.Repeat(" ", len(str)-size) + str
	}

	return str[len(str)-size:]
}

func trimStringZero(str string) string {
	output := bytes.NewBuffer(nil)

	leading := true

	for _, char := range str {
		if leading && char == '0' {
			output.WriteString(" ")

			continue
		}

		leading = false

		output.WriteByte(byte(char))
	}

	return output.String()
}
