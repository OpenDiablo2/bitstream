package bitstream

import (
	"crypto/rand"
	"errors"
	"io"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitStream_NilData(t *testing.T) {
	bs := New(nil)

	bit, err := bs.readBit()

	if !errors.Is(err, io.EOF) {
		t.Error("expecting end of file for nil read seeker")
	}

	if bit {
		t.Error("expecting false bit value read from nil data source")
	}
}

func TestBitStream_BitPosition(t *testing.T) {
	bs := FromBytes(
		0x01, 0x02, 0x03, 0x04,
		0x01, 0x02, 0x03, 0x04,
	)

	tests := []struct {
		bitsToRead           int
		expectedBytePosition int
		expectedError        error
	}{
		{3, 0, nil},
		{3, 0, nil},
		{2, 1, nil},
		{8, 2, nil},
		{8, 3, nil},
		{64, 8, io.EOF},
	}

	for _, test := range tests {
		res := bs.Next(test.bitsToRead).Bits()
		if !errors.Is(res.Error, test.expectedError) {
			t.Error(res.Error)
		}

		bytePosition := bs.Position()

		if bytePosition != test.expectedBytePosition {
			const fmtMsg = "read %v bits, expecting byte position %v, but got %v"

			t.Errorf(fmtMsg, test.bitsToRead, test.expectedBytePosition, bytePosition)
		}
	}
}

func TestBitStream_OffsetBitPosition(t *testing.T) {
	bs := FromBytes(
		0x01, 0x02, 0x03, 0x04,
		0x01, 0x02, 0x03, 0x04,
		0x01, 0x02, 0x03, 0x04,
		0x01, 0x02, 0x03, 0x04,
		0x01, 0x02, 0x03, 0x04,
		0x01, 0x02, 0x03, 0x04,
		0x01, 0x02, 0x03, 0x04,
		0x01, 0x02, 0x03, 0x04,
	)

	tests := []struct {
		offset int
		expect int
	}{
		{3, 3},
		{-3, 0},
		{-3, 0},
		{9, 1},
		{-2, 7},
	}

	for idx := range tests {
		expected := tests[idx].expect

		bs.OffsetBitPosition(tests[idx].offset)

		if bs.bitPosition != expected {
			t.Errorf("expected bit position to be %v, but got %v", expected, bs.bitPosition)
		}
	}
}

func TestBitStream_OffsetPosition(t *testing.T) {
	bs := FromBytes(
		0x01, 0x02, 0x03, 0x04,
		0x01, 0x02, 0x03, 0x04,
		0x01, 0x02, 0x03, 0x04,
		0x01, 0x02, 0x03, 0x04,
		0x01, 0x02, 0x03, 0x04,
		0x01, 0x02, 0x03, 0x04,
		0x01, 0x02, 0x03, 0x04,
		0x01, 0x02, 0x03, 0x04,
	)

	tests := []struct {
		offset int
		expect int
	}{
		{3, 3},
		{-3, 0},
		{-3, 0},
		{9, 9},
		{2, 11},
	}

	for idx := range tests {
		expected := tests[idx].expect

		got := bs.OffsetPosition(tests[idx].offset)

		if got != expected {
			t.Errorf("expected bit position to be %v, but got %v", expected, got)
		}
	}
}

func TestBitStream_Position(t *testing.T) {
	bytes := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	bs := FromBytes(bytes...)

	tests := []struct {
		set    int
		expect int
	}{
		{-1, 0},
		{0, 0},
		{1, 1},
		{100, 100},
	}

	for _, test := range tests {
		if bs.SetPosition(test.set).Position() != test.expect {
			t.Errorf("set position to %v, expected %v but got %v", test.set, test.expect, bs.Position())
		}
	}
}

func TestBitStream_ReadBit(t *testing.T) {
	bs := FromBytes(130)

	tests := []struct {
		expect bool
		err    error
	}{
		{false, nil},
		{true, nil},
		{false, nil},
		{false, nil},
		{false, nil},
		{false, nil},
		{false, nil},
		{true, nil},
		{false, nil},
	}

	for idx := range tests {
		b, _ := bs.readBit()

		expected := tests[idx].expect

		if b != expected {
			t.Errorf("expected bit at position %v to be %v, got %v", idx, expected, b)
		}
	}

	if _, err := bs.readBit(); !errors.Is(err, io.EOF) {
		t.Errorf("expected EOF error, got %v", err)
	}
}

func TestBitStream_ReadBits(t *testing.T) {
	bs := FromBytes(130)

	tests := []struct {
		numBitsToRead int
		expectedError error
	}{
		{4, nil},
		{1, nil},
		{10, io.EOF},
	}

	for _, test := range tests {
		b := bs.Next(test.numBitsToRead).Bits()
		if !errors.Is(b.Error, test.expectedError) {
			t.Error(b.Error)
		}

		if len(b.Bits) != test.numBitsToRead {
			const fmtErr = "expected bits length of %v, got length %v"

			t.Errorf(fmtErr, test.numBitsToRead, len(b.Bits))
		}
	}
}

func TestBitStream_SetBitPosition(t *testing.T) {
	bs := FromBytes(128, 2)

	tests := []struct {
		bitPosition int
		expected    bool
	}{
		{7, true},   // starts at 7, reads bit and ends at 0 in next byte
		{1, true},   // sets to 1, reads bit and ends at 2
		{-1, true},  // goes back one byte to position 7 of first byte, ends in second byte
		{-2, false}, // goes back to first byte, ends in first byte
		{9, true},   // goes to second bit of second byte
	}

	for _, test := range tests {
		v, err := bs.SetBitPosition(test.bitPosition).readBit()
		if err != nil {
			t.Error(err)
			continue
		}

		if v != test.expected {
			t.Errorf("expected bit at position %v to be %v, got %v", test.bitPosition, test.expected, v)
		}
	}
}

func TestBitStream_SetPosition(t *testing.T) {
	bs := FromBytes(0, 0, 0, 0, 0, 0)

	tests := []struct {
		position int
		expected int
	}{
		{-1, 0},
		{0, 0},
		{2, 2},
		{7, 7},
	}

	for _, test := range tests {
		v := bs.SetPosition(test.position).Position()

		if v != test.expected {
			const fmtErr = "expected to be at position %v after offsetting by %v, got %v"

			t.Errorf(fmtErr, test.expected, test.position, v)
		}
	}
}

func TestBitStream_ReadByte_AsByte(t *testing.T) {
	bs := FromBytes(
		128, 1, 15, 204,
	)

	tests := []struct {
		bitOffset    int
		expectedByte byte
	}{
		{24, 0b_1100_1100}, // the 4th byte
		{22, 0b_0011_0000}, // two MSB's of 3rd byte, six LSB's of 4th
		{23, 0b_1001_1000}, // one MSB of 3rd byte, seven LSB's of 4th
		{4, 0b_0001_1000},  // 4 MSB's of 1st byte, 4 LSB's of 2nd
	}

	for _, test := range tests {
		// set byte posiiton to 0
		bs.SetPosition(0)

		// offset n bits from 0.
		// going past 7 offsets the current byte
		bs.SetBitPosition(test.bitOffset)

		val, err := bs.Next(1).Bytes().AsByte()
		if err != nil {
			t.Error(err)
		}

		if val != test.expectedByte {
			const fmtErr = "expected bits as byte %v, but got %v"

			t.Errorf(fmtErr, test.expectedByte, val)
		}
	}
}

func BenchmarkBitStream_ReadBits(b *testing.B) {
	bytes := make([]byte, 1024)

	if _, err := rand.Read(bytes); err != nil {
		b.Error(err)
	}

	bs := FromBytes(bytes...)

	b.Run("readbit", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = bs.SetPosition(0).SetBitPosition(0).readBit()
		}
	})

	b.Run("1bit reads", func(b *testing.B) {
		readbits(b, bs, 1)
	})

	b.Run("8bit reads", func(b *testing.B) {
		readbits(b, bs, 8)
	})

	b.Run("16bit reads", func(b *testing.B) {
		readbits(b, bs, 16)
	})

	b.Run("32bit reads", func(b *testing.B) {
		readbits(b, bs, 32)
	})

	b.Run("64bit reads", func(b *testing.B) {
		readbits(b, bs, 32)
	})

	bs = New(nil)

	b.Run("8bit reads (from nil stream)", func(b *testing.B) {
		readbits(b, bs, 8)
	})

	b.Run("16bit reads (from nil stream)", func(b *testing.B) {
		readbits(b, bs, 16)
	})

	b.Run("32bit reads (from nil stream)", func(b *testing.B) {
		readbits(b, bs, 32)
	})

	b.Run("64bit reads (from nil stream)", func(b *testing.B) {
		readbits(b, bs, 32)
	})
}

func readbits(b *testing.B, bs *Reader, numBits int) {
	for i := 0; i < b.N; i++ {
		bs.SetPosition(0)
		bs.SetBitPosition(0)

		numBitsToRead, err := rand.Int(rand.Reader, big.NewInt(int64(numBits)))
		if err != nil {
			b.Error(err)
		}

		_ = bs.Next(int(numBitsToRead.Int64())).Bits()
	}
}

func TestFromBytes(t *testing.T) {
	s := FromBytes(0, 32)

	s.OffsetBitPosition(12) // second byte, will read from 5th LSB next

	val, err := s.Next(2).Bits().AsByte()
	if err != nil {
		t.Error("Reading bits returned End Of File")
	}

	assert.Equal(t, byte(2), val, "unexpected value returned")
}

func TestBitstream_Copy(t *testing.T) {
	original := FromBytes(
		0b_0000_0000,
		0b_0000_0001,
		0b_0000_0010,
		0b_1111_0011,
		0b_0000_0100,
		)

	_ = original.Next(8).Bits() // skip

	position, bitPosition := original.Position(), original.BitPosition()

	s := original.Copy()
	b1, err := s.Next(1).Bytes().AsUInt()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, original.Position(), position, "position of original stream after copy is changed")
	assert.Equal(t, original.BitPosition(), bitPosition, "bit position of original stream after copy is changed")

	assert.Equal(t, uint(1), b1, "unexpected value returned")
	assert.Equal(t, 8, s.BitsRead(), "unexpected value returned")

	s = s.Copy()
	b2, err := s.Next(12).Bits().AsUInt()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, uint(0b_0011_0000_0010), b2, "unexpected value returned")
	assert.Equal(t, 12, s.BitsRead(), "unexpected value returned")

	s = s.Copy()
	b3, err := s.Next(1).Bytes().AsUInt()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, uint(0b_0100_1111), b3, "unexpected value returned")
	assert.Equal(t, 8, s.BitsRead(), "unexpected value returned")

	s = s.Copy()
	b4, err := s.Next(1).Bytes().AsUInt()
	if err == nil {
		t.Error("expecting End Of File")
	}

	assert.Equal(t, uint(0), b4, "unexpected value returned")
	assert.Equal(t, 4, s.BitsRead(), "unexpected value returned")

	s = s.Copy()
	s.OffsetBitPosition(-12)
	b5, err := s.Next(7).Bits().AsUInt()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, uint(0b_0100_1111), b5, "unexpected value returned")
	assert.Equal(t, 7, s.BitsRead(), "unexpected value returned")
}

// Reader is supposed to replace Bitmuncher, this was a test to ensure they worked the same.
// func TestAgainstBitmuncher(t *testing.T) {
//	rand.Seed(time.Now().UnixNano())
//
//	b := make([]byte, rand.Intn(1024*1024))
//
//	for idx := range b {
//		b[idx] = byte(rand.Int())
//	}
//
//	bm := d2datautils.CreateBitMuncher(b, 0)
//	bs := FromBytes(b...)
//
//	funcMap := []func() error {
//		func() error {
//			b1 := bm.GetBit()
//			b2, _ := bs.Next(1).Bits().AsUInt32()
//
//			if b1 != b2 {
//				return errors.New("bits dont match")
//			}
//
//			return nil
//		},
//		func() error {
//			b1 := bm.GetByte()
//			b2, _ := bs.Next(1).Bytes().AsByte()
//
//			if bm.Offset() + 1 >= len(b) {
//				return nil
//			}
//
//			if b1 != b2 {
//				return fmt.Errorf("bytes dont match: %v != %v", b1, b2)
//			}
//
//			return nil
//		},
//		func() error {
//			if bm.Offset() + 3 >= len(b) {
//				return nil
//			}
//
//			b1 := bm.GetBits(12)
//			b2, _ := bs.Next(12).Bits().AsUInt32()
//
//			if b1 != b2 {
//				return fmt.Errorf("12 bits as uint32 dont match: %v != %v", b1, b2)
//			}
//
//			return nil
//		},
//		func() error {
//			bm = bm.Copy()
//			bs = bs.Copy()
//
//			if bm.BitsRead() != bs.BitsRead() {
//				return fmt.Errorf("tally of bits read not equal: %v != %v", bm.BitsRead(), bs.BitsRead())
//			}
//
//			absoluteBitOffset := (bs.Position() * 8) + bs.BitPosition()
//			if bm.Offset() != absoluteBitOffset {
//				return fmt.Errorf("absolute bit position not equal: %v != %v", bm.Offset(), absoluteBitOffset)
//			}
//
//			return nil
//		},
//		func() error {
//			return nil
//		},
//		func() error {
//			return nil
//		},
//	}
//
//	numBits := len(b) * 8
//	for numBits > 8 {
//		randIdx := rand.Intn(len(funcMap))
//		if err := funcMap[randIdx](); err != nil {
//			t.Error(err)
//			break
//		}
//
//		if bm.Offset() >= len(b) - 1 {
//			break
//		}
//
//		numBits--
//	}
// }
