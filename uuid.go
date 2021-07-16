package uuid

type UUIDv6 uuidBase

type UUIDv8 uuidBase

// UUIDv6FromBytes creates a new UUIDv6 from a slice of bytes and returns an error, if an array length does not equal 16.
func UUIDv6FromBytes(b []byte) (uuid UUIDv6, err error) {
	var utmp uuidBase
	err = utmp.UnmarshalBinary(b)
	return UUIDv6(utmp), err
}

// UUIDv8FromBytes creates a new UUIDv8 from a slice of bytes and returns an error, if an array length does not equal 16.
func UUIDv8FromBytes(b []byte) (uuid UUIDv8, err error) {
	var utmp uuidBase
	err = utmp.UnmarshalBinary(b)
	return UUIDv8(utmp), err
}
