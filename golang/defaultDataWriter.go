package abi

import (
	"encoding/hex"
	"strings"
)

type defaultDataWriter struct {
	parts [][]byte
}

func NewDefaultDataWriter() *defaultDataWriter {
	return &defaultDataWriter{
		parts: [][]byte{},
	}
}

func (d *defaultDataWriter) Write(data []byte) error {
	if len(d.parts) == 0 {
		return errWriterCannotWriteSinceThereIsNoPart
	}

	partIndex := len(d.parts) - 1
	d.parts[partIndex] = append(d.parts[partIndex], data...)
	return nil
}

func (d *defaultDataWriter) GotoNextPart() {
	d.parts = append(d.parts, []byte{})
}

func (d *defaultDataWriter) String() string {
	partsHex := make([]string, len(d.parts))

	for i, part := range d.parts {
		partsHex[i] = hex.EncodeToString(part)
	}

	return strings.Join(partsHex, partsSeparator)
}
