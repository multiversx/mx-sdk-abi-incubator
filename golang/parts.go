package abi

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

type partsHolder struct {
	parts [][]byte
}

// A newly-created holder has no parts.
// Parts are created by calling appendEmptyPart().
func newEmptyPartsHolder() *partsHolder {
	return &partsHolder{
		parts: [][]byte{},
	}
}

func newPartsHolderFromHex(encoded string) (*partsHolder, error) {
	partsHex := strings.Split(encoded, partsSeparator)
	parts := make([][]byte, len(partsHex))

	for i, partHex := range partsHex {
		part, err := hex.DecodeString(partHex)
		if err != nil {
			return nil, err
		}

		parts[i] = part
	}

	return &partsHolder{
		parts: parts,
	}, nil
}

func (holder *partsHolder) hasAnyPart() bool {
	return len(holder.parts) > 0
}

func (holder *partsHolder) getNumParts() uint32 {
	return uint32(len(holder.parts))
}

func (holder *partsHolder) getPart(index uint32) ([]byte, error) {
	if index >= holder.getNumParts() {
		return nil, fmt.Errorf("part index %d is out of range", index)
	}

	return holder.parts[index], nil
}

func (holder *partsHolder) appendToLastPart(data []byte) error {
	if !holder.hasAnyPart() {
		return errors.New("cannot write, since there is no part to write to")
	}

	holder.parts[len(holder.parts)-1] = append(holder.parts[len(holder.parts)-1], data...)
	return nil
}

func (holder *partsHolder) appendEmptyPart() {
	holder.parts = append(holder.parts, []byte{})
}

func (holder *partsHolder) getParts() [][]byte {
	return holder.parts
}

func (holder *partsHolder) encodeToHex() string {
	partsHex := make([]string, len(holder.parts))

	for i, part := range holder.parts {
		partsHex[i] = hex.EncodeToString(part)
	}

	return strings.Join(partsHex, partsSeparator)
}

type partsReader struct {
	holder           *partsHolder
	currentPartIndex uint32
}

func newPartsReader(holder *partsHolder) *partsReader {
	return &partsReader{
		holder:           holder,
		currentPartIndex: 0,
	}
}

func (reader *partsReader) readWholePart() ([]byte, error) {
	if reader.isEndOfData() {
		return nil, fmt.Errorf("cannot wholly read part %d: unexpected end of data", reader.currentPartIndex)
	}

	part, err := reader.holder.getPart(uint32(reader.currentPartIndex))
	if err != nil {
		return nil, err
	}

	return part, nil
}

func (reader *partsReader) gotoNextPart() error {
	if reader.isEndOfData() {
		return fmt.Errorf(
			"cannot advance to next part, since the reader is already beyond the last part; current part is %d",
			reader.currentPartIndex,
		)
	}

	reader.currentPartIndex++
	return nil
}

// isEndOfData returns true if the reader is already beyond the last part.
func (reader *partsReader) isEndOfData() bool {
	return reader.currentPartIndex >= reader.holder.getNumParts()
}
