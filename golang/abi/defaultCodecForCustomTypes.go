package abi

import (
	"bytes"
	"fmt"
	"io"
)

// https://docs.multiversx.com/developers/data/custom-types
func (c *defaultCodec) encodeNestedStruct(writer io.Writer, value StructValue) error {
	for _, field := range value.Fields {
		err := c.doEncodeNested(writer, field.Value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *defaultCodec) encodeTopLevelStruct(writer io.Writer, value StructValue) error {
	return c.encodeNestedStruct(writer, value)
}

func (c *defaultCodec) decodeNestedStruct(reader io.Reader, value *StructValue) error {
	for _, field := range value.Fields {
		err := c.doDecodeNested(reader, field.Value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *defaultCodec) decodeTopLevelStruct(data []byte, value *StructValue) error {
	reader := bytes.NewReader(data)
	return c.decodeNestedStruct(reader, value)
}

func (c *defaultCodec) encodeNestedEnum(writer io.Writer, value EnumValue) error {
	err := c.doEncodeNested(writer, U8Value{value.Discriminant})
	if err != nil {
		return err
	}

	for _, field := range value.Fields {
		err := c.doEncodeNested(writer, field.Value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *defaultCodec) encodeTopLevelEnum(writer io.Writer, value EnumValue) error {
	if value.Discriminant == 0 && len(value.Fields) == 0 {
		// Write nothing
		return nil
	}

	return c.encodeNestedEnum(writer, value)
}

func (c *defaultCodec) decodeNestedEnum(reader io.Reader, value *EnumValue) error {
	discriminant := &U8Value{}
	err := c.doDecodeNested(reader, discriminant)
	if err != nil {
		return err
	}

	value.Discriminant = discriminant.Value

	for _, field := range value.Fields {
		err := c.doDecodeNested(reader, field.Value)
		if err != nil {
			return fmt.Errorf("cannot decode field '%s' of struct, because of: %w", field.Name, err)
		}
	}

	return nil
}

func (c *defaultCodec) decodeTopLevelEnum(data []byte, value *EnumValue) error {
	if len(data) == 0 {
		value.Discriminant = 0
		return nil
	}

	reader := bytes.NewReader(data)
	return c.decodeNestedEnum(reader, value)
}
