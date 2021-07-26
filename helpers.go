package uuid

import (
	"encoding/binary"
	"math"
	"time"
)

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

// toUint64 converts []byte to uint64
func toBytes(data uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(data))
	return b
}

func timeToBytes(t time.Time) []byte {
	tsBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(tsBytes, uint64(t.Unix()))
	return tsBytes
}
