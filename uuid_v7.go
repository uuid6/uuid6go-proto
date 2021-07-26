package uuid

import (
	"bytes"
	"crypto/rand"
	"time"
)

type UUIDv7 uuidBase

// Generator is a primary structure that allows you to create new UUIDv7
// SubsecondPrecisionLength is a number of bits to carry sub-second information. On a systems with a high-resulution sub-second precision (like x64 Linux/Windows/MAC) it can go up to 48 bits
// NodePrecisionBits is a number of bits for node information. If this is not set to 0 [default], a Node variable must be set to a non-zero value
// Node information about the node generating UUIDs
// CounterPrecisionBits how many bits are dedicated to a counter. If two UUIDs were generated at the same time an internal counter would increase, distinguishing those UUIDs
type UUIDv7Generator struct {
	SubsecondPrecisionLength int

	NodePrecisionBits int
	Node              uint64

	CounterPrecisionBits int

	//Internal constants used during the generation process.
	currentTs       []byte
	counter         uint64
	currentPosition int
}

// UUIDv7FromBytes creates a new UUIDv7 from a slice of bytes and returns an error, if an array length does not equal 16.
func UUIDv7FromBytes(b []byte) (uuid UUIDv7, err error) {
	var utmp uuidBase
	err = utmp.UnmarshalBinary(b)
	return UUIDv7(utmp), err
}

// Timestamp returns unix epoch stored in the struct without millisecond precision
func (u UUIDv7) Timestamp() uint64 {
	bytes := [16]byte(u)
	tmp := toUint64(bytes[0:5])
	tmp = tmp >> 4 //We are off by 4 last bits of the byte there.
	return tmp
}

// Timestamp returns unix epoch stored in the struct without millisecond precision
func (u UUIDv7) Time() time.Time {
	bytes := [16]byte(u)
	tmp := toUint64(bytes[0:5])
	tmp = tmp >> 4 //We are off by 4 last bits of the byte there.
	return time.Unix(int64(tmp), 0)
}

// Ver returns a version of UUID, 07 in this case
func (u UUIDv7) Ver() uint16 {
	bytes := [16]byte(u)
	var tmp uint16 = uint16(bytes[6:7][0])
	tmp = tmp >> 4 //We are off by 4 last bits of the byte there.
	return tmp
}

// Var doing something described in the draft, but I don't know what
func (u UUIDv7) Var() uint16 {
	bytes := [16]byte(u)
	var tmp uint16 = uint16(bytes[8:9][0])
	tmp = tmp >> 6 //We are off by 4 last bits of the byte there.
	return tmp
}

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

// UUIDv7FromBytes creates a new UUIDv7 from a slice of bytes and returns an error, if an array length does not equal 16.
func (u *UUIDv7Generator) Next() (uuid UUIDv7) {
	var retval UUIDv7

	//Getting current time
	then := time.Now()
	tsBytes := timeToBytes(then)

	//Copy unix timestamp where it suppose to go, bits 0 - 36
	for i := 0; i < 36; i++ {
		retval = retval.setBit(35-i, getBit(tsBytes, 63-i))
	}

	//Getting nanoseconds
	var precisitonBytes []byte
	//Use the below to test your counter
	// precisitonBytes := []byte{0xff, 0xff}
	if u.SubsecondPrecisionLength != 0 {
		precisitonBytes, _ = encodeDecimal(float64(then.Nanosecond()), u.SubsecondPrecisionLength)
	}

	//Adding sub-second precision length
	u.currentPosition = retval.stack(u.currentPosition, precisitonBytes, u.SubsecondPrecisionLength)

	//Checks if we are going to use counter at all, or we don't need it
	useCounter := false

	//If we are using precision and bytes on the last tick equal current bytes,
	//use counter and append it.
	//But do this only if we are using precision
	if u.SubsecondPrecisionLength != 0 {
		if bytes.Equal(u.currentTs, precisitonBytes) {
			u.counter++
			useCounter = true
		} else {
			u.counter = 0
		}
		u.currentTs = precisitonBytes
	}

	//If we are using the counter, it goes right after precision bytes
	if useCounter {
		//counter bits
		u.currentPosition = retval.stack(u.currentPosition, toBytes(u.counter), u.CounterPrecisionBits)

	}

	//Adding node data after bytes
	if u.NodePrecisionBits != 0 {
		u.currentPosition = retval.stack(u.currentPosition, toBytes(u.Node), u.NodePrecisionBits)
	}

	//Create some random crypto data for the tail end
	rnd := make([]byte, 16)
	rand.Read(rnd)

	//Copy crypto data from the array to the end of the GUID
	cnt := 0
	limit := indexer(u.currentPosition)
	for i := 127; i > limit; i-- {
		//Ommiting bits 48-51 and 64, 65. Those contain version and variant information
		if i == 48 || i == 49 || i == 50 || i == 51 || i == 64 || i == 65 {
			continue
		}
		bit := getBit(rnd, cnt)
		cnt++
		retval = retval.setBit(i, bit)
	}

	//Adding version data [0111 = 7]
	retval = retval.setBit(48, false)
	retval = retval.setBit(49, true)
	retval = retval.setBit(50, true)
	retval = retval.setBit(51, true)

	//Adding variant data [10]
	retval = retval.setBit(64, true)
	retval = retval.setBit(65, false)

	return UUIDv7(retval)
}

func (u UUIDv7) ToString() string {
	return uuidBase(u).ToString()
}

func (u UUIDv7) ToMicrosoftString() string {
	return uuidBase(u).ToMicrosoftString()
}

func (u UUIDv7) ToBinaryString() string {
	return uuidBase(u).ToBinaryString()
}

func (u UUIDv7) ToBitArray() []bool {
	return uuidBase(u).ToBitArray()
}
