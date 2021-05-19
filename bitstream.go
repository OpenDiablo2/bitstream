package bitstream

import (
	"io"
)

const (
	bitsPerByte = 8
	byteMask    = 0xFF
	bitMask     = 0x01
)

type endianness bool

const (
	LittleEndian endianness = false
	BigEndian    endianness = true
)

func New(rs io.ReadSeeker) *BitStream {
	result := &BitStream{
		readSeeker:   rs,
		bytePosition: 0,
		bitPosition:  0,
	}

	result.setDefaultOptions()

	return result
}

func FromBytes(bytes ...byte) *BitStream {
	return New(nil).FromBytes(bytes...)
}

type BitStream struct {
	readSeeker   io.ReadSeeker
	endianness       // determines which end the bits are read from the byte (from biggest end or smallest end)
	bytePosition int // the byte index in the stream, 0 to length of stream bytes
	bitPosition  int // the bit index within the current byte, 0 to 7
	Options options
}

type options struct{
	ReadBeyondEOF bool // allows returning 0's when reading bits past EOF
}

func (bs *BitStream) setDefaultOptions() {
	bs.Options.ReadBeyondEOF = true
}

func (bs *BitStream) FromBytes(bytes ...byte) *BitStream {
	bs.readSeeker = &byteSeeker{bytes: bytes}
	return bs
}

func (bs *BitStream) ReadBit() (bool, error) {
	if bs.readSeeker == nil {
		return false, io.EOF
	}

	if _, err := bs.readSeeker.Seek(int64(bs.bytePosition), io.SeekStart); err != nil {
		if !bs.Options.ReadBeyondEOF {
			err = nil
		}

		return false, err
	}

	tmp := []byte{0}

	position := bs.bitPosition // we store a copy, it gets altered during the read
	if numRead, err := bs.readSeeker.Read(tmp); numRead < 1 || err != nil {
		if bs.Options.ReadBeyondEOF {
			err = nil
		}

		return false, err
	}

	bs.OffsetBitPosition(+1)

	shift := 0

	switch bs.endianness {
	case LittleEndian:
		shift = position
	case BigEndian:
		shift = bitsPerByte - position
	}

	return ((tmp[0] >> shift) & bitMask) > 0, nil
}

func (bs *BitStream) ReadBits(n int) Bits {
	bits := make(Bits, n)

	for idx := 0; idx < n; idx++ {
		b, err := bs.ReadBit()

		if err != nil {
			bits = bits[:idx]
			break
		}

		bits[idx] = b
	}

	return bits
}

func (bs *BitStream) Position() int {
	if bs.bytePosition < 0 {
		bs.bytePosition = 0
	}

	return bs.bytePosition
}

func (bs *BitStream) SetPosition(i int) *BitStream {
	bs.bytePosition = i

	if bs.bytePosition < 0 {
		bs.bytePosition = 0
	}

	return bs
}

func (bs *BitStream) OffsetPosition(i int) {
	bs.SetPosition(bs.Position() + i)
}

func (bs *BitStream) BitPosition() int {
	return bs.bitPosition
}

func (bs *BitStream) SetBitPosition(i int) *BitStream {
	// corner case, can't go back any further
	if i < 0 && bs.bytePosition == 0 {
		i = 0
	}

	// going negative when not at first byte
	for i < 0 && bs.bytePosition > 0 {
		i += bitsPerByte
		bs.OffsetPosition(-1)
	}

	// going further than current byte
	bs.OffsetPosition(i / bitsPerByte)

	// ensure bit position within 0..7
	bs.bitPosition = i % bitsPerByte

	return bs
}

func (bs *BitStream) OffsetBitPosition(i int) *BitStream {
	bs.SetBitPosition(bs.BitPosition() + i)

	return bs
}

func (bs *BitStream) SetLittleEndian() *BitStream {
	bs.endianness = LittleEndian
	return bs
}

func (bs *BitStream) SetBigEndian() *BitStream {
	bs.endianness = BigEndian
	return bs
}
