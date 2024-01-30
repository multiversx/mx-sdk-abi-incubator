package abi

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

type partsHolder struct {
	parts            [][]byte
	focusedPartIndex uint32
}

// A newly-created holder has no parts.
// Parts are created by calling appendEmptyPart().
func newEmptyPartsHolder() *partsHolder {
	return &partsHolder{
		parts:            [][]byte{},
		focusedPartIndex: 0,
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

func (holder *partsHolder) getParts() [][]byte {
	return holder.parts
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

func (holder *partsHolder) hasAnyPart() bool {
	return len(holder.parts) > 0
}

func (holder *partsHolder) appendEmptyPart() {
	holder.parts = append(holder.parts, []byte{})
}

func (holder *partsHolder) encodeToHex() string {
	partsHex := make([]string, len(holder.parts))

	for i, part := range holder.parts {
		partsHex[i] = hex.EncodeToString(part)
	}

	return strings.Join(partsHex, partsSeparator)
}

func (holder *partsHolder) readWholePart() ([]byte, error) {
	if holder.isFocusedBeyondLastPart() {
		return nil, fmt.Errorf("cannot wholly read part %d: unexpected end of data", holder.focusedPartIndex)
	}

	part, err := holder.getPart(uint32(holder.focusedPartIndex))
	if err != nil {
		return nil, err
	}

	return part, nil
}

func (holder *partsHolder) focusOnNextPart() error {
	if holder.isFocusedBeyondLastPart() {
		return fmt.Errorf(
			"cannot focus on next part, since the focus is already beyond the last part; focused part index is %d",
			holder.focusedPartIndex,
		)
	}

	holder.focusedPartIndex++
	return nil
}

// isFocusedBeyondLastPart returns true if the reader is already beyond the last part.
func (holder *partsHolder) isFocusedBeyondLastPart() bool {
	return holder.focusedPartIndex >= holder.getNumParts()
}
