package abi

import "fmt"

type partReader struct {
	part   []byte
	offset int
}

func newPartReader(part []byte) *partReader {
	return &partReader{
		part: part,
	}
}

func (r *partReader) Read(numBytes int) ([]byte, error) {
	if r.offset+numBytes > len(r.part) {
		return nil, fmt.Errorf("cannot read %d bytes from offset %d: unexpected end of data", numBytes, r.offset)
	}

	data := r.part[r.offset : r.offset+numBytes]
	r.offset += numBytes

	return data, nil
}

func (r *partReader) ReadWhole() ([]byte, error) {
	if r.offset != 0 {
		return nil, fmt.Errorf("cannot wholly read part, since current reading offset (%d) is not zero", r.offset)
	}

	part := r.part
	r.offset = len(part)
	return part, nil
}

func (r *partReader) IsPartEmpty() bool {
	return len(r.part) == 0
}
