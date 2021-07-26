package uuid

// getBit returns the value of a bit at a specified positon in UUID
func (b uuidBase) getBit(index int) bool {
	pos := index / 8
	j := uint(index % 8)
	j = 7 - j
	return (b[pos] & (uint8(1) << j)) != 0
}

// getBit returns the value of a bit at a specified positon in UUIDv7
func (b UUIDv7) getBit(index int) bool {
	return uuidBase(b).getBit(index)
}

// getBit returns the value of a bit at a specified positon in given byte array
func getBit(b []byte, index int) bool {
	pos := index / 8
	j := uint(index % 8)
	j = 7 - j
	return (b[pos] & (uint8(1) << j)) != 0
}

// setBit sets the value of a bit at a specified positon in UUID
func (b uuidBase) setBit(index int, value bool) uuidBase {
	pos := index / 8
	j := uint(index % 8)
	j = 7 - j
	if value {
		b[pos] |= (uint8(1) << j)
	} else {
		b[pos] &= ^(uint8(1) << j)
	}
	return b
}

// setBit sets the value of a bit at a specified positon in UUIDv7
func (b UUIDv7) setBit(index int, value bool) UUIDv7 {
	return UUIDv7(uuidBase(b).setBit(index, value))
}

// indexer returns updated index of a bit in the array of bits. It skips bits 48-51 and 64,65
// for those containt information about Version and Variant and can't be populated by the
// precision bits. It also omits first 36 bits of timestamp at the beginning of the GUID
func indexer(input int) int {
	out := 35 + input //Skip the TS block and start counting right after ts block
	if input > 11 {   //If we are bumbing into a ver block, skip it
		out += 4
	}

	if input > 23 { //If we are bumping into a var block
		out += 2
	}
	return out
}

// stack adds a chunk of bits, encoded as []byte at the selected started position, with respect to the timestamp, version and variant values.
func (b UUIDv7) stack(startingPosition int, value []byte, length int) (UUIDv7, int) {
	rettype, retval := (uuidBase(b)).stack(startingPosition, value, length)
	return UUIDv7(rettype), retval
}

// stack adds a chunk of bits, encoded as []byte at the selected started position, with respect to the timestamp, version and variant values.
func (b uuidBase) stack(startingPosition int, value []byte, length int) (uuidBase, int) {
	cnt := 0
	for i := startingPosition; i < startingPosition+length; i++ {
		bit := getBit(value, (len(value)*8-1)-cnt)
		b = b.setBit(indexer(i), bit)
		cnt++
	}
	return b, startingPosition + length
}
