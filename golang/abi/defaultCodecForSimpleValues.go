package abi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/big"

	twos "github.com/multiversx/mx-components-big-int/twos-complement"
)

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

func (c *defaultCodec) decodeTopLevelSignedNumber(data []byte, maxValue int64) (int64, error) {
	b := big.NewInt(0).SetBytes(data)
	if !b.IsInt64() {
		return 0, fmt.Errorf("decoded value is too large (does not fit an int64): %s", b)
	}

	n := b.Int64()
	if n > maxValue {
		return 0, fmt.Errorf("decoded value is too large: %d > %d", n, maxValue)
	}

	return n, nil
}

func (c *defaultCodec) encodeNestedBigNumber(writer io.Writer, value *big.Int) error {
	data := twos.ToBytes(value)
	dataLength := len(data)

	// Write the length of the payload
	err := c.encodeLength(writer, uint32(dataLength))
	if err != nil {
		return err
	}

	// Write the payload
	_, err = writer.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (c *defaultCodec) encodeTopLevelBigNumber(writer io.Writer, value *big.Int) error {
	data := twos.ToBytes(value)
	_, err := writer.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (c *defaultCodec) decodeNestedBigNumber(reader io.Reader) (*big.Int, error) {
	// Read the length of the payload
	length, err := c.decodeLength(reader)
	if err != nil {
		return nil, err
	}

	// Read the payload
	data, err := readBytesExactly(reader, int(length))
	if err != nil {
		return nil, err
	}

	return twos.FromBytes(data), nil
}

func (c *defaultCodec) decodeTopLevelBigNumber(data []byte) *big.Int {
	return twos.FromBytes(data)
}

func (c *defaultCodec) encodeLength(writer io.Writer, length uint32) error {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, length)

	_, err := writer.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}

func (c *defaultCodec) decodeLength(reader io.Reader) (uint32, error) {
	bytes, err := readBytesExactly(reader, 4)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint32(bytes), nil
}

func (c *defaultCodec) encodeNestedString(writer io.Writer, value StringValue) error {
	data := []byte(value.Value)
	err := c.encodeLength(writer, uint32(len(data)))
	if err != nil {
		return err
	}

	_, err = writer.Write(data)
	return err
}

func (c *defaultCodec) decodeNestedString(reader io.Reader, value *StringValue) error {
	length, err := c.decodeLength(reader)
	if err != nil {
		return err
	}

	data, err := readBytesExactly(reader, int(length))
	if err != nil {
		return err
	}

	value.Value = string(data)
	return nil
}

func (c *defaultCodec) encodeNestedBytes(writer io.Writer, value BytesValue) error {
	err := c.encodeLength(writer, uint32(len(value.Value)))
	if err != nil {
		return err
	}

	_, err = writer.Write(value.Value)
	return err
}

func (c *defaultCodec) decodeNestedBytes(reader io.Reader, value *BytesValue) error {
	length, err := c.decodeLength(reader)
	if err != nil {
		return err
	}

	data, err := readBytesExactly(reader, int(length))
	if err != nil {
		return err
	}

	value.Value = data
	return nil
}
