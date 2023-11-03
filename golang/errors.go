package abi

import "errors"

var errUnexpectedEndOfData = errors.New("unexpected end of data")
var errOffsetIsNotAtStartOfPart = errors.New("offset is not zero")
var errUnreadDataInPart = errors.New("unread data in part")
var errNoDataPart = errors.New("no data part")
var errVariadicMustBeLast = errors.New("variadic must be last")
