package modes

func ReadBits(data []byte, bitCursor uint64, count uint8) uint64 {
	var out uint64

	shift := int(count)

	for {
		byteIdx := bitCursor / 8
		bitIdx := bitCursor % 8
		shift = shift + int(bitIdx) - 8

		mask := byte(0xff) >> bitIdx

		if shift > 0 {
			out = out | (uint64((data[byteIdx] & mask)) << shift)
		} else {
			out = out | (uint64(data[byteIdx]) >> -shift)
		}

		if shift <= 0 {
			return out
		}

		bitCursor = bitCursor + 8 - bitIdx
	}
}

func WriteBits(data []byte, toWrite uint64, bitCursor uint64, count uint8) {
	shift := int(count)

	for {
		byteIdx := bitCursor / 8
		bitIdx := byte(bitCursor % 8)
		shift = shift - 8 + int(bitIdx)

		mask := byte(0xff << (8 - bitIdx))

		var byteWrite byte
		if shift > 0 {
			byteWrite = byte(toWrite >> shift)
		} else {
			mask = 0xff >> (8 + shift)
			byteWrite = byte(toWrite << -shift)
		}

		(data)[byteIdx] = (data)[byteIdx] & mask

		(data)[byteIdx] = (data)[byteIdx] | byteWrite

		bitCursor = bitCursor + 8 - uint64(bitIdx)

		if shift <= 0 {
			return
		}
	}
}
