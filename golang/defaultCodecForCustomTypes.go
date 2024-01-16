package abi

// https://docs.multiversx.com/developers/data/custom-types
func (c *defaultCodec) encodeNestedStruct(writer dataWriter, value StructValue) error {
	for _, field := range value.Fields {
		err := c.EncodeNested(writer, field.Value)
		if err != nil {
			return err
		}
	}

	return nil
}

// https://docs.multiversx.com/developers/data/custom-types
func (c *defaultCodec) encodeTopLevelStruct(writer dataWriter, value StructValue) error {
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
func (c *defaultCodec) encodeNestedEnum(writer dataWriter, value EnumValue) error {
	err := c.EncodeNested(writer, U8Value{value.Discriminant})
	if err != nil {
		return err
	}

	for _, field := range value.Fields {
		err := c.EncodeNested(writer, field.Value)
		if err != nil {
			return err
		}
	}

	return nil
}

// See: https://docs.multiversx.com/developers/data/custom-types
func (c *defaultCodec) encodeTopLevelEnum(writer dataWriter, value EnumValue) error {
	if value.Discriminant == 0 && len(value.Fields) == 0 {
		// Write nothing
		return nil
	}

	return c.encodeNestedEnum(writer, value)
}

// See: https://docs.multiversx.com/developers/data/custom-types
func (c *defaultCodec) decodeNestedEnum(reader dataReader, value *EnumValue) error {
	err := c.DecodeNested(reader, &value.Discriminant)
	if err != nil {
		return err
	}

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
