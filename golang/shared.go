package abi

import (
	"fmt"
	"io"
)

func readBytesExactly(reader io.Reader, numBytes int) ([]byte, error) {
	data := make([]byte, numBytes)
	n, err := reader.Read(data)
	if err != nil {
		return nil, err
	}
	if n != numBytes {
		return nil, fmt.Errorf("cannot read exactly %d bytes", numBytes)
	}

	return data, err
}
