package modes

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/pkg/errors"
)

var checksumTable = [112]uint32{
	0x3935ea, 0x1c9af5, 0xf1b77e, 0x78dbbf, 0xc397db, 0x9e31e9, 0xb0e2f0, 0x587178,
	0x2c38bc, 0x161c5e, 0x0b0e2f, 0xfa7d13, 0x82c48d, 0xbe9842, 0x5f4c21, 0xd05c14,
	0x682e0a, 0x341705, 0xe5f186, 0x72f8c3, 0xc68665, 0x9cb936, 0x4e5c9b, 0xd8d449,
	0x939020, 0x49c810, 0x24e408, 0x127204, 0x093902, 0x049c81, 0xfdb444, 0x7eda22,
	0x3f6d11, 0xe04c8c, 0x702646, 0x381323, 0xe3f395, 0x8e03ce, 0x4701e7, 0xdc7af7,
	0x91c77f, 0xb719bb, 0xa476d9, 0xadc168, 0x56e0b4, 0x2b705a, 0x15b82d, 0xf52612,
	0x7a9309, 0xc2b380, 0x6159c0, 0x30ace0, 0x185670, 0x0c2b38, 0x06159c, 0x030ace,
	0x018567, 0xff38b7, 0x80665f, 0xbfc92b, 0xa01e91, 0xaff54c, 0x57faa6, 0x2bfd53,
	0xea04ad, 0x8af852, 0x457c29, 0xdd4410, 0x6ea208, 0x375104, 0x1ba882, 0x0dd441,
	0xf91024, 0x7c8812, 0x3e4409, 0xe0d800, 0x706c00, 0x383600, 0x1c1b00, 0x0e0d80,
	0x0706c0, 0x038360, 0x01c1b0, 0x00e0d8, 0x00706c, 0x003836, 0x001c1b, 0xfff409,
	0x000000, 0x000000, 0x000000, 0x000000, 0x000000, 0x000000, 0x000000, 0x000000,
	0x000000, 0x000000, 0x000000, 0x000000, 0x000000, 0x000000, 0x000000, 0x000000,
	0x000000, 0x000000, 0x000000, 0x000000, 0x000000, 0x000000, 0x000000, 0x000000,
}

func (e ExtendedSquitter) checksum() error {
	data := append(e.ModeS.Raw[:len(e.ModeS.Raw)-3], []byte{0, 0, 0}...)

	crc := checksumSquitter(data, uint16(8*len(data)))

	if crc != e.ModeS.ParityInterrogator {
		return errors.Wrap(ErrWrongCRC, fmt.Sprintf("expecting %08X, got %08X", e.ModeS.ParityInterrogator, crc))
	}

	return nil
}

func (e ExtendedSquitter) checksumv2() error {
	data := append(e.ModeS.Raw[:len(e.ModeS.Raw)-3], []byte{0, 0, 0, 0, 0, 0, 0, 0, 0}...)

	// x^24 + x^23 +x^22 +x^21 +x^20 +x^19 +x^18 +x^17 +x^16 +x^15 +x^14 +x^13 +x^12 +x^10 +x^3 + 1
	generator := uint64(0x01fff409)

	for idx := 0; idx < len(e.ModeS.Raw)-24; idx++ {
		bitIdx := byte(idx % 8)
		byteIdx := idx / 8
		mask := byte(80) >> bitIdx

		if data[byteIdx]&mask != 0 {
			var localdata uint64
			generatorMask := generator << (32 - bitIdx)
			_ = binary.Read(bytes.NewReader(data[byteIdx:]), binary.LittleEndian, localdata)
			localdata = localdata ^ generatorMask
			buf := bytes.NewBuffer(nil)
			binary.Write(buf, binary.BigEndian, localdata^generatorMask)

			for idx, b := range buf.Bytes() {
				data[idx+byteIdx] = b
			}
		}
	}

	return nil
}

func checksumSquitter(data []byte, bits uint16) uint32 {
	crc := uint32(0)
	offset := uint16(0)
	if len(data)*8 != 112 {
		offset = 112 - 56
	}

	for j := uint16(0); j < bits; j++ {
		byteIdx := j / 8
		bitIdx := byte(j % 8)
		bitmask := byte(1 << (7 - bitIdx))

		/* If bit is set, xor with corresponding table entry. */
		if data[byteIdx]&bitmask != 0 {
			crc = crc ^ checksumTable[j+offset]
		}
	}

	return crc
}
