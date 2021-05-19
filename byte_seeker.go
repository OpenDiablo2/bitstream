package bitstream

import (
	"errors"
	"io"
)

type byteSeeker struct {
	bytes  []byte
	offset int
}

func (d *byteSeeker) Read(result []byte) (n int, err error) {
	numRead := 0

	for idx := 0; idx < len(result); idx++ {
		if (idx + d.offset) >= len(d.bytes) {
			break
		}

		result[idx] = d.bytes[idx+d.offset]
		numRead++
	}

	if numRead < len(result) {
		err = io.EOF
	}

	return numRead, err
}

func (d *byteSeeker) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		d.offset = int(offset)
	case io.SeekEnd:
		d.offset = len(d.bytes) - int(offset)
	case io.SeekCurrent:
		d.offset += int(offset)
	}

	var err error

	if d.offset > len(d.bytes) {
		const errStr = "offset greater than max possible offset"

		d.offset = len(d.bytes)
		err = errors.New(errStr)
	}

	if d.offset < 0 {
		const errStr = "offset was negative"

		d.offset = 0
		err = errors.New(errStr)
	}

	return int64(d.offset), err
}
