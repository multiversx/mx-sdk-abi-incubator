package abi

import (
	"encoding/hex"
	"strings"
)

type defaultDataReader struct {
	parts        [][]byte
	partIndex    int
	offsetInPart int
}

func NewDefaultDataReader(parts [][]byte) *defaultDataReader {
	return &defaultDataReader{
		parts: parts,
	}
}

func NewDefaultDataReaderFromString(encoded string) (*defaultDataReader, error) {
	partsHex := strings.Split(encoded, partsSeparator)
	parts := make([][]byte, len(partsHex))

	for i, partHex := range partsHex {
		part, err := hex.DecodeString(partHex)
		if err != nil {
			return nil, err
		}

		parts[i] = part
	}

	return &defaultDataReader{
		parts: parts,
	}, nil
}

func (d *defaultDataReader) Read(numBytes int) ([]byte, error) {
	part := d.parts[d.partIndex]

	if d.offsetInPart+numBytes > len(part) {
		return nil, errUnexpectedEndOfData
	}

	data := part[d.offsetInPart : d.offsetInPart+numBytes]
	d.offsetInPart += numBytes

	return data, nil
}

func (d *defaultDataReader) ReadWholePart() ([]byte, error) {
	if d.offsetInPart != 0 {
		return nil, errOffsetIsNotAtStartOfPart
	}

	part := d.parts[d.partIndex]
	d.offsetInPart = len(part)

	return part, nil
}

func (d *defaultDataReader) HasNextPart() bool {
	return d.partIndex < len(d.parts)-1
}

func (d *defaultDataReader) GotoNextPart() error {
	if d.offsetInPart != len(d.parts[d.partIndex]) {
		return errUnreadDataInPart
	}

	d.partIndex++
	d.offsetInPart = 0
	return nil
}
