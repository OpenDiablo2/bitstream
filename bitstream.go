// Package bitstream provides a stream reader implementation that can
// read data that is not byte aligned.
package bitstream

import (
	"bytes"
	"fmt"
	"io"
)

const (
	bitsPerByte       = 8
	bitsPerWord       = bitsPerByte
	bitsPerDoubleWord = bitsPerWord << 1
	//nolint:deadcode,unused,varcheck // may be used
	bitsPerQuadWord = bitsPerDoubleWord << 1

	bitMask = 0x01
)

// New creates a new BitStream using the given io.ReadSeeker
func New(rs io.ReadSeeker) *BitStream {
	bs := &BitStream{stream: rs}

	return bs
}

// FromBytes creates a new BitStream with the given bytes
func FromBytes(data ...byte) *BitStream {
	return New(nil).FromBytes(data...)
}

// Copy creates a deep copy of the source BitStream
func Copy(src *BitStream) *BitStream {
	return src.Copy()
}

// BitStream is used for reading structured data that is not byte aligned.
// It can read one or many bits and yield a BitInterpreter.
// The BitInterpreter can then be used to interpret the bits as another type.
type BitStream struct {
	stream      io.ReadSeeker
	bitPosition int // the bit index within the current byte, 0 to 7
	bitsRead    int // the number of bits read since the last seek
	unitsToRead int
	Options     options
}

type endianness int

// endianess types
const (
	LittleEndian endianness = iota
	BigEndian
)

type options struct {
	endianness // determines which end the bits are read from the byte (from biggest end or smallest end)
}

// FromBytes yields a new BitStream, using the given bytes as the stream source
func (bs *BitStream) FromBytes(b ...byte) *BitStream {
	bs.stream = bytes.NewReader(b)
	return bs
}

// Copy creates a copy of this bitstream, including the data.
// We do some funky stuff here regarding byte and bit offset to mimic the
// original implementation.
func (bs *BitStream) Copy() *BitStream {
	currentPosition, _ := bs.Seek(0, io.SeekCurrent)
	currentBitPosition := bs.bitPosition
	absoluteBitPosition := (int(currentPosition) * bitsPerByte) + currentBitPosition

	length, _ := bs.Seek(0, io.SeekEnd)

	_, _ = bs.Seek(0, io.SeekStart)
	buf := make([]byte, length)

	if bs.stream != nil {
		_, _ = bs.stream.Read(buf)
	}

	dst := FromBytes(buf...)

	dst.SetPosition(0).OffsetBitPosition(absoluteBitPosition)

	return dst
}

// BitsRead returns a number of readed bits
func (bs *BitStream) BitsRead() int {
	return bs.bitsRead
}

var tmpBit = make([]byte, 1)

// readBit reads a single bit from the stream source. It will always yield a boolean,
// even if the stream is empty or at the end of the file. However, if BitStream.Options.ReadBeyondEOF
// is false, readBit will return an io.EOF error during this read.
//
// It is also important to note that this read operation mutates
// the BitStream BytePosition and BitPosition.
func (bs *BitStream) readBit() (bool, error) {
	if bs.stream == nil {
		return false, io.EOF
	}

	bp := bs.bitPosition // we store a copy, it gets altered during the read

	if numRead, err := bs.stream.Read(tmpBit); numRead < 1 || err != nil {
		return false, fmt.Errorf("error reading bits: %w", err)
	}

	bs.bitsRead++
	bs.OffsetBitPosition(+1)

	// we only seek to the next byte in the stream when we
	// read the last bit (index 7) of the current byte
	if bp < bitsPerByte {
		// otherwise, we seek backward by one byte so that
		// the next read yields the same byte
		_, _ = bs.Seek(-1, io.SeekCurrent)
	}

	shift := 0

	switch bs.Options.endianness {
	case LittleEndian:
		shift = bp
	case BigEndian:
		shift = bitsPerByte - bp
	}

	return ((tmpBit[0] >> shift) & bitMask) > 0, nil
}

// readBits will read n bits into a BitInterpreter. If the BitStream.Options.ReadBeyondEOF
// option is false, the resultant BitInterpreter will be truncated to only contain bits
// that were successfully read before encountering the end of file.
func (bs *BitStream) readBits(n int) (Bits, error) {
	bits := make(Bits, n) // preallocate

	// read each bit, one by one
	for idx := 0; idx < n; idx++ {
		b, err := bs.readBit()
		// if there is an error (EOF), we truncate if ReadBeyondEOF is false
		if err != nil {
			return bits, err
		}

		bits[idx] = b
	}

	return bits, nil
}

// Seek sets the byte position within the stream.
func (bs *BitStream) Seek(offset int64, whence int) (int64, error) {
	if bs.stream == nil {
		return 0, io.EOF
	}

	result, err := bs.stream.Seek(offset, whence)
	if err != nil {
		return 0, fmt.Errorf("error seeking Bitstream: %w", err)
	}

	return result, nil
}

// Position returns the byte-position within the stream
func (bs *BitStream) Position() int {
	p, _ := bs.Seek(0, io.SeekCurrent)
	return int(p)
}

// SetPosition sets the byte position within the stream.
// The final position will be a positive integer.
func (bs *BitStream) SetPosition(i int) *BitStream {
	_, _ = bs.Seek(int64(i), io.SeekStart)

	return bs
}

// OffsetPosition will offset the current position by the given integer.
// The final position will be a positive integer.
func (bs *BitStream) OffsetPosition(i int) int {
	position, _ := bs.Seek(int64(i), io.SeekCurrent)
	return int(position)
}

// BitPosition returns the current bit index that the reader will read from.
// This is relative to the current byte within the stream, it is not an
// absolute bit position within the entire stream.
func (bs *BitStream) BitPosition() int {
	return bs.bitPosition
}

// SetBitPosition sets the bit position.
// It's IMPORTANT to understand that this is NOT relative to the current position! This is relative to
// bit position 0 of the current byte!
//
// Example: setting to -1 is the same as calling OffsetPosition(-1) and then SetBitPosition(7)
func (bs *BitStream) SetBitPosition(i int) *BitStream {
	position, _ := bs.Seek(0, io.SeekCurrent)

	// corner case, can't go back any further
	if i < 0 && position <= 0 {
		i = 0
	}

	// going negative when not at first byte
	for i < 0 && position > 0 {
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
	bs.Options.endianness = LittleEndian
	return bs
}

// SetBigEndian makes the BitStream read bits from the current byte from most-significant to least-significant.
func (bs *BitStream) SetBigEndian() *BitStream {
	bs.Options.endianness = BigEndian
	return bs
}

// Next sets the integer count for the next "unit" of data read.
// ex:
//		instance.Next(2).Bytes().AsUInt64() will read 2 bytes and interpret as uint64
//		instance.Next(4).Bits().AsInt() will read 4 bits
func (bs *BitStream) Next(count int) *BitStream {
	bs.unitsToRead = count

	return bs
}

// Bits will read a number of bits from the stream into a Response
//
// NOTE: The number is specified by by calling bitstream.Next
//
// example:
//
// val, err := bitstream.Next(2).Bytes()
func (bs *BitStream) Bits() Response {
	bits, err := bs.readBits(bs.unitsToRead)

	return Response{bits, err}
}

// Bytes will read (bs.unitsToRead * 8) bits into a Response
func (bs *BitStream) Bytes() Response {
	bits, err := bs.readBits(bs.unitsToRead * bitsPerByte)

	return Response{bits, err}
}
