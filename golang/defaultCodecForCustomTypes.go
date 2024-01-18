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
func (c *defaultCodec) decodeNestedStruct(reader dataReader, value *StructValue) error {
	for _, field := range value.Fields {
		err := c.DecodeNested(reader, field.Value)
		if err != nil {
			return err
		}
	}

	return nil
}

// See: https://docs.multiversx.com/developers/data/custom-types
func (c *defaultCodec) decodeTopLevelStruct(reader dataReader, value *StructValue) error {
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
func (c *defaultCodec) decodeNestedEnum(reader dataReader, value *EnumValue) error {
	discriminant := &U8Value{}
	err := c.DecodeNested(reader, discriminant)
	if err != nil {
		return err
	}

	value.Discriminant = discriminant.Value

	for _, field := range value.Fields {
		err := c.DecodeNested(reader, field.Value)
		if err != nil {
			return err
		}
	}

	return nil
}

// See: https://docs.multiversx.com/developers/data/custom-types
func (c *defaultCodec) decodeTopLevelEnum(reader dataReader, value *EnumValue) error {
	if reader.IsCurrentPartEmpty() {
		value.Discriminant = 0
		return nil
	}

	return c.decodeNestedEnum(reader, value)
}
