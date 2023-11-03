package modes

func (e ExtendedSquitter) ChecksumSquitter() error {
	remainder := ChecksumSquitter(e.ModeS.Raw[:len(e.ModeS.Raw)-3])

	if remainder != e.ParityInterrogator {
		return ErrWrongCRC
	}

	return nil
}

func ChecksumSquitter(input []byte) uint32 {
	data := make([]byte, len(input)+3)
	copy(data, input)

	// generator is 25 bits (1111111111111010000001001).
	// x^24 + x^23 +x^22 +x^21 +x^20 +x^19 +x^18 +x^17 +x^16 +x^15 +x^14 +x^13 +x^12 +x^10 +x^3 + 1
	generator := uint64(0x01fff409)

	for idx := 0; idx < len(input)*8; idx++ {
		bitIdx := byte(idx % 8)
		byteIdx := idx / 8
		mask := byte(0x80) >> bitIdx

		if data[byteIdx]&mask != 0 {
			val := ReadBits(data, uint64(idx), 25)
			WriteBits(data, val^generator, uint64(idx), 25)
		}
	}

	return uint32(ReadBits(data, uint64(len(input)*8), 24))
}
