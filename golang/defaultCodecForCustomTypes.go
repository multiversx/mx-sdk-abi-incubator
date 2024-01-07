package abi

func (c *defaultCodec) encodeNestedStruct(writer dataWriter, value StructValue) error {
	for _, field := range value.Fields {
		err := c.EncodeNested(writer, field.Value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *defaultCodec) encodeTopLevelStruct(writer dataWriter, value StructValue) error {
	return c.encodeNestedStruct(writer, value)
}

func (c *defaultCodec) decodeNestedStruct(reader dataReader, value *StructValue) error {
	for _, field := range value.Fields {
		err := c.DecodeNested(reader, field.Value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *defaultCodec) decodeTopLevelStruct(reader dataReader, value *StructValue) error {
	return c.decodeNestedStruct(reader, value)
}
