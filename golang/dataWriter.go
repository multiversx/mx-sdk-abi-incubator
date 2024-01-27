package abi

import (
	"encoding/hex"
	"errors"
	"strings"
)

type dataWriter struct {
	parts [][]byte
}

// newDataWriter creates a new dataWriter.
// A newly-created writer has no parts.
// Parts are created by calling GotoNextPart().
func newDataWriter() *dataWriter {
	return &dataWriter{
		parts: [][]byte{},
	}
}

func (writer *dataWriter) Write(data []byte) (int, error) {
	if len(writer.parts) == 0 {
		return 0, errors.New("cannot write, since there is no part to write to")
	}

	partIndex := len(writer.parts) - 1
	writer.parts[partIndex] = append(writer.parts[partIndex], data...)
	return len(data), nil
}

func (writer *dataWriter) GotoNextPart() {
	writer.parts = append(writer.parts, []byte{})
}

func (writer *dataWriter) GetParts() [][]byte {
	return writer.parts
}

func (writer *dataWriter) String() string {
	partsHex := make([]string, len(writer.parts))

	for i, part := range writer.parts {
		partsHex[i] = hex.EncodeToString(part)
	}

	return strings.Join(partsHex, partsSeparator)
}
