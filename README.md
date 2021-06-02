# What is a bitstream?

It is a utility for reading, writing, and interpreting data that is not always byte-aligned. This module provides a
Reader and Writer bitstream implementation.

## How can I use a bitstream.Reader?

Suppose you have a binary file with the following format:

* First 6 bits is an unsigned integer for the major version
* following 7 bits is padded 0's
* there will be three 4-byte strings, each with 4 bits of zero padding

The file could be parsed in this way:

```golang
package main

import (
	"github.com/gravestench/bitstream"
)

func main() {
	r := bitstream.NewReader().FromBytes(fileBytes)

	var err error

	version, err := r.Next(6).Bits().AsUInt()
	if err != nil {
		// handle it
	}

	// check the zero pad
	pad, err := r.Next(7).Bits().AsInt()
	if pad > 0 || err != nil {
		// handle it
	}

	strings := make([]string, 3)

	for strIdx := range strings {
		chars := make([]byte, 4)

		// read 4 characters
		for charIdx := range chars {
			chars[charIdx], err = r.Next(1).Bytes().AsByte()
			if err != nil {
				// handle it
			}
		}

		// check the zero pad
		pad, err := r.Next(4).Bits().AsInt()
		if pad > 0 || err != nil {
			// handle it
		}

		// cast the 4 characters as a string
		strings[strIdx] = (string)(chars)
	}
}
```

## How can I use a bitstream.Writer?

Assuming the same file format as above:

```golang
package main

import (
	"github.com/gravestench/bitstream"
)

type example struct {
	version uint8
	strings [4]string
}

func main() {
	e := example{
		version: 5,
		strings: [4]string{
			"abcd",
			"efgh",
			"ijkl",
		},
	}

	w := &bitstream.Writer{}
	bits := bitstream.BitsFromByte(e.version)[:6]

	if bitsWritten, err := w.Write(bits); err != nil {
		// handle it
	} else if bitsWritten != 6 {
		// handle it
	}

	pad7 := bitstream.BitsFromByte(0)[:7]

	if bitsWritten, err := w.Write(pad7); err != nil {
		// handle it
	} else if bitsWritten != 7 {
		// handle it
	}

	pad4 := bitstream.BitsFromByte(0)[:4]

	for idx := range e.strings {
		bytes := ([]byte)(e.strings[idx])[:4]
		w.Write(bytes, pad4)
	}
}
```