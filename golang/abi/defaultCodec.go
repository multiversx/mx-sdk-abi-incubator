package abi

import (
	"bytes"
	"fmt"
	"io"
	"math"
)

// defaultCodec is the default codec for encoding and decoding values.
//
// See:
// - https://docs.multiversx.com/developers/data/simple-values
// - https://docs.multiversx.com/developers/data/composite-values
// - https://docs.multiversx.com/developers/data/custom-types
type defaultCodec struct {
}

// NewDefaultCodec creates a new default codec.
func NewDefaultCodec() *defaultCodec {
	return &defaultCodec{}
}

func (c *defaultCodec) EncodeNested(value any) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	err := c.doEncodeNested(buffer, value)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (c *defaultCodec) doEncodeNested(writer io.Writer, value any) error {
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
	case BigIntValue:
		return c.encodeNestedBigNumber(writer, value.(BigIntValue).Value)
	case StructValue:
		return c.encodeNestedStruct(writer, value.(StructValue))
	case EnumValue:
		return c.encodeNestedEnum(writer, value.(EnumValue))
	default:
		return fmt.Errorf("unsupported type for nested encoding: %T", value)
	}
}

func (c *defaultCodec) EncodeTopLevel(value any) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	err := c.doEncodeTopLevel(buffer, value)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (c *defaultCodec) doEncodeTopLevel(writer io.Writer, value any) error {
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
	case BigIntValue:
		return c.encodeTopLevelBigNumber(writer, value.(BigIntValue).Value)
	case StructValue:
		return c.encodeTopLevelStruct(writer, value.(StructValue))
	case EnumValue:
		return c.encodeTopLevelEnum(writer, value.(EnumValue))
	default:
		return fmt.Errorf("unsupported type for top-level encoding: %T", value)
	}
}

func (c *defaultCodec) DecodeNested(data []byte, value any) error {
	reader := bytes.NewReader(data)
	err := c.doDecodeNested(reader, value)
	if err != nil {
		return fmt.Errorf("cannot decode (nested) %T, because of: %w", value, err)
	}

	return nil
}

func (c *defaultCodec) doDecodeNested(reader io.Reader, value any) error {
	switch value.(type) {
	case *U8Value:
		return c.decodeNestedNumber(reader, &value.(*U8Value).Value, 1)
	case *U16Value:
		return c.decodeNestedNumber(reader, &value.(*U16Value).Value, 2)
	case *U32Value:
		return c.decodeNestedNumber(reader, &value.(*U32Value).Value, 4)
	case *U64Value:
		return c.decodeNestedNumber(reader, &value.(*U64Value).Value, 8)
	case *I8Value:
		return c.decodeNestedNumber(reader, &value.(*I8Value).Value, 1)
	case *I16Value:
		return c.decodeNestedNumber(reader, &value.(*I16Value).Value, 2)
	case *I32Value:
		return c.decodeNestedNumber(reader, &value.(*I32Value).Value, 4)
	case *I64Value:
		return c.decodeNestedNumber(reader, &value.(*I64Value).Value, 8)
	case *BigIntValue:
		n, err := c.decodeNestedBigNumber(reader)
		if err != nil {
			return err
		}

		value.(*BigIntValue).Value = n
		return nil
	case *StructValue:
		return c.decodeNestedStruct(reader, value.(*StructValue))
	case *EnumValue:
		return c.decodeNestedEnum(reader, value.(*EnumValue))
	case *OptionValue:
		return c.decodeNestedOption(reader, value.(*OptionValue))
	case *OutputListValue:
		return c.decodeNestedList(reader, value.(*OutputListValue))
	default:
		return fmt.Errorf("unsupported type for nested decoding: %T", value)
	}
}

func (c *defaultCodec) DecodeTopLevel(data []byte, value any) error {
	err := c.doDecodeTopLevel(data, value)
	if err != nil {
		return fmt.Errorf("cannot decode (top-level) %T, because of: %w", value, err)
	}

	return nil
}

func (c *defaultCodec) doDecodeTopLevel(data []byte, value any) error {
	switch value.(type) {
	case *U8Value:
		n, err := c.decodeTopLevelUnsignedNumber(data, math.MaxUint8)
		if err != nil {
			return err
		}

		value.(*U8Value).Value = uint8(n)
	case *U16Value:
		n, err := c.decodeTopLevelUnsignedNumber(data, math.MaxUint16)
		if err != nil {
			return err
		}

		value.(*U16Value).Value = uint16(n)
	case *U32Value:
		n, err := c.decodeTopLevelUnsignedNumber(data, math.MaxUint32)
		if err != nil {
			return err
		}

		value.(*U32Value).Value = uint32(n)
	case *U64Value:
		n, err := c.decodeTopLevelUnsignedNumber(data, math.MaxUint64)
		if err != nil {
			return err
		}

		value.(*U64Value).Value = uint64(n)
	case *I8Value:
		n, err := c.decodeTopLevelSignedNumber(data, math.MaxInt8)
		if err != nil {
			return err
		}

		value.(*I8Value).Value = int8(n)
	case *I16Value:
		n, err := c.decodeTopLevelSignedNumber(data, math.MaxInt16)
		if err != nil {
			return err
		}

		value.(*I16Value).Value = int16(n)
	case *I32Value:
		n, err := c.decodeTopLevelSignedNumber(data, math.MaxInt32)
		if err != nil {
			return err
		}

		value.(*I32Value).Value = int32(n)

	case *I64Value:
		n, err := c.decodeTopLevelSignedNumber(data, math.MaxInt64)
		if err != nil {
			return err
		}

		value.(*I64Value).Value = int64(n)
	case *BigIntValue:
		n := c.decodeTopLevelBigNumber(data)
		value.(*BigIntValue).Value = n
	case *StructValue:
		return c.decodeTopLevelStruct(data, value.(*StructValue))
	case *EnumValue:
		return c.decodeTopLevelEnum(data, value.(*EnumValue))
	default:
		return fmt.Errorf("unsupported type for top-level decoding: %T", value)
	}

	return nil
}
