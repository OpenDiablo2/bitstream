package bitstream

import (
	"reflect"
	"testing"
)

const (
	T = true
	F = false
)

func TestBits_AsByte(t *testing.T) {
	tests := []struct {
		name string
		b    Bits
		want byte
	}{
		{"empty", Bits{}, 0},
		{"4th is T", Bits{F, F, F, T}, 8},
		{"8th bit T", Bits{F, F, F, F, F, F, F, F, T}, 0},
		{"9th bit T", Bits{F, F, F, F, F, F, F, T}, 128},
		{"8 bits T", Bits{T, T, T, T, T, T, T, T}, 255},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.AsByte(); got != tt.want {
				t.Errorf("AsByte() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBits_AsInt(t *testing.T) {
	tests := []struct {
		name string
		b    Bits
		want int
	}{
		{"empty", Bits{}, 0},

		{"negative 1 (4-bit)", Bits{
			// LSB is ON THE LEFT!!!!
			// LSB is ON THE LEFT!!!!
			// LSB is ON THE LEFT!!!!
			true, true, true, true,
		}, -1},

		{"negative 1 (8 bit)", Bits{true, true, true, true, true, true, true, true}, -1},
		{"negative 7 (8 bit)", Bits{true, false, false, true, true, true, true, true}, -7},
		{"positive 7 (4 bit)", Bits{true, true, true, false}, 7},
		{"negative 1 (1 bit)", Bits{true}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.AsInt(); got != tt.want {
				t.Errorf("AsInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBits_AsInt16(t *testing.T) {
	tests := []struct {
		name string
		b    Bits
		want int16
	}{
		{"empty", Bits{}, 0},
		{"negative 1 (1 bit)", Bits{T}, -1},
		{"negative 1 (3 bit)", Bits{T, T, T}, -1},
		{"negative 1 (16 bit)", Bits{T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T}, -1},
		{"negative 1 (17 bit)", Bits{T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T, T}, -1},
		{"0 (17 bit)", Bits{F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, T}, 0},
		{"positive 2 (3 bit)", Bits{F, T, F}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.AsInt16(); got != tt.want {
				t.Errorf("AsInt16() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestBits_AsInt32(t *testing.T) {
//	tests := []struct {
//		name string
//		b    Bits
//		want int32
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := tt.b.AsInt32(); got != tt.want {
//				t.Errorf("AsInt32() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//func TestBits_AsInt64(t *testing.T) {
//	tests := []struct {
//		name string
//		b    Bits
//		want int64
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := tt.b.AsInt64(); got != tt.want {
//				t.Errorf("AsInt64() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//func TestBits_AsInt8(t *testing.T) {
//	tests := []struct {
//		name string
//		b    Bits
//		want int8
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := tt.b.AsInt8(); got != tt.want {
//				t.Errorf("AsInt8() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//func TestBits_AsUInt(t *testing.T) {
//	tests := []struct {
//		name string
//		b    Bits
//		want uint
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := tt.b.AsUInt(); got != tt.want {
//				t.Errorf("AsUInt() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//func TestBits_AsUInt16(t *testing.T) {
//	tests := []struct {
//		name string
//		b    Bits
//		want uint16
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := tt.b.AsUInt16(); got != tt.want {
//				t.Errorf("AsUInt16() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//func TestBits_AsUInt32(t *testing.T) {
//	tests := []struct {
//		name string
//		b    Bits
//		want uint32
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := tt.b.AsUInt32(); got != tt.want {
//				t.Errorf("AsUInt32() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//func TestBits_AsUInt64(t *testing.T) {
//	tests := []struct {
//		name string
//		b    Bits
//		want uint64
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := tt.b.AsUInt64(); got != tt.want {
//				t.Errorf("AsUInt64() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//func TestBits_AsUInt8(t *testing.T) {
//	tests := []struct {
//		name string
//		b    Bits
//		want uint8
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := tt.b.AsUInt8(); got != tt.want {
//				t.Errorf("AsUInt8() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func TestBits_AsBytes(t *testing.T) {
	tests := []struct {
		name string
		b    Bits
		want []byte
	}{
		{"", Bits{F, T}, []byte{2}},
		{"", Bits{F, F, F, F, F, F, F, F, T}, []byte{0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.AsBytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AsBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}