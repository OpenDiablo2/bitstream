package bitstream

import (
	"io"
)

const (
	bitsPerByte = 8
	byteMask    = 0xFF
	bitMask     = 0x01
)

type endianness int

const (
	LittleEndian endianness = iota
	BigEndian
)

// New creates a new BitStream using the given io.ReadSeeker
func New(rs io.ReadSeeker) *BitStream {
	result := &BitStream{
		readSeeker:   rs,
		bytePosition: 0,
		bitPosition:  0,
	}

	result.setDefaultOptions()

	return result
}

// FromBytes creates a new BitStream with the given bytes
func FromBytes(bytes ...byte) *BitStream {
	return New(nil).FromBytes(bytes...)
}

// BitStream is used for reading structured data that is not byte aligned.
// It can read one or many bits and yield a BitInterpreter.
// The BitInterpreter can then be used to convert the bits into another number type.
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

// FromBytes yields a new BitStream, using the given bytes as the stream source
func (bs *BitStream) FromBytes(bytes ...byte) *BitStream {
	bs.readSeeker = &byteSeeker{bytes: bytes}
	return bs
}

// ReadBit reads a single bit from the stream source. It will always yield a boolean,
// even if the stream is empty or at the end of the file. However, if BitStream.Options.ReadBeyondEOF
// is false, ReadBit will return an io.EOF error during this read.
//
// It is also important to note that this read operation mutates
// the BitStream BytePosition and BitPosition.
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

// ReadBits will read n bits into a BitInterpreter. If the BitStream.Options.ReadBeyondEOF
// option is false, the resultant BitInterpretter will be truncated to only contain bits
// that were successfully read before encountering the end of file.
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

// Position returns the byte-position within the stream
func (bs *BitStream) Position() int {
	if bs.bytePosition < 0 {
		bs.bytePosition = 0
	}

	return bs.bytePosition
}

// SetPosition sets the byte position within the stream.
// The final position will be a positive integer.
func (bs *BitStream) SetPosition(i int) *BitStream {
	bs.bytePosition = i

	if bs.bytePosition < 0 {
		bs.bytePosition = 0
	}

	return bs
}

// OffsetPosition will offset the current position by the given integer.
// The final position will be a positive integer.
func (bs *BitStream) OffsetPosition(i int) {
	bs.SetPosition(bs.Position() + i)
}

// BitPosition returns the current bit index that the reader will read from.
// This is relative to the current byte within the stream, it is not an
// absolute bit position within the entire stream.
func (bs *BitStream) BitPosition() int {
	return bs.bitPosition
}

// SetBitPosition sets the bit position.
// NOTE: setting a bit position <0 or >7 will offset the byte position
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

// OffsetBitPosition will offset the current bit position, updating the byte position if (0 > i > 7)
func (bs *BitStream) OffsetBitPosition(i int) *BitStream {
	bs.SetBitPosition(bs.BitPosition() + i)

	return bs
}

// SetLittleEndian makes the BitStream read bits from the current byte from least-significant to most-significant.
func (bs *BitStream) SetLittleEndian() *BitStream {
	bs.endianness = LittleEndian
	return bs
}

// SetBigEndian makes the BitStream read bits from the current byte from most-significant to least-significant.
func (bs *BitStream) SetBigEndian() *BitStream {
	bs.endianness = BigEndian
	return bs
}
