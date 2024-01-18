package abi

import (
	"bytes"
	"io"
)

type defaultCodec struct {
}

func NewDefaultCodec() *defaultCodec {
	return &defaultCodec{}
}

func (c *defaultCodec) EncodeNested(value interface{}) ([]byte, error) {
	writer := bytes.NewBuffer(nil)
	err := c.doEncodeNested(writer, value)
	if err != nil {
		return nil, err
	}

	return writer.Bytes(), nil
}

func (c *defaultCodec) doEncodeNested(writer io.Writer, value interface{}) error {
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
	case EnumValue:
		return c.encodeNestedEnum(writer, value.(EnumValue))
	default:
		return newErrUnsupportedType("defaultCodec.EncodeNested()", value)
	}
}

func (c *defaultCodec) EncodeTopLevel(value interface{}) ([]byte, error) {
	writer := bytes.NewBuffer(nil)
	err := c.doEncodeTopLevel(writer, value)
	if err != nil {
		return nil, err
	}

	return writer.Bytes(), nil
}

func (c *defaultCodec) doEncodeTopLevel(writer io.Writer, value interface{}) error {
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
	case EnumValue:
		return c.encodeTopLevelEnum(writer, value.(EnumValue))
	default:
		return newErrUnsupportedType("defaultCodec.EncodeTopLevel()", value)
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
	case *EnumValue:
		return c.decodeNestedEnum(reader, value.(*EnumValue))
	default:
		return newErrUnsupportedType("defaultCodec.DecodeNested()", value)
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
	case *EnumValue:
		return c.decodeTopLevelEnum(reader, value.(*EnumValue))
	default:
		return newErrUnsupportedType("defaultCodec.DecodeTopLevel()", value)
	}
}
