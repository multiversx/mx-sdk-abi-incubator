package abi

import "io"

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

// https://docs.multiversx.com/developers/data/custom-types
func (c *defaultCodec) encodeTopLevelStruct(writer io.Writer, value StructValue) error {
	return c.encodeNestedStruct(writer, value)
}

// See: https://docs.multiversx.com/developers/data/custom-types
func (c *defaultCodec) decodeNestedStruct(reader *partReader, value *StructValue) error {
	for _, field := range value.Fields {
		err := c.doDecodeNested(reader, field.Value)
		if err != nil {
			return err
		}
	}

	return nil
}

// See: https://docs.multiversx.com/developers/data/custom-types
func (c *defaultCodec) decodeTopLevelStruct(reader *partReader, value *StructValue) error {
	return c.decodeNestedStruct(reader, value)
}

// See: https://docs.multiversx.com/developers/data/custom-types
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

// See: https://docs.multiversx.com/developers/data/custom-types
func (c *defaultCodec) encodeTopLevelEnum(writer io.Writer, value EnumValue) error {
	if value.Discriminant == 0 && len(value.Fields) == 0 {
		// Write nothing
		return nil
	}

	return c.encodeNestedEnum(writer, value)
}

// See: https://docs.multiversx.com/developers/data/custom-types
func (c *defaultCodec) decodeNestedEnum(reader *partReader, value *EnumValue) error {
	discriminant := &U8Value{}
	err := c.doDecodeNested(reader, discriminant)
	if err != nil {
		return err
	}

	value.Discriminant = discriminant.Value

	for _, field := range value.Fields {
		err := c.doDecodeNested(reader, field.Value)
		if err != nil {
			return err
		}
	}

	return nil
}

// See: https://docs.multiversx.com/developers/data/custom-types
func (c *defaultCodec) decodeTopLevelEnum(reader *partReader, value *EnumValue) error {
	if reader.IsPartEmpty() {
		value.Discriminant = 0
		return nil
	}

	return c.decodeNestedEnum(reader, value)
}
