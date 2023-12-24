package binary

// ChecksumSquitter computes the checksum.
func ChecksumSquitter(input []byte) uint32 {
	data := make([]byte, len(input)+3) //nolint: gomnd
	copy(data, input)

	// generator is 25 bits (1111111111111010000001001).
	// x^24 + x^23 +x^22 +x^21 +x^20 +x^19 +x^18 +x^17 +x^16 +x^15 +x^14 +x^13 +x^12 +x^10 +x^3 + 1
	generator := uint64(0x01fff409) //nolint: gomnd

	for idx := 0; idx < len(input)*8; idx++ {
		bitIdx := byte(idx % 8)      //nolint: gomnd
		byteIdx := idx / 8           //nolint: gomnd
		mask := byte(0x80) >> bitIdx //nolint: gomnd

		if data[byteIdx]&mask != 0 {
			val := ReadBits(data, uint64(idx), 25)          //nolint: gomnd
			WriteBits(data, val^generator, uint64(idx), 25) //nolint: gomnd
		}
	}

	return uint32(ReadBits(data, uint64(len(input)*8), 24)) //nolint: gomnd
}