package uuid

import (
	"encoding/binary"
	"math"
	"time"
)

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

// encodeDecimal takes nanoseconds and converts them to the binary-encoded arbitrary-precision
// byte array.
func encodeDecimal(sec float64, bits int) (val []byte, err error) {
	len := int(math.Log10(sec)) + 1
	sec = sec / math.Pow10(len)
	num := math.Pow(2, float64(bits))
	var part uint64 = uint64(sec * float64(num))
	val = make([]byte, 8)
	binary.BigEndian.PutUint64(val, part)
	return val, nil
}

// toUint64 converts []byte to uint64
func toUint64(data []byte) uint64 {
	var arr [8]byte
	copy(arr[len(arr)-len(data):], data)
	return binary.BigEndian.Uint64(arr[:])
}

func timeToBytes(t time.Time) []byte {
	tsBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(tsBytes, uint64(t.Unix()))
	return tsBytes
}
