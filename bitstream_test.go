package bitstream

import (
	"io"
	"math/rand"
	"testing"
)

func TestBitStream_NilData(t *testing.T) {
	bs := New(nil)

	bit, err := bs.ReadBit()

	if err != io.EOF {
		t.Error("expecting end of file for nil read seeker")
	}

	if bit {
		t.Error("expecting false bit value read from nil data source")
	}
}

func TestBitStream_BitPosition(t *testing.T) {
	bs :=  FromBytes(
			0x01, 0x02, 0x03, 0x04,
			0x01, 0x02, 0x03, 0x04,
			0x01, 0x02, 0x03, 0x04,
			0x01, 0x02, 0x03, 0x04,
			0x01, 0x02, 0x03, 0x04,
			0x01, 0x02, 0x03, 0x04,
			0x01, 0x02, 0x03, 0x04,
			0x01, 0x02, 0x03, 0x04,
		)

	if bs.BitPosition() != 0 {
		t.Error("expected bit position to be 0 after init")
	}

	_, _ = bs.ReadBit()

	if bs.BitPosition() != 1 {
		t.Error("expected bit position to be 1 after reading a bit")
	}

	_ = bs.ReadBits(16)

	if bs.BitPosition() != 1 {
		t.Error("expected bit position to still be 1")
	}

	_ = bs.ReadBits(17)

	if bs.BitPosition() != 2 {
		t.Errorf("expected bit position to be 2, got %v", bs.bitPosition)
	}
}

func TestBitStream_OffsetBitPosition(t *testing.T) {
	bs :=  FromBytes(
		0x01, 0x02, 0x03, 0x04,
		0x01, 0x02, 0x03, 0x04,
		0x01, 0x02, 0x03, 0x04,
		0x01, 0x02, 0x03, 0x04,
		0x01, 0x02, 0x03, 0x04,
		0x01, 0x02, 0x03, 0x04,
		0x01, 0x02, 0x03, 0x04,
		0x01, 0x02, 0x03, 0x04,
	)

	tests := []struct{
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
	bs :=  FromBytes(
		0x01, 0x02, 0x03, 0x04,
		0x01, 0x02, 0x03, 0x04,
		0x01, 0x02, 0x03, 0x04,
		0x01, 0x02, 0x03, 0x04,
		0x01, 0x02, 0x03, 0x04,
		0x01, 0x02, 0x03, 0x04,
		0x01, 0x02, 0x03, 0x04,
		0x01, 0x02, 0x03, 0x04,
	)

	tests := []struct{
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

		bs.OffsetPosition(tests[idx].offset)

		if bs.bytePosition != expected {
			t.Errorf("expected bit position to be %v, but got %v", expected, bs.bytePosition)
		}
	}
}

func TestBitStream_Position(t *testing.T) {
	bs := New(nil)

	if bs.Position() != 0 {
		t.Error("expected byte position to be 0")
	}

	bs.SetPosition(-1)

	if bs.Position() != 0 {
		t.Error("expected byte position to be 0")
	}

	bs.bytePosition = 8

	if bs.Position() != 8 {
		t.Error("expected byte position to be 8")
	}

	bs.bytePosition = -3

	if bs.Position() != 0 {
		t.Error("expected byte position to be 0")
	}
}

func TestBitStream_ReadBit(t *testing.T) {
	bs := FromBytes(0b10000010)

	tests := []struct{
		expect bool
		err error
	}{
		{false, nil},
		{true, nil},
		{false, nil},
		{false, nil},
		{false, nil},
		{false, nil},
		{false, nil},
		{true, nil},
		{false, io.EOF},
	}

	bs.Options.ReadBeyondEOF = true

	for idx := range tests {
		b, _ := bs.ReadBit()

		expected := tests[idx].expect

		if b != expected {
			t.Errorf("expected bit at position %v to be %v, got %v", idx, expected, b)
		}
	}

	bs.Options.ReadBeyondEOF = false

	if _, err := bs.ReadBit(); err != io.EOF {
		t.Error("expected EOF")
	}
}

func TestBitStream_ReadBits(t *testing.T) {
	bs := FromBytes(0b10000010)

	tests := []struct{
		numBitsToRead int
	}{
		{4},
		{1},
		{10},
	}

	for _, test := range tests {
		b := bs.ReadBits(test.numBitsToRead)

		if len(b) != test.numBitsToRead {
			const fmtErr = "expected bits length of %v, got length %v"
			t.Errorf(fmtErr, test.numBitsToRead, len(b))
		}
	}
}

func TestBitStream_SetBitPosition(t *testing.T) {
	bs := FromBytes(0b1000_0000, 0b0000_0010)

	tests := []struct{
		bitPosition int
		expected bool
	}{
		{7, true},
		{1, true},
		{-1, true},
		{-2, false},
		{9, true},
	}

	for _, test := range tests {
		v, err := bs.SetBitPosition(test.bitPosition).ReadBit()
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
}


func TestBitStream_ReadByte_AsByte(t *testing.T) {
	bs := FromBytes(
		0b_1000_0000,
		0b_0000_0001,
		0b_0000_1111,
		0b_1100_1100,
	)

	tests := []struct{
		bitOffset int
		expectedByte byte
	}{
		{24, 0b_1100_1100}, // the 4th byte
		{22, 0b_0011_0000}, // two MSB's of 3rd byte, six LSB's of 4th
		{23, 0b_1001_1000}, // one MSB of 3rd byte, seven LSB's of 4th
		{4, 0b_0001_1000}, // 4 MSB's of 1st byte, 4 LSB's of 2nd
	}

	for _, test := range tests {
		// set byte posiiton to 0
		bs.SetPosition(0)

		// offset n bits from 0.
		// going past 7 offsets the current byte
		bs.SetBitPosition(test.bitOffset)

		val := bs.ReadBits(bitsPerByte).AsByte()

		if val != test.expectedByte {
			const fmtErr = "expected bits as byte %v, but got %v"

			t.Errorf(fmtErr, test.expectedByte, val)
		}
	}
}

func BenchmarkBitStream_ReadBits(b *testing.B) {
	bytes := make([]byte, 1024)
	rand.Read(bytes)

	bs := FromBytes(bytes...)

	b.Run("readbit", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = bs.SetPosition(0).SetBitPosition(0).ReadBit()
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

func readbits(b *testing.B, bs *BitStream, numBits int) {
	for i := 0; i < b.N; i++ {
		bs.SetPosition(0)
		bs.SetBitPosition(0)

		numBitsToRead := rand.Intn(numBits)
		_ = bs.ReadBits(numBitsToRead)
	}
}

//func TestBits_AsInt(t *testing.T) {
//
//}
//
//func TestBits_AsInt16(t *testing.T) {
//
//}
//
//func TestBits_AsUInt16(t *testing.T) {
//
//}