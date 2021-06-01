package bitstream

import "fmt"

// Writer is a stream writer, capable of writing data which is not byte-aligned.
// CAVEAT: the resulting byte buffer WILL be byte-aligned, as the underlying representation
// of the bits is in a byte slice.
type Writer struct {
	// all of the bytes written by this Writer
	bytes     []byte

	// bitBuffer is a single byte used for storing individual bits that are written.
	// When a byte-boundary is passed, this byte is appended to the bytes slice
	//
	// NOTE: calling Bytes() while bitOffset != 0 will append this bitBuffer as a final byte.
	bitBuffer byte

	// bitOffset is incremented every time a bit is written. It wraps modulo 8 (bits per byte).
	bitOffset int

	// endianness determines the order in which bits are written into the bitBuffer
	endianness
}


// Bytes returns a copy of the byte buffer
func (w *Writer) Bytes() []byte {
	bytes := append([]byte{}, w.bytes...)

	if w.bitOffset != 0 {
		bytes = append(bytes, w.bitBuffer)
	}

	return bytes
}

// Write the given args, yielding the number of bits written.
//
// NOTE: the arguments can be bool, byte, []byte, or Bits
func (w *Writer) Write(args ...interface{}) (bitsWritten int, err error) {
	for idx := range args {
		switch v := args[idx].(type) {
		case byte :
			if num, err := w.WriteByte(v); err != nil {
				break
			} else {
				bitsWritten += num
			}
		case []byte :
			if num, err := w.WriteBytes(v); err != nil {
				break
			} else {
				bitsWritten += num
			}
		case Bits:
			if num, err := w.WriteBits(v); err != nil {
				break
			} else {
				bitsWritten += num
			}
		case bool:
			if num, err := w.WriteBit(v); err != nil {
				break
			} else {
				bitsWritten += num
			}
		default:
			err = fmt.Errorf("bad type supplied for argument, index %v, value %v", idx, v)
		}
	}

	return bitsWritten, err
}

// WriteBytes writes the given bytes
func (w *Writer) WriteBytes(b []byte) (bitsWritten int, err error) {
	for idx := range b {
		numWritten, err := w.WriteByte(b[idx])

		bitsWritten += numWritten

		if err != nil {
			break
		}
	}

	return bitsWritten, err
}

// WriteByte writes the given byte
func (w *Writer) WriteByte(b byte) (bitsWritten int, err error) {
	bits := NewReader().FromBytes(b).Next(1).Bytes().Bits

	return w.WriteBits(bits)
}

// WriteBits writes the given Bits
func (w *Writer) WriteBits(b Bits) (bitsWritten int, err error) {
	for idx := range b {
		numWritten, err := w.WriteBit(b[idx])

		bitsWritten += numWritten

		if err != nil {
			break
		}
	}

	return bitsWritten, err
}

// WriteBit writes the given bit
func (w *Writer) WriteBit(b bool) (bitsWritten int, err error) {
	shift := uint8(0)
	bp := uint8(w.bitOffset)

	switch w.endianness {
	case LittleEndian:
		shift = bp
	case BigEndian:
		shift = uint8(bitsPerByte) - bp - 1
	}

	w.bitOffset++

	if b {
		w.bitBuffer |= 1<<shift
	}

	if w.bitOffset >= bitsPerByte {
		w.bitOffset = 0
		w.bytes = append(w.bytes, w.bitBuffer)
		w.bitBuffer = 0
	}

	return 1, err
}

// WriteBool writes the given Bits, an alias for WriteBit
func (w *Writer) WriteBool(b bool) (bitsWritten int, err error) {
	return w.WriteBit(b)
}
