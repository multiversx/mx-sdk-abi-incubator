package abi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/big"
)

func (c *defaultCodec) decodeTopLevelU8(data []byte, value *U8Value) error {
	n, err := c.decodeTopLevelUnsignedNumber(data, math.MaxUint8)
	if err != nil {
		return err
	}

	value.Value = uint8(n)
	return nil
}

func (c *defaultCodec) decodeTopLevelU16(data []byte, value *U16Value) error {
	n, err := c.decodeTopLevelUnsignedNumber(data, math.MaxUint16)
	if err != nil {
		return err
	}

	value.Value = uint16(n)
	return nil
}

func (c *defaultCodec) decodeTopLevelU32(data []byte, value *U32Value) error {
	n, err := c.decodeTopLevelUnsignedNumber(data, math.MaxUint32)
	if err != nil {
		return err
	}

	value.Value = uint32(n)
	return nil
}

func (c *defaultCodec) decodeTopLevelU64(data []byte, value *U64Value) error {
	n, err := c.decodeTopLevelUnsignedNumber(data, math.MaxUint64)
	if err != nil {
		return err
	}

	value.Value = uint64(n)
	return nil
}

func (c *defaultCodec) encodeNestedNumber(writer io.Writer, value any, numBytes int) error {
	buffer := new(bytes.Buffer)

	err := binary.Write(buffer, binary.BigEndian, value)
	if err != nil {
		return err
	}

	data := buffer.Bytes()
	if len(data) != numBytes {
		return fmt.Errorf("unexpected number of bytes: %d != %d", len(data), numBytes)
	}

	_, err = writer.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (c *defaultCodec) decodeNestedNumber(reader io.Reader, value any, numBytes int) error {
	data, err := readBytesExactly(reader, numBytes)
	if err != nil {
		return err
	}

	buffer := bytes.NewReader(data)
	err = binary.Read(buffer, binary.BigEndian, value)
	if err != nil {
		return err
	}

	return nil
}

func (c *defaultCodec) encodeTopLevelUnsignedNumber(writer io.Writer, value uint64) error {
	b := big.NewInt(0).SetUint64(value)
	data := b.Bytes()
	_, err := writer.Write(data)
	return err
}

func (c *defaultCodec) encodeTopLevelSignedNumber(writer io.Writer, value int64) error {
	b := big.NewInt(0).SetInt64(value)
	data := b.Bytes()
	_, err := writer.Write(data)
	return err
}

func (c *defaultCodec) decodeTopLevelUnsignedNumber(data []byte, maxValue uint64) (uint64, error) {
	b := big.NewInt(0).SetBytes(data)
	if !b.IsUint64() {
		return 0, fmt.Errorf("decoded value is too large (does not fit an uint64): %s", b)
	}

	n := b.Uint64()
	if n > maxValue {
		return 0, fmt.Errorf("decoded value is too large: %d > %d", n, maxValue)
	}

	return n, nil
}
