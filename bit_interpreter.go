package bitstream

import "math"

const (
	twosComplimentNegativeOne = 4294967295
)

// BitInterpreter is responsible for converting an array of booleans to various number formats.
type BitInterpreter []bool

// Bits is an alias for BitInterpreter, for brevity.
type Bits = BitInterpreter

// AsByte interprets the bits as a byte
func (b Bits) AsBool() bool {
	return b.AsInt() > 0
}

// AsByte interprets the bits as a byte
func (b Bits) AsByte() byte {
	return byte(b.AsUInt())
}

// AsBytes interprets the bits as a slice of bytes
func (b Bits) AsBytes() []byte {
	numBits := len(b)
	numBytes := int(math.Ceil(float64(numBits) / float64(bitsPerByte)))
	result := make([]byte, numBytes)

	for idx := 0; idx < numBytes; idx++ {
		startBit, stopBit := idx*bitsPerByte, (idx+1)*bitsPerByte
		if stopBit > numBits {
			stopBit = numBits
		}
		result[idx] = Bits(b[startBit:stopBit]).AsByte()
	}

	return result
}
// AsInt8 interprets the bits as a signed 8-bit integer
func (b Bits) AsInt8() int8 {
	return int8(makeSigned32(uint32(b.AsUInt8()), len(b)))
}

// AsUInt8 interprets the bits as a usnigned 8-bit integer
func (b Bits) AsUInt8() uint8 {
	return uint8(b.AsInt8())
}

// AsInt16 interprets the bits as a signed 16-bit integer
func (b Bits) AsInt16() int16 {
	return int16(makeSigned32(uint32(b.AsUInt16()), len(b)))
}

// AsUInt16 interprets the bits as an unsugned 16-bit integer
func (b Bits) AsUInt16() uint16 {
	return uint16(b.AsUInt())
}

// AsInt32 interprets the bits as a signed 32-bit integer
func (b Bits) AsInt32() int32 {
	return makeSigned32(b.AsUInt32(), len(b))
}

// AsUInt32 interprets the bits as an unsigned 32-bit integer
func (b Bits) AsUInt32() uint32 {
	return uint32(b.AsUInt())
}

// AsInt64 interprets the bits as a signed 64-bit integer
func (b Bits) AsInt64() int64 {
	return int64(b.AsInt32()) // lol hack
}

// AsUInt64 interprets the bits as an unsigned 64-bit integer
func (b Bits) AsUInt64() uint64 {
	return uint64(b.AsUInt())
}

// AsInt interprets the bits as a signed integer
func (b Bits) AsInt() int {
	return int(makeSigned32(uint32(b.AsUInt()), len(b)))
}

// AsUInt interprets the bits as an unsigned integer
func (b Bits) AsUInt() uint {
	result := uint(0)

	for idx := 0; idx < len(b); idx++ {
		if b[idx] {
			result = result | (1 << idx)
		}
	}

	return result
}

func makeSigned32(unsignedValue uint32, signBitIndex int) int32 {
	if signBitIndex == 0 {
		return 0
	}

	// If its a single bit, a unsignedValue of 1 is -1 automagically
	if signBitIndex == 1 {
		return -int32(unsignedValue)
	}

	signMask := uint32(1 << (signBitIndex - 1))

	// If there is no sign bit, return the unsignedValue as is
	if (unsignedValue & signMask) == 0 {
		return int32(unsignedValue)
	}

	// We need to extend the signed bit out so that the negative unsignedValue
	// representation still works with the 2s compliment rule.
	result := uint32(twosComplimentNegativeOne)

	for i := byte(0); i < byte(signBitIndex); i++ {
		if ((unsignedValue >> uint(i)) & 1) == 0 {
			result -= uint32(1 << uint(i))
		}
	}

	// Force casting to a signed unsignedValue
	return int32(result)
}