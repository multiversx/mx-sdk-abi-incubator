package abi

import (
	"errors"
	"io"
)

func (c *defaultCodec) decodeNestedOption(reader io.Reader, value *OptionValue) error {
	bytes, err := readBytesExactly(reader, 1)
	if err != nil {
		return err
	}

	if bytes[0] == 0 {
		value.Value = nil
		return nil
	}

	return c.doDecodeNested(reader, value.Value)
}

func (c *defaultCodec) decodeNestedList(reader io.Reader, value *OutputListValue) error {
	if value.ItemCreator == nil {
		return errors.New("cannot deserialize list: item creator is nil")
	}

	length, err := c.decodeLength(reader)
	if err != nil {
		return err
	}

	value.Items = make([]any, 0, length)

	for i := uint32(0); i < length; i++ {
		newItem := value.ItemCreator()

		err := c.doDecodeNested(reader, newItem)
		if err != nil {
			return err
		}

		value.Items = append(value.Items, newItem)
	}

	return nil
}
