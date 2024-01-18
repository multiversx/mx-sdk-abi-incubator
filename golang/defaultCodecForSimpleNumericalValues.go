package abi

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/big"
)

func (c *defaultCodec) encodeNestedU8(writer io.Writer, value U8Value) error {
	_, err := writer.Write([]byte{value.Value})
	return err
}

func (c *defaultCodec) encodeNestedI8(writer io.Writer, value I8Value) error {
	_, err := writer.Write([]byte{byte(value.Value)})
	return err
}

func (c *defaultCodec) encodeTopLevelU8(writer io.Writer, value U8Value) error {
	return c.encodeTopLevelUnsignedNumber(writer, uint64(value.Value))
}

func (c *defaultCodec) encodeTopLevelI8(writer io.Writer, value I8Value) error {
	return c.encodeTopLevelSignedNumber(writer, int64(value.Value))
}

func (c *defaultCodec) decodeNestedU8(reader dataReader, value *U8Value) error {
	data, err := reader.Read(1)
	if err != nil {
		return newErrCodecCannotDecodeType("u8", err)
	}

	value.Value = data[0]
	return nil
}

func (c *defaultCodec) decodeTopLevelU8(reader dataReader, value *U8Value) error {
	n, err := c.decodeTopLevelUnsignedNumber(reader, math.MaxUint8)
	if err != nil {
		return newErrCodecCannotDecodeType("u8", err)
	}

	value.Value = uint8(n)
	return nil
}

func (c *defaultCodec) encodeNestedU16(writer io.Writer, value U16Value) error {
	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data, value.Value)
	_, err := writer.Write(data)
	return err
}

func (c *defaultCodec) encodeTopLevelU16(writer io.Writer, value U16Value) error {
	return c.encodeTopLevelUnsignedNumber(writer, uint64(value.Value))
}

func (c *defaultCodec) decodeNestedU16(reader dataReader, value *U16Value) error {
	data, err := reader.Read(2)
	if err != nil {
		return newErrCodecCannotDecodeType("u16", err)
	}

	value.Value = binary.BigEndian.Uint16(data)
	return nil
}

func (c *defaultCodec) decodeTopLevelU16(reader dataReader, value *U16Value) error {
	n, err := c.decodeTopLevelUnsignedNumber(reader, math.MaxUint16)
	if err != nil {
		return newErrCodecCannotDecodeType("u16", err)
	}

	value.Value = uint16(n)
	return nil
}

func (c *defaultCodec) encodeNestedU32(writer io.Writer, value U32Value) error {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, value.Value)
	_, err := writer.Write(data)
	return err
}

func (c *defaultCodec) encodeTopLevelU32(writer io.Writer, value U32Value) error {
	return c.encodeTopLevelUnsignedNumber(writer, uint64(value.Value))
}

func (c *defaultCodec) decodeNestedU32(reader dataReader, value *U32Value) error {
	data, err := reader.Read(4)
	if err != nil {
		return newErrCodecCannotDecodeType("u32", err)
	}

	value.Value = binary.BigEndian.Uint32(data)
	return nil
}

func (c *defaultCodec) decodeTopLevelU32(reader dataReader, value *U32Value) error {
	n, err := c.decodeTopLevelUnsignedNumber(reader, math.MaxUint32)
	if err != nil {
		return newErrCodecCannotDecodeType("u32", err)
	}

	value.Value = uint32(n)
	return nil
}

func (c *defaultCodec) encodeNestedU64(writer io.Writer, value U64Value) error {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, value.Value)
	_, err := writer.Write(data)
	return err
}

func (c *defaultCodec) encodeTopLevelU64(writer io.Writer, value U64Value) error {
	return c.encodeTopLevelUnsignedNumber(writer, uint64(value.Value))
}

func (c *defaultCodec) decodeNestedU64(reader dataReader, value *U64Value) error {
	data, err := reader.Read(8)
	if err != nil {
		return newErrCodecCannotDecodeType("u64", err)
	}

	value.Value = binary.BigEndian.Uint64(data)
	return nil
}

func (c *defaultCodec) decodeTopLevelU64(reader dataReader, value *U64Value) error {
	n, err := c.decodeTopLevelUnsignedNumber(reader, math.MaxUint64)
	if err != nil {
		return newErrCodecCannotDecodeType("u64", err)
	}

	value.Value = uint64(n)
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

func (c *defaultCodec) decodeTopLevelUnsignedNumber(reader dataReader, maxValue uint64) (uint64, error) {
	data, err := reader.ReadWholePart()
	if err != nil {
		return 0, err
	}

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
