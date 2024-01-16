package abi

import "bytes"

func trimLeadingZeros(data []byte) []byte {
	return bytes.TrimLeft(data, string([]byte{0x00}))
}
