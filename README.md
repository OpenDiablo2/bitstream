# What is a bitstream?
It is a utility for reading and interpreting data that is not always byte-aligned.

## How can I use it?
Suppose you have a binary file with the following format:
* First 6 bits is an unsigned integer for the major version
* following 7 bits is padded 0's
* there will be eight 4-byte strings, each with 4 bits of zero padding 

The file could be parsed in this way:
```golang
stream := bitstream.FromBytes(fileBytes)

var err error

version, err := stream.Next(6).Bits().AsUInt()
if err != nil {
    // handle it
}

// check the zero pad
pad, err := stream.Next(7).Bits().AsInt()
if pad > 0 || err != nil {
    // handle it
}

strings := make([]string, 8)

for strIdx := range strings {
    chars := make([]byte, 4)
    
    // read 4 characters
    for charIdx := range chars {
        chars[charIdx], err = stream.Next(1).Bytes().AsByte()
        if err != nil {
            // handle it
        }
    }

    // check the zero pad
    pad, err := stream.Next(4).Bits().AsInt()
    if pad > 0 || err != nil {
        // handle it
    }

    // cast the 4 characters as a string
    strings[strIdx] = (string)(chars)
}
```
