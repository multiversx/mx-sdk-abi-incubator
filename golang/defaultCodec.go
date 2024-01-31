package abi

import (
	"bytes"
	"fmt"
	"io"
)

type defaultCodec struct {
}

func NewDefaultCodec() *defaultCodec {
	return &defaultCodec{}
}

func (c *defaultCodec) EncodeNested(value interface{}) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	err := c.doEncodeNested(buffer, value)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (c *defaultCodec) doEncodeNested(writer io.Writer, value interface{}) error {
	switch value.(type) {
	case U8Value:
		return c.encodeNestedNumber(writer, value.(U8Value).Value, 1)
	case U16Value:
		return c.encodeNestedNumber(writer, value.(U16Value).Value, 2)
	case U32Value:
		return c.encodeNestedNumber(writer, value.(U32Value).Value, 4)
	case U64Value:
		return c.encodeNestedNumber(writer, value.(U64Value).Value, 8)
	case I8Value:
		return c.encodeNestedNumber(writer, value.(I8Value).Value, 1)
	case I16Value:
		return c.encodeNestedNumber(writer, value.(I16Value).Value, 2)
	case I32Value:
		return c.encodeNestedNumber(writer, value.(I32Value).Value, 4)
	case I64Value:
		return c.encodeNestedNumber(writer, value.(I64Value).Value, 8)
	case StructValue:
		return c.encodeNestedStruct(writer, value.(StructValue))
	case EnumValue:
		return c.encodeNestedEnum(writer, value.(EnumValue))
	default:
		return fmt.Errorf("unsupported type for nested encoding: %T", value)
	}
}

func (c *defaultCodec) EncodeTopLevel(value interface{}) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	err := c.doEncodeTopLevel(buffer, value)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (c *defaultCodec) doEncodeTopLevel(writer io.Writer, value interface{}) error {
	switch value.(type) {
	case U8Value:
		return c.encodeTopLevelUnsignedNumber(writer, uint64(value.(U8Value).Value))
	case U16Value:
		return c.encodeTopLevelUnsignedNumber(writer, uint64(value.(U16Value).Value))
	case U32Value:
		return c.encodeTopLevelUnsignedNumber(writer, uint64(value.(U32Value).Value))
	case U64Value:
		return c.encodeTopLevelUnsignedNumber(writer, value.(U64Value).Value)
	case I8Value:
		return c.encodeTopLevelSignedNumber(writer, int64(value.(I8Value).Value))
	case I16Value:
		return c.encodeTopLevelSignedNumber(writer, int64(value.(I16Value).Value))
	case I32Value:
		return c.encodeTopLevelSignedNumber(writer, int64(value.(I32Value).Value))
	case I64Value:
		return c.encodeTopLevelSignedNumber(writer, value.(I64Value).Value)
	case StructValue:
		return c.encodeTopLevelStruct(writer, value.(StructValue))
	case EnumValue:
		return c.encodeTopLevelEnum(writer, value.(EnumValue))
	default:
		return fmt.Errorf("unsupported type for top-level encoding: %T", value)
	}
}

func (c *defaultCodec) DecodeNested(data []byte, value interface{}) error {
	reader := bytes.NewReader(data)
	err := c.doDecodeNested(reader, value)
	if err != nil {
		return fmt.Errorf("cannot decode (nested) %T, because of: %w", value, err)
	}

	return nil
}

func (c *defaultCodec) doDecodeNested(reader io.Reader, value interface{}) error {
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
		return fmt.Errorf("unsupported type for nested decoding: %T", value)
	}
}

func (c *defaultCodec) DecodeTopLevel(data []byte, value interface{}) error {
	err := c.doDecodeTopLevel(data, value)
	if err != nil {
		return fmt.Errorf("cannot decode (top-level) %T, because of: %w", value, err)
	}

	return nil
}

func (c *defaultCodec) doDecodeTopLevel(data []byte, value interface{}) error {
	switch value.(type) {
	case *U8Value:
		return c.decodeTopLevelU8(data, value.(*U8Value))
	case *U16Value:
		return c.decodeTopLevelU16(data, value.(*U16Value))
	case *U32Value:
		return c.decodeTopLevelU32(data, value.(*U32Value))
	case *U64Value:
		return c.decodeTopLevelU64(data, value.(*U64Value))
	case *StructValue:
		return c.decodeTopLevelStruct(data, value.(*StructValue))
	case *EnumValue:
		return c.decodeTopLevelEnum(data, value.(*EnumValue))
	default:
		return fmt.Errorf("unsupported type for top-level decoding: %T", value)
	}
}
