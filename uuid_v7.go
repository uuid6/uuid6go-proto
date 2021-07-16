package uuid

import (
	"bytes"
	"crypto/rand"
	"time"
)

type UUIDv7 uuidBase

type UUIDv7Generator struct {
	Precision int
	currentTs []byte
	counter   uint8
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

	//Getting current time
	then := time.Now()
	tsBytes := timeToBytes(then)

	//Getting nanoseconds
	var precisitonBytes []byte
	//Use the below to test your counter
	// precisitonBytes := []byte{0xff, 0xff}
	if u.Precision != 0 {
		precisitonBytes, _ = encodeDecimal(float64(then.Nanosecond()), u.Precision)
	}

	//Checks if we are going to use counter at all, or we don't need it
	useCounter := false

	//If we are using precision and bytes on the last tick equal current bytes,
	//use counter and append it.
	//But do this only if we are using precision
	if u.Precision != 0 {
		if bytes.Equal(u.currentTs, precisitonBytes) {
			u.counter++
			useCounter = true
		} else {
			u.counter = 0
		}
		u.currentTs = precisitonBytes
	}

	//Creating the array to return
	var retval UUIDv7

	//Limit will show you how much of the final lenght the precision bytes will take
	limit := 0
	//If we are using the counter, it goes right after precision bytes
	if useCounter {
		limit = u.Precision + 8
		precisitonBytes = append(precisitonBytes, u.counter)
	} else {
		limit = u.Precision
	}

	//Create some random crypto data for the tail end
	rnd := make([]byte, 16)
	rand.Read(rnd)

	//Copy unix timestamp where it suppose to go, bits 0 - 36
	for i := 0; i < 36; i++ {
		retval = retval.SetBit(35-i, GetBit(tsBytes, 63-i))
	}

	//Copy the bits from precision array to the return array
	//Indexer omits bits 48-51 and 64, 65. Those contain version and variant information
	//We have to start copying from the end of the bit array and go up,
	//because precision is not always aligned to BYTE and can contain empty zeroes
	//at the beginning.
	if u.Precision != 0 {
		for i := 0; i < limit; i++ {
			bit := GetBit(precisitonBytes, (len(precisitonBytes)*8-1)-i)
			retval = retval.SetBit(indexer(limit-i), bit)
		}
	}

	//Copy crypto data from the array to the end of the GUID
	cnt := 0
	limit = indexer(limit)
	for i := 127; i > limit; i-- {
		//Ommiting bits 48-51 and 64, 65. Those contain version and variant information
		if i == 48 || i == 49 || i == 50 || i == 51 || i == 64 || i == 65 {
			continue
		}
		bit := GetBit(rnd, cnt)
		cnt++
		retval = retval.SetBit(i, bit)
	}

	//Adding version data [0111 = 7]
	retval = retval.SetBit(48, false)
	retval = retval.SetBit(49, true)
	retval = retval.SetBit(50, true)
	retval = retval.SetBit(51, true)

	//Adding variant data [10]
	retval = retval.SetBit(64, true)
	retval = retval.SetBit(65, false)

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
