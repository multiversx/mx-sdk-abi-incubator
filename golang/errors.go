package abi

import (
	"errors"
	"fmt"
)

var errNilInputValue = errors.New("nil input value")
var errNilOutputValue = errors.New("nil output value")
var errNilItemCreator = errors.New("nil item creator")
var errWriterCannotWriteSinceThereIsNoPart = errors.New("cannot write, since there is no part to write to")
var errVariadicMustBeLast = errors.New("variadic must be last")

func newErrReaderCannotReadDueToEndOfDataInPart(numBytes int, offsetInPart int, partIndex int) error {
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

func newErrReaderCannotReadWholePartDueToEndOfData(partIndex int) error {
	return fmt.Errorf(
		"cannot wholly read part %d: unexpected end of data",
		partIndex,
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

func newErrReaderCannotGotoNextPartDueToEndOfData(currentPartIndex int) error {
	return fmt.Errorf(
		"cannot advance to next part, since the reader is already beyond the last part; current part is %d",
		currentPartIndex,
	)
}

func newErrUnsupportedType(when string, value interface{}) error {
	return fmt.Errorf("%s, unsupported type: %T", when, value)
}

func newErrCodecCannotDecodeType(typeAlias string, originalError error) error {
	return fmt.Errorf("cannot decode %s, because of: %w", typeAlias, originalError)
}
