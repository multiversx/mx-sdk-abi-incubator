package abi

import (
	"encoding/hex"
	"fmt"
	"strings"
)

type dataReader struct {
	parts     [][]byte
	partIndex int
}

func newDataReader(parts [][]byte) *dataReader {
	return &dataReader{
		parts: parts,
	}
}

func newDataReaderFromString(encoded string) (*dataReader, error) {
	partsHex := strings.Split(encoded, partsSeparator)
	parts := make([][]byte, len(partsHex))

	for i, partHex := range partsHex {
		part, err := hex.DecodeString(partHex)
		if err != nil {
			return nil, err
		}

		parts[i] = part
	}

	return &dataReader{
		parts: parts,
	}, nil
}

func (reader *dataReader) ReadWholePart() ([]byte, error) {
	if reader.IsEndOfData() {
		return nil, fmt.Errorf("cannot wholly read part %d: unexpected end of data", reader.partIndex)
	}

	part := reader.parts[reader.partIndex]
	return part, nil
}

func (reader *dataReader) GotoNextPart() error {
	if reader.IsEndOfData() {
		return fmt.Errorf(
			"cannot advance to next part, since the reader is already beyond the last part; current part is %d",
			reader.partIndex,
		)
	}

	reader.partIndex++
	return nil
}

// IsEndOfData returns true if the reader is already beyond the last part.
func (reader *dataReader) IsEndOfData() bool {
	return reader.partIndex >= len(reader.parts)
}
