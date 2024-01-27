package abi

import (
	"encoding/hex"
	"errors"
	"strings"
)

type defaultDataWriter struct {
	parts [][]byte
}

// NewDefaultDataWriter creates a new defaultDataWriter.
// A newly-created writer has no parts.
// Parts are created by calling GotoNextPart().
func NewDefaultDataWriter() *defaultDataWriter {
	return &defaultDataWriter{
		parts: [][]byte{},
	}
}

func (d *defaultDataWriter) Write(data []byte) (int, error) {
	if len(d.parts) == 0 {
		return 0, errors.New("cannot write, since there is no part to write to")
	}

	partIndex := len(d.parts) - 1
	d.parts[partIndex] = append(d.parts[partIndex], data...)
	return len(data), nil
}

func (d *defaultDataWriter) GotoNextPart() {
	d.parts = append(d.parts, []byte{})
}

func (d *defaultDataWriter) GetParts() [][]byte {
	return d.parts
}

func (d *defaultDataWriter) String() string {
	partsHex := make([]string, len(d.parts))

	for i, part := range d.parts {
		partsHex[i] = hex.EncodeToString(part)
	}

	return strings.Join(partsHex, partsSeparator)
}
