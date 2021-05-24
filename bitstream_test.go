package bitstream

import (
	"io"
	"math/rand"
	"testing"
)

func TestBitStream_NilData(t *testing.T) {
	bs := New(nil)

	bit, err := bs.readBit()

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

	_ = bs.Read(1).Bits()

	if bs.BitPosition() != 1 {
		t.Error("expected bit position to be 1 after reading a bit")
	}

	_ = bs.Read(16).Bits()

	if bs.BitPosition() != 1 {
		t.Error("expected bit position to still be 1")
	}

	_ = bs.Read(17).Bits()

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

		got := bs.OffsetPosition(tests[idx].offset)

		if got != expected {
			t.Errorf("expected bit position to be %v, but got %v", expected, got)
		}
	}
}

func TestBitStream_Position(t *testing.T) {
	bytes := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	bs := FromBytes(bytes...)

	tests := []struct{
		set int
		expect int
	}{
		{-1, 0},
		{0, 0},
		{1, 1},
		{100, len(bytes)},
	}

	for _, test := range tests {
		if bs.SetPosition(test.set).Position() != test.expect {
			t.Errorf("set position to %v, expected %v but got %v", test.set, test.expect, bs.Position())
		}
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
		b, _ := bs.readBit()

		expected := tests[idx].expect

		if b != expected {
			t.Errorf("expected bit at position %v to be %v, got %v", idx, expected, b)
		}
	}

	bs.Options.ReadBeyondEOF = false

	if _, err := bs.readBit(); err != io.EOF {
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
		b := bs.Read(test.numBitsToRead).Bits()
		if b.Error != nil {
			t.Error(b.Error)
		}

		if len(b.Bits) != test.numBitsToRead {
			const fmtErr = "expected bits length of %v, got length %v"
			t.Errorf(fmtErr, test.numBitsToRead, len(b.Bits))
		}
	}
}

func TestBitStream_SetBitPosition(t *testing.T) {
	bs := FromBytes(0b1000_0000, 0b0000_0010)

	tests := []struct{
		bitPosition int
		expected bool
	}{
		{7, true}, // starts at 7, reads bit and ends at 0 in next byte
		{1, true}, // sets to 1, reads bit and ends at 2
		{-1, true}, // goes back one byte to position 7 of first byte, ends in second byte
		{-2, false}, // goes back to first byte, ends in first byte
		{9, true}, // goes to second bit of second byte
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

	tests := []struct{
		position int
		expected int
	}{
		{-1, 0},
		{0, 0},
		{2, 2},
		{7, 6},
	}

	for _, test := range tests {
		v := bs.SetPosition(test.position).Position()

		if v != test.expected {
			t.Errorf("expected bit at position %v to be %v, got %v", test.position, test.expected, v)
		}
	}
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

		res := bs.Read(1).Bytes()
		if res.Error != nil {
			t.Error(res.Error)
		}

		if res.AsByte() != test.expectedByte {
			const fmtErr = "expected bits as byte %v, but got %v"

			t.Errorf(fmtErr, test.expectedByte, res.AsByte())
		}
	}
}

func BenchmarkBitStream_ReadBits(b *testing.B) {
	bytes := make([]byte, 1024)
	rand.Read(bytes)

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

func readbits(b *testing.B, bs *BitStream, numBits int) {
	for i := 0; i < b.N; i++ {
		bs.SetPosition(0)
		bs.SetBitPosition(0)

		numBitsToRead := rand.Intn(numBits)
		_ = bs.Read(numBitsToRead).Bits()
	}
}