package abi

import (
	"encoding/hex"
	"fmt"
	"strings"
)

type defaultDataReader struct {
	parts     [][]byte
	partIndex int
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

func (d *defaultDataReader) ReadWholePart() ([]byte, error) {
	if d.IsEndOfData() {
		return nil, fmt.Errorf("cannot wholly read part %d: unexpected end of data", d.partIndex)
	}

	part := d.parts[d.partIndex]
	return part, nil
}

func (d *defaultDataReader) GotoNextPart() error {
	if d.IsEndOfData() {
		return fmt.Errorf(
			"cannot advance to next part, since the reader is already beyond the last part; current part is %d",
			d.partIndex,
		)
	}

	d.partIndex++
	return nil
}

// IsEndOfData returns true if the reader is already beyond the last part.
func (d *defaultDataReader) IsEndOfData() bool {
	return d.partIndex >= len(d.parts)
}
