package uuid

import (
	"encoding/binary"
	"fmt"
)

type uuidBase [16]byte

type UUIDv6 uuidBase

type UUIDv7 uuidBase

/*
bytes      	0               1               2               3
		    0                   1                   2                   3
		    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
		   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		   |                            unixts                             |
		   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

bytes      	4               5               6               7
		   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		   |unixts |       subsec_a        |  ver  |       subsec_b        |
		   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

bytes      	8               9               10              11
		   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		   |var|                   subsec_seq_node                         |
		   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

bytes      	12              13              14              15
		   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		   |                       subsec_seq_node                         |
		   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

*/

type UUIDv8 uuidBase

// UUIDv6FromBytes creates a new UUIDv6 from a slice of bytes and returns an error, if an array length does not equal 16.
func UUIDv6FromBytes(b []byte) (uuid UUIDv6, err error) {
	var utmp uuidBase
	err = utmp.UnmarshalBinary(b)
	return UUIDv6(utmp), err
}

// UUIDv7FromBytes creates a new UUIDv7 from a slice of bytes and returns an error, if an array length does not equal 16.
func UUIDv7FromBytes(b []byte) (uuid UUIDv7, err error) {
	var utmp uuidBase
	err = utmp.UnmarshalBinary(b)
	return UUIDv7(utmp), err
}

// UUIDv8FromBytes creates a new UUIDv8 from a slice of bytes and returns an error, if an array length does not equal 16.
func UUIDv8FromBytes(b []byte) (uuid UUIDv8, err error) {
	var utmp uuidBase
	err = utmp.UnmarshalBinary(b)
	return UUIDv8(utmp), err
}

func (u UUIDv7) UnixTS() uint64 {
	bytes := [16]byte(u)
	var tmp uint64 = binary.BigEndian.Uint64(bytes[0:4])
	tmp = tmp >> 4 //We are off by 4 last bits of the byte there.
	return tmp
}

func (u UUIDv7) Ver() uint16 {
	bytes := [16]byte(u)
	var tmp uint16 = uint16(bytes[6:7][0])
	tmp = tmp >> 4 //We are off by 4 last bits of the byte there.
	return tmp
}

func (u UUIDv7) Var() uint16 {
	bytes := [16]byte(u)
	var tmp uint16 = uint16(bytes[8:9][0])
	tmp = tmp >> 6 //We are off by 4 last bits of the byte there.
	return tmp
}

func (u UUIDv7) Subseq() uint64 {
	bytes := [16]byte(u)
	var tmp uint64 = binary.BigEndian.Uint64(bytes[8:15])
	tmp = tmp & 0b0011_1111_1111_1111_1111_1111_1111_1111_1111_1111_1111_1111_1111_1111_1111_1111
	return tmp
}

func (u UUIDv7) SubsecA() uint16 {
	bytes := [16]byte(u)
	var tmp uint16 = binary.BigEndian.Uint16(bytes[4:5])
	tmp = tmp & 0b0000_1111_1111_1111
	return tmp
}

func (u UUIDv7) SubsecB() uint16 {
	bytes := [16]byte(u)
	var tmp uint16 = binary.BigEndian.Uint16(bytes[6:7])
	tmp = tmp & 0b0000_1111_1111_1111
	return tmp
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (uuid *uuidBase) UnmarshalBinary(data []byte) error {
	if len(data) != 16 {
		return fmt.Errorf("invalid UUID (got %d bytes)", len(data))
	}
	copy(uuid[:], data)
	return nil
}
