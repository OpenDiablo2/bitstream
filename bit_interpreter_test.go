//nolint:gocritic // (appear for 'commentFormating) TODO
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
			T, T, T, T,
		}, -1},

		{"negative 1 (8 bit)", Bits{T, T, T, T, T, T, T, T}, -1},
		{"negative 7 (8 bit)", Bits{T, F, F, T, T, T, T, T}, -7},
		{"positive 7 (4 bit)", Bits{T, T, T, F}, 7},
		{"negative 1 (1 bit)", Bits{T}, -1},
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

func TestBits_AsInt32(t *testing.T) {
	tests := []struct {
		name string
		b    Bits
		want int32
	}{
		{"empty", Bits{}, 0},
		{"negative 310 (10 bits)", Bits{F, T, F, T, F, F, T, T, F, T}, -310},
		{"positive 1024 (12 bits)", Bits{F, F, F, F, F, F, F, F, F, F, T, F}, 1024},
		{"negative 1024 (12 bits)", Bits{F, F, F, F, F, F, F, F, F, F, T, T}, -1024},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.AsInt32(); got != tt.want {
				t.Errorf("AsInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func TestBits_AsInt8(t *testing.T) {
	tests := []struct {
		name string
		b    Bits
		want int8
	}{
		{"positive 100 (8 bits)", Bits{F, F, T, F, F, T, T, F}, 100},
		{"negative 100 (8 bits)", Bits{F, F, T, T, T, F, F, T}, -100},
		{"positive 8 (5 bits)", Bits{F, F, F, T, F}, 8},
		{"negative 8 (5 bits)", Bits{F, F, F, T, T}, -8},
		{"positive 64 (8 bits)", Bits{F, F, F, F, F, F, T, F}, 64},
		{"negative 64 (7 bits)", Bits{F, F, F, F, F, F, T}, -64},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.AsInt8(); got != tt.want {
				t.Errorf("AsInt8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBits_AsUInt(t *testing.T) {
	tests := []struct {
		name string
		b    Bits
		want uint
	}{
		{"2 (2 bits)", Bits{F, T}, 2},
		{"70 (7 bits)", Bits{F, T, T, F, F, F, T}, 70},
		{"123 (7 bits)", Bits{T, T, F, T, T, T, T}, 123},
		{"256 (8 bits)", Bits{T, T, T, T, T, T, T, T}, 255},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.AsUInt(); got != tt.want {
				t.Errorf("AsUInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBits_AsUInt16(t *testing.T) {
	tests := []struct {
		name string
		b    Bits
		want uint16
	}{
		{"1945 (11 bits)", Bits{T, F, F, T, T, F, F, T, T, T, T}, 1945},
		{"1945 (16 bits)", Bits{T, F, F, T, T, F, F, T, T, T, T, F, F, F, F, F}, 1945},
		{"2021 (11 bits)", Bits{T, F, T, F, F, T, T, T, T, T, T}, 2021},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.AsUInt16(); got != tt.want {
				t.Errorf("AsUInt16() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBits_AsUInt32(t *testing.T) {
	tests := []struct {
		name string
		b    Bits
		want uint32
	}{
		{"2^31 (5 bits)", Bits{T, F, T, T, T}, 2 ^ 31},
		{"2^31 (20 bits)", Bits{T, F, T, T, T, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F}, 2 ^ 31},
		{"2^31 (32 bits)", Bits{T, F, T, T, T, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F, F}, 2 ^ 31},
		{"1548876 (uint32)", Bits{F, F, T, T, F, F, T, F, F, T, F, F, F, T, F, T, T, T, T, F, T, F, F, F, F, F, F, F, F, F, F, F}, 1548876},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.AsUInt32(); got != tt.want {
				t.Errorf("AsUInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
