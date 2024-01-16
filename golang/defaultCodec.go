package abi

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

type defaultCodec struct {
}

func NewDefaultCodec() *defaultCodec {
	return &defaultCodec{}
}

func (c *defaultCodec) EncodeNested(writer dataWriter, value interface{}) error {
	switch value.(type) {
	case U8Value:
		return c.encodeNestedU8(writer, value.(U8Value))
	case U16Value:
		return c.encodeNestedU16(writer, value.(U16Value))
	case U32Value:
		return c.encodeNestedU32(writer, value.(U32Value))
	case U64Value:
		return c.encodeNestedU64(writer, value.(U64Value))
	case StructValue:
		return c.encodeNestedStruct(writer, value.(StructValue))
	default:
		return newErrUnsupportedType(value)
	}
}

func (c *defaultCodec) EncodeTopLevel(writer dataWriter, value interface{}) error {
	switch value.(type) {
	case U8Value:
		return c.encodeTopLevelU8(writer, value.(U8Value))
	case U16Value:
		return c.encodeTopLevelU16(writer, value.(U16Value))
	case U32Value:
		return c.encodeTopLevelU32(writer, value.(U32Value))
	case U64Value:
		return c.encodeTopLevelU64(writer, value.(U64Value))
	case StructValue:
		return c.encodeTopLevelStruct(writer, value.(StructValue))
	default:
		return newErrUnsupportedType(value)
	}
}

func (c *defaultCodec) DecodeNested(reader dataReader, value interface{}) error {
	switch value.(type) {
	case *U8Value:
		return c.decodeNestedU8(reader, value.(*U8Value))
	case *U16Value:
		return c.decodeNestedU16(reader, value.(*U16Value))
	case *U32Value:
		return c.decodeNestedU32(reader, value.(*U32Value))
	case *U64Value:
		return c.decodeNestedU64(reader, value.(*U64Value))
	case *StructValue:
		return c.decodeNestedStruct(reader, value.(*StructValue))
	default:
		return newErrUnsupportedType(value)
	}
}

func (c *defaultCodec) DecodeTopLevel(reader dataReader, value interface{}) error {
	switch value.(type) {
	case *U8Value:
		return c.decodeTopLevelU8(reader, value.(*U8Value))
	case *U16Value:
		return c.decodeTopLevelU16(reader, value.(*U16Value))
	case *U32Value:
		return c.decodeTopLevelU32(reader, value.(*U32Value))
	case *U64Value:
		return c.decodeTopLevelU64(reader, value.(*U64Value))
	case *StructValue:
		return c.decodeTopLevelStruct(reader, value.(*StructValue))
	default:
		return newErrUnsupportedType(value)
	}
}

func (c *defaultCodec) encodeNestedU8(writer dataWriter, value U8Value) error {
	return writer.Write([]byte{value.Value})
}

func (c *defaultCodec) encodeTopLevelU8(writer dataWriter, value U8Value) error {
	return writer.Write(trimLeadingZeros([]byte{value.Value}))
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

func (c *defaultCodec) encodeNestedU16(writer dataWriter, value U16Value) error {
	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data, value.Value)
	return writer.Write(data)
}

func (c *defaultCodec) encodeTopLevelU16(writer dataWriter, value U16Value) error {
	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data, value.Value)
	return writer.Write(trimLeadingZeros(data))
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

func (c *defaultCodec) encodeNestedU32(writer dataWriter, value U32Value) error {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, value.Value)
	return writer.Write(data)
}

func (c *defaultCodec) encodeTopLevelU32(writer dataWriter, value U32Value) error {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, value.Value)
	return writer.Write(trimLeadingZeros(data))
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

func (c *defaultCodec) encodeNestedU64(writer dataWriter, value U64Value) error {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, value.Value)
	return writer.Write(data)
}

func (c *defaultCodec) encodeTopLevelU64(writer dataWriter, value U64Value) error {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, value.Value)
	return writer.Write(trimLeadingZeros(data))
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
