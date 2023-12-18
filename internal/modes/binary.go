package modes

// ReadBits reads a set of bits.
func ReadBits(data []byte, bitCursor uint64, count uint8) uint64 {
	var out uint64

	shift := int(count)

	for {
		byteIdx := bitCursor / 8        //nolint: gomnd
		bitIdx := bitCursor % 8         //nolint: gomnd
		shift = shift + int(bitIdx) - 8 //nolint: gomnd

		mask := byte(0xff) >> bitIdx //nolint: gomnd

		if shift > 0 {
			out |= (uint64((data[byteIdx] & mask)) << shift)
		} else {
			out |= (uint64(data[byteIdx]) >> -shift)
		}

		if shift <= 0 {
			return out
		}

		bitCursor = bitCursor + 8 - bitIdx //nolint: gomnd
	}
}

// WriteBits writes a set of bits.
func WriteBits(data []byte, toWrite uint64, bitCursor uint64, count uint8) {
	shift := int(count)

	for {
		byteIdx := bitCursor / 8        //nolint: gomnd
		bitIdx := byte(bitCursor % 8)   //nolint: gomnd
		shift = shift - 8 + int(bitIdx) //nolint: gomnd

		mask := byte(0xff << (8 - bitIdx)) //nolint: gomnd

		var byteWrite byte
		if shift > 0 {
			byteWrite = byte(toWrite >> shift)
		} else {
			mask = 0xff >> (8 + shift) //nolint: gomnd
			byteWrite = byte(toWrite << -shift)
		}

		(data)[byteIdx] &= mask

		(data)[byteIdx] &= byteWrite

		bitCursor = bitCursor + 8 - uint64(bitIdx) //nolint: gomnd

		if shift <= 0 {
			return
		}
	}
}
