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
		return nil, newErrReaderCannotReadDueToEndOfData(numBytes, d.offsetInPart, d.partIndex)
	}

	data := part[d.offsetInPart : d.offsetInPart+numBytes]
	d.offsetInPart += numBytes

	return data, nil
}

func (d *defaultDataReader) ReadWholePart() ([]byte, error) {
	if d.offsetInPart != 0 {
		return nil, newErrReaderCannotReadWholePartDueToNonZeroOffset(d.partIndex, d.offsetInPart)
	}
	if d.partIndex >= len(d.parts) {
		return nil, newErrReaderCannotReadWholePartDueToEndOfData(d.partIndex)
	}

	part := d.parts[d.partIndex]
	d.offsetInPart = len(part)

	return part, nil
}

func (d *defaultDataReader) HasUnreadData() bool {
	if len(d.parts) == 0 {
		return false
	}
	if d.partIndex >= len(d.parts) {
		return false
	}

	isLastPart := d.partIndex == len(d.parts)-1
	isPartRead := d.offsetInPart == len(d.parts[d.partIndex])
	if isLastPart && isPartRead {
		return false
	}

	return true
}

func (d *defaultDataReader) GotoNextPart() error {
	lengthOfCurrentPart := len(d.parts[d.partIndex])
	numUnreadBytes := lengthOfCurrentPart - d.offsetInPart

	if numUnreadBytes > 0 {
		return newErrReaderCannotGotoNextPartDueToUnreadDataInCurrentPart(numUnreadBytes, d.partIndex, lengthOfCurrentPart)
	}

	d.partIndex++
	d.offsetInPart = 0
	return nil
}
