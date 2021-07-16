package uuid

// GetBit returns the value of a bit at a specified positon in UUID
func (b uuidBase) GetBit(index int) bool {
	pos := index / 8
	j := uint(index % 8)
	j = 7 - j
	return (b[pos] & (uint8(1) << j)) != 0
}

// GetBit returns the value of a bit at a specified positon in UUIDv7
func (b UUIDv7) GetBit(index int) bool {
	return uuidBase(b).GetBit(index)
}

// GetBit returns the value of a bit at a specified positon in given byte array
func GetBit(b []byte, index int) bool {
	pos := index / 8
	j := uint(index % 8)
	j = 7 - j
	return (b[pos] & (uint8(1) << j)) != 0
}

// SetBit sets the value of a bit at a specified positon in UUID
func (b uuidBase) SetBit(index int, value bool) uuidBase {
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

// SetBit sets the value of a bit at a specified positon in UUIDv7
func (b UUIDv7) SetBit(index int, value bool) UUIDv7 {
	return UUIDv7(uuidBase(b).SetBit(index, value))
}
