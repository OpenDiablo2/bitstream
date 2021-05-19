package bitstream

const (
	twosComplimentNegativeOne = 4294967295
)

type Bits []bool

func (b Bits) AsByte() byte {
	return byte(b.AsUInt())
}

func (b Bits) AsInt8() int8 {
	return int8(makeSigned32(uint32(b.AsUInt8()), len(b)))
}

func (b Bits) AsUInt8() uint8 {
	return uint8(b.AsInt8())
}

func (b Bits) AsInt16() int16 {
	return int16(makeSigned32(uint32(b.AsUInt16()), len(b)))
}

func (b Bits) AsUInt16() uint16 {
	return uint16(b.AsUInt())
}

func (b Bits) AsInt32() int32 {
	return makeSigned32(b.AsUInt32(), len(b))
}

func (b Bits) AsUInt32() uint32 {
	return uint32(b.AsUInt())
}

func (b Bits) AsInt64() int64 {
	return int64(b.AsInt32()) // lol hack
}

func (b Bits) AsUInt64() uint64 {
	return uint64(b.AsUInt())
}

func (b Bits) AsInt() int {
	return int(makeSigned32(uint32(b.AsUInt()), len(b)))
}

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