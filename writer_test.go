package bitstream

import (
	"reflect"
	"testing"
)

func TestWriter_Bytes(t *testing.T) {
	type fields struct {
		bytes       []byte
		currentByte byte
		bitOffset   int
	}
	tests := []struct {
		name string
		fields
		want []byte
	}{
		{"empty", fields{[]byte{}, 0, 0}, []byte{}},
		{"not empty", fields{[]byte{1, 2, 3, 4}, 1 << 6, 1}, []byte{1, 2, 3, 4, 1 << 6}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Writer{
				bytes:     tt.fields.bytes,
				bitOffset: tt.fields.bitOffset,
				bitBuffer: tt.fields.currentByte,
			}
			if got := w.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriter_WriteBit(t *testing.T) {
	tests := []struct {
		name            string
		bit             bool
		expectedWritten int
		expectedBytes   []byte
	}{
		{"false", F, 1, []byte{0}},
		{"true", T, 1, []byte{1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Writer{}

			gotBitsWritten, err := w.WriteBit(tt.bit)
			if err != nil {
				t.Errorf("WriteBit() error = %v", err)
				return
			}

			if gotBitsWritten != tt.expectedWritten {
				t.Errorf("WriteBit() gotBitsWritten = %v, want %v", gotBitsWritten, tt.expectedWritten)
			}

			if got := w.Bytes(); !reflect.DeepEqual(got, tt.expectedBytes) {
				t.Errorf("Bytes() = %v, want %v", got, tt.expectedBytes)
			}
		})
	}
}

func TestWriter_WriteBits(t *testing.T) {
	type state struct {
		bytes       []byte
		currentByte byte
		bitOffset   int
	}
	tests := []struct {
		name string
		state
		bitsToWrite   Bits
		expectedBytes []byte
	}{
		{
			"empty",
			state{[]byte{}, 0b110, 3},
			Bits{F, F, F, T},
			[]byte{0b_1000_110},
		},
		{
			"start unaligned, end aligned",
			state{[]byte{1, 2, 3, 4}, 0b11, 3},
			Bits{T, T, T, T},
			[]byte{1, 2, 3, 4, 0b_0111_1011},
		},
		{
			"start unaligned, end unaligned",
			state{[]byte{1, 2, 3, 4}, 0b10, 4},
			Bits{F, F, F, F, T},
			[]byte{1, 2, 3, 4, 0b_0000_0010, 0b1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Writer{
				bytes:     tt.state.bytes,
				bitBuffer: tt.state.currentByte,
				bitOffset: tt.state.bitOffset,
			}

			gotBitsWritten, err := w.WriteBits(tt.bitsToWrite)
			if err != nil {
				t.Errorf("WriteBits() error = %v", err)
				return
			}

			if gotBitsWritten != len(tt.bitsToWrite) {
				t.Errorf("WriteBits() gotBitsWritten = %v, want %v", gotBitsWritten, len(tt.bitsToWrite))
			}

			if got := w.Bytes(); !reflect.DeepEqual(got, tt.expectedBytes) {
				t.Errorf("Bytes() = %v, want %v", got, tt.expectedBytes)
			}
		})
	}
}

func TestWriter_WriteBool(t *testing.T) {
	tests := []struct {
		name            string
		bitsToWrite     []bool
		wantBitsWritten int
		expectedBytes   []byte
	}{
		{"write single false", []bool{F}, 1, []byte{0}},
		{"write single true", []bool{T}, 1, []byte{1}},
		{"write 8 true", []bool{T, T, T, T, T, T, T, T}, 8, []byte{255}},
		{"write 9 true", []bool{T, T, T, T, T, T, T, T, T}, 9, []byte{255, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Writer{}

			totalBitsWritten := 0

			for idx := range tt.bitsToWrite {
				bitsWritten, err := w.WriteBool(tt.bitsToWrite[idx])

				totalBitsWritten += bitsWritten

				if err != nil {
					t.Errorf("WriteBool() error = %v", err)
					return
				}
			}

			if totalBitsWritten != tt.wantBitsWritten {
				t.Errorf("WriteBool() gotBitsWritten = %v, want %v", totalBitsWritten, tt.wantBitsWritten)
			}

			if got := w.Bytes(); !reflect.DeepEqual(got, tt.expectedBytes) {
				t.Errorf("Bytes() = %v, want %v", got, tt.expectedBytes)
			}
		})
	}
}

func TestWriter_WriteByte(t *testing.T) {
	type state struct {
		existingBytes []byte
		bitBuffer     byte
		bitOffset     int
		endianness
	}
	tests := []struct {
		name string
		state
		bytesToWrite    []byte
		wantBitsWritten int
		expectedBytes   []byte
	}{
		{
			"from empty, write 0x1",
			state{
				existingBytes: nil,
			},
			[]byte{0x1},
			8,
			[]byte{0x1},
		},

		{
			"with existing bytes, aligned, multiple bytes",
			state{
				existingBytes: []byte{3, 2},
			},
			[]byte{2, 1},
			16,
			[]byte{3, 2, 2, 1},
		},

		{
			"with existing bytes, unaligned",
			state{
				existingBytes: []byte{6, 9, 4},
				bitOffset:     1,
			},
			[]byte{1},
			8,
			[]byte{6, 9, 4, 2, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Writer{
				bytes:      tt.state.existingBytes,
				bitOffset:  tt.state.bitOffset,
				bitBuffer:  tt.state.bitBuffer,
				endianness: tt.state.endianness,
			}

			totalBitsWritten := 0

			for idx := range tt.bytesToWrite {
				bitsWritten, err := w.WriteByte(tt.bytesToWrite[idx])

				totalBitsWritten += bitsWritten

				if err != nil {
					t.Errorf("WriteByte() error = %v", err)
					return
				}
			}

			if totalBitsWritten != tt.wantBitsWritten {
				t.Errorf("WriteByte() gotBitsWritten = %v, want %v", totalBitsWritten, tt.wantBitsWritten)
			}

			if got := w.Bytes(); !reflect.DeepEqual(got, tt.expectedBytes) {
				t.Errorf("WriteByte() = %v, want %v", got, tt.expectedBytes)
			}
		})
	}
}

func TestWriter_WriteBytes(t *testing.T) {
	type state struct {
		existingBytes []byte
		bitBuffer     byte
		bitOffset     int
		endianness
	}
	tests := []struct {
		name string
		state
		bytesToWrite    []byte
		wantBitsWritten int
		expectedBytes   []byte
	}{
		{
			"from empty, write 0x1",
			state{
				existingBytes: nil,
			},
			[]byte{0x1},
			8,
			[]byte{0x1},
		},

		{
			"with existing bytes, aligned, multiple bytes",
			state{
				existingBytes: []byte{3, 2},
			},
			[]byte{2, 1},
			16,
			[]byte{3, 2, 2, 1},
		},

		{
			"with existing bytes, unaligned",
			state{
				existingBytes: []byte{6, 9, 4},
				bitOffset:     1,
			},
			[]byte{1},
			8,
			[]byte{6, 9, 4, 2, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Writer{
				bytes:      tt.state.existingBytes,
				bitOffset:  tt.state.bitOffset,
				bitBuffer:  tt.state.bitBuffer,
				endianness: tt.state.endianness,
			}

			totalBitsWritten, err := w.WriteBytes(tt.bytesToWrite)
			if err != nil {
				t.Errorf("WriteBytes() error = %v", err)
				return
			}

			if totalBitsWritten != tt.wantBitsWritten {
				t.Errorf("WriteBytes() gotBitsWritten = %v, want %v", totalBitsWritten, tt.wantBitsWritten)
			}

			if got := w.Bytes(); !reflect.DeepEqual(got, tt.expectedBytes) {
				t.Errorf("WriteBytes() = %v, want %v", got, tt.expectedBytes)
			}
		})
	}
}

func TestWriter_Write(t *testing.T) {
	type fields struct {
		bytes      []byte
		bitBuffer  byte
		bitOffset  int
		endianness endianness
	}
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantBitsWritten int
		expectedBytes   []byte
		wantErr         bool
	}{
		{
			"from empty, combo write, end unaligned",
			fields{},
			args{args: []interface{}{F, F, byte(3)}},
			10,
			[]byte{0b_0000_1100, 0},
			false,
		},
		{
			"from non-empty, combo write, end unaligned",
			fields{},
			args{args: []interface{}{F, byte(3), []byte{0b0101_0101, 0b0011_0011, 0b0000_1111}}},
			33,
			[]byte{0b_0000_0110, 0b_1010_1010, 0b_0110_0110, 0b_0001_1110, 0},
			false,
		},
		{
			"bad arg string",
			fields{},
			args{args: []interface{}{"123"}},
			0,
			nil,
			true,
		},
		{
			"bad arg int64",
			fields{},
			args{args: []interface{}{int64(1234567)}},
			0,
			nil,
			true,
		},
		{
			"mixed with bad args",
			fields{},
			args{args: []interface{}{F, T, F, int64(1234567)}},
			3,
			[]byte{2},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Writer{
				bytes:      tt.fields.bytes,
				bitBuffer:  tt.fields.bitBuffer,
				bitOffset:  tt.fields.bitOffset,
				endianness: tt.fields.endianness,
			}
			gotBitsWritten, err := w.Write(tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotBitsWritten != tt.wantBitsWritten {
				t.Errorf("Write() gotBitsWritten = %v, want %v", gotBitsWritten, tt.wantBitsWritten)
			}

			if got := w.Bytes(); !reflect.DeepEqual(got, tt.expectedBytes) {
				bothEmpty := len(got) == 0 && len(tt.expectedBytes) == 0
				if !bothEmpty {
					t.Errorf("WriteBytes() = %v, want %v", got, tt.expectedBytes)
				}
			}
		})
	}
}

func TestWriter_WriteUint(t *testing.T) {
	tests := []struct {
		name            string
		uintsToWrite    []interface{}
		wantBitsWritten int
		expectedBytes   []byte
	}{
		{"write uint16", []interface{}{uint16(18), uint16(100), uint16(480)}, 48, []byte{18, 0, 100, 0, 224, 1}},
		{"write uint32", []interface{}{uint32(1024), uint32(40000), uint32(256)}, 96, []byte{0, 4, 0, 0, 64, 156, 0, 0, 0, 1, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Writer{}

			totalBitsWritten := 0

			for idx := range tt.uintsToWrite {
				bitsWritten, err := w.WriteUint(tt.uintsToWrite[idx])

				totalBitsWritten += bitsWritten

				if err != nil {
					t.Errorf("WriteBool() error = %v", err)
					return
				}
			}

			if totalBitsWritten != tt.wantBitsWritten {
				t.Errorf("WriteUint() gotBitsWritten = %v, want %v", totalBitsWritten, tt.wantBitsWritten)
			}

			if got := w.Bytes(); !reflect.DeepEqual(got, tt.expectedBytes) {
				t.Errorf("Bytes() = %v, want %v", got, tt.expectedBytes)
			}
		})
	}
}
