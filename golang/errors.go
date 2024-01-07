package abi

import (
	"errors"
	"fmt"
)

var errWriterCannotWriteSinceThereIsNoPart = errors.New("cannot write, since there is no part to write to")

var errVariadicMustBeLast = errors.New("variadic must be last")

func newErrReaderCannotReadDueToEndOfData(numBytes int, offsetInPart int, partIndex int) error {
	return fmt.Errorf(
		"cannot read %d bytes from offset %d, within part %d: unexpected end of data",
		numBytes,
		offsetInPart,
		partIndex,
	)
}

func newErrReaderCannotReadWholePartDueToNonZeroOffset(partIndex, offsetInPart int) error {
	return fmt.Errorf(
		"cannot wholly read part %d, since current reading offset (%d) is not zero",
		partIndex,
		offsetInPart,
	)
}

func newErrReaderCannotGotoNextPartDueToUnreadDataInCurrentPart(numUnreadBytes int, currentPartIndex int, lengthOfCurrentPart int) error {
	return fmt.Errorf(
		"cannot advance to next part, since there is still unread data (%d bytes) in part %d (of length: %d)",
		numUnreadBytes,
		currentPartIndex,
		lengthOfCurrentPart,
	)
}

func newErrUnsupportedType(value interface{}) error {
	return fmt.Errorf("unsupported type: %T", value)
}

func newErrCodecCannotDecodeType(typeAlias string, originalError error) error {
	return fmt.Errorf("cannot decode %s, because of: %w", typeAlias, originalError)
}
