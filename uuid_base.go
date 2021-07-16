package uuid

import "fmt"

type uuidBase [16]byte

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (uuid *uuidBase) UnmarshalBinary(data []byte) error {
	if len(data) != 16 {
		return fmt.Errorf("invalid UUID (got %d bytes)", len(data))
	}
	copy(uuid[:], data)
	return nil
}

// ToString returns UUID as a string in the following format
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
func (u uuidBase) ToString() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:5], u[5:7], u[7:9], u[9:11], u[11:16])
}

// ToMicrosoftString returns UUID as a string in the following format
// {XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX}
func (u uuidBase) ToMicrosoftString() string {
	return fmt.Sprintf("{%X-%X-%X-%X-%X}", u[0:5], u[5:7], u[7:9], u[9:11], u[11:16])
}

// ToBinaryString returns UUID as a string of binary digits grouped in 8
func (u uuidBase) ToBinaryString() string {
	var ret string
	for _, n := range u {
		ret += fmt.Sprintf("%0*b ", 8, n) // prints 00000000 11111101
	}
	return ret
}
