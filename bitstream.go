package bitstream

import (
	"io"
)

const (
	bitsPerByte = 8
	bitMask     = 0x01
)

// New creates a new BitStream using the given io.ReadSeeker
func New(rs io.ReadSeeker) *BitStream {
	bs := &BitStream{ReadSeeker: rs}

	bs.setDefaultOptions()

	return bs
}

// FromBytes creates a new BitStream with the given bytes
func FromBytes(bytes ...byte) *BitStream {
	return New(nil).FromBytes(bytes...)
}

// Copy creates a deep copy of the source BitStream
func Copy(src *BitStream) *BitStream {
	return src.Copy()
}

// BitStream is used for reading structured data that is not byte aligned.
// It can read one or many bits and yield a BitInterpreter.
// The BitInterpreter can then be used to interpret the bits as another type.
type BitStream struct {
	io.ReadSeeker
	bitPosition  int // the bit index within the current byte, 0 to 7
	bitsRead int // the number of bits read since the last seek
	Options      options
}

type endianness int

const (
	LittleEndian endianness = iota
	BigEndian
)

type options struct {
	endianness         // determines which end the bits are read from the byte (from biggest end or smallest end)
	ReadBeyondEOF bool // allows returning 0's when reading bits past EOF
}

func (bs *BitStream) setDefaultOptions() {
	bs.Options.ReadBeyondEOF = true
}

// FromBytes yields a new BitStream, using the given bytes as the stream source
func (bs *BitStream) FromBytes(bytes ...byte) *BitStream {
	bs.ReadSeeker = &byteSeeker{bytes: bytes}
	return bs
}

// Copy creates a copy of this bitstream, including the data
func (bs *BitStream) Copy() *BitStream {
	currentPosition, _ := bs.Seek(0, io.SeekCurrent)
	currentBitPosition := bs.bitPosition
	bitsRead := bs.bitsRead
	length, _ := bs.Seek(0, io.SeekEnd)

	_, _ = bs.Seek(0, io.SeekStart)
	buf := make([]byte, length)

	_, _ = bs.Read(buf)

	dst := FromBytes(buf...)

	bs.SetPosition(int(currentPosition)).SetBitPosition(currentBitPosition)
	bs.bitsRead = bitsRead

	return dst
}

func (bs *BitStream) BitsRead() int {
	return bs.bitsRead
}

// ReadBit reads a single bit from the stream source. It will always yield a boolean,
// even if the stream is empty or at the end of the file. However, if BitStream.Options.ReadBeyondEOF
// is false, ReadBit will return an io.EOF error during this read.
//
// It is also important to note that this read operation mutates
// the BitStream BytePosition and BitPosition.
var tmpBit = make([]byte, 1)

func (bs *BitStream) ReadBit() (bool, error) {
	if bs.ReadSeeker == nil {
		return false, io.EOF
	}

	bs.bitsRead++

	position := bs.bitPosition // we store a copy, it gets altered during the read
	if numRead, err := bs.ReadSeeker.Read(tmpBit); numRead < 1 || err != nil {
		if bs.Options.ReadBeyondEOF {
			err = nil
			bs.bitsRead--
		}

		return false, err
	}

	bs.OffsetBitPosition(+1)

	shift := 0

	switch bs.Options.endianness {
	case LittleEndian:
		shift = position
	case BigEndian:
		shift = bitsPerByte - position
	}

	return ((tmpBit[0] >> shift) & bitMask) > 0, nil
}

// ReadBits will read n bits into a BitInterpreter. If the BitStream.Options.ReadBeyondEOF
// option is false, the resultant BitInterpreter will be truncated to only contain bits
// that were successfully read before encountering the end of file.
func (bs *BitStream) ReadBits(n int) Bits {
	bits := make(Bits, n)

	for idx := 0; idx < n; idx++ {
		b, err := bs.ReadBit()

		if err != nil {
			if !bs.Options.ReadBeyondEOF {
				bits = bits[:idx]
			}

			break
		}

		bits[idx] = b
	}

	return bits
}

// Seek sets the byte position within the stream.
func (bs *BitStream) Seek(offset int64, whence int) (int64, error) {
	if bs.ReadSeeker == nil {
		return 0, io.EOF
	}

	return bs.ReadSeeker.Seek(offset, whence)
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

// ReadBytes is a helper method for reading a slice of bytes from the bitstream
func (bs *BitStream) ReadByte() byte {
	return bs.ReadBits(bitsPerByte).AsByte()
}

// ReadBytes is a helper method for reading a slice of bytes from the bitstream
func (bs *BitStream) ReadBytes(count int) []byte {
	bytes := make([]byte, count)

	for idx := 0; idx < count; idx++ {
		bytes[idx] = bs.ReadByte()
	}

	return bytes
}
