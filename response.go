package bitstream

// Response represents a response of Reader
type Response struct {
	Bits
	Error error
}

// AsBool interprets the bits as a bool
func (res Response) AsBool() (bool, error) {
	return res.Bits.AsBool(), res.Error
}

// AsByte interprets the bits as a byte
func (res Response) AsByte() (byte, error) {
	return res.Bits.AsByte(), res.Error
}

// AsBytes interprets the bits as a slice of bytes
func (res Response) AsBytes() ([]byte, error) {
	return res.Bits.AsBytes(), res.Error
}

// AsInt8 interprets the bits as a signed 8-bit integer
func (res Response) AsInt8() (int8, error) {
	return res.Bits.AsInt8(), res.Error
}

// AsUInt8 interprets the bits as a usnigned 8-bit integer
func (res Response) AsUInt8() (uint8, error) {
	return res.Bits.AsUInt8(), res.Error
}

// AsInt16 interprets the bits as a signed 16-bit integer
func (res Response) AsInt16() (int16, error) {
	return res.Bits.AsInt16(), res.Error
}

// AsUInt16 interprets the bits as an unsugned 16-bit integer
func (res Response) AsUInt16() (uint16, error) {
	return res.Bits.AsUInt16(), res.Error
}

// AsInt32 interprets the bits as a signed 32-bit integer
func (res Response) AsInt32() (int32, error) {
	return res.Bits.AsInt32(), res.Error
}

// AsUInt32 interprets the bits as an unsigned 32-bit integer
func (res Response) AsUInt32() (uint32, error) {
	return res.Bits.AsUInt32(), res.Error
}

// AsInt64 interprets the bits as a signed 64-bit integer
func (res Response) AsInt64() (int64, error) {
	return res.Bits.AsInt64(), res.Error
}

// AsUInt64 interprets the bits as an unsigned 64-bit integer
func (res Response) AsUInt64() (uint64, error) {
	return res.Bits.AsUInt64(), res.Error
}

// AsInt interprets the bits as a signed integer
func (res Response) AsInt() (int, error) {
	return res.Bits.AsInt(), res.Error
}

// AsUInt interprets the bits as an unsigned integer
func (res Response) AsUInt() (uint, error) {
	return res.Bits.AsUInt(), res.Error
}
