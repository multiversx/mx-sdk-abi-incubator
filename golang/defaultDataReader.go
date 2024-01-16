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
		return nil, newErrReaderCannotReadDueToEndOfDataInPart(numBytes, d.offsetInPart, d.partIndex)
	}

	data := part[d.offsetInPart : d.offsetInPart+numBytes]
	d.offsetInPart += numBytes

	return data, nil
}

func (d *defaultDataReader) ReadWholePart() ([]byte, error) {
	if d.offsetInPart != 0 {
		return nil, newErrReaderCannotReadWholePartDueToNonZeroOffset(d.partIndex, d.offsetInPart)
	}
	if d.IsEndOfData() {
		return nil, newErrReaderCannotReadWholePartDueToEndOfData(d.partIndex)
	}

	part := d.parts[d.partIndex]
	d.offsetInPart = len(part)

	return part, nil
}

func (d *defaultDataReader) GotoNextPart() error {
	lengthOfCurrentPart := len(d.parts[d.partIndex])
	numUnreadBytes := lengthOfCurrentPart - d.offsetInPart

	if numUnreadBytes > 0 {
		return newErrReaderCannotGotoNextPartDueToUnreadDataInCurrentPart(numUnreadBytes, d.partIndex, lengthOfCurrentPart)
	}
	if d.IsEndOfData() {
		return newErrReaderCannotGotoNextPartDueToEndOfData(d.partIndex)
	}

	d.partIndex++
	d.offsetInPart = 0
	return nil
}

func (d *defaultDataReader) IsCurrentPartEmpty() bool {
	return len(d.parts) > 0 && len(d.parts[d.partIndex]) == 0
}

// IsEndOfData returns true if the reader is already beyond the last part.
func (d *defaultDataReader) IsEndOfData() bool {
	return d.partIndex >= len(d.parts)
}
