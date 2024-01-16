package abi

type serializer struct {
	codec codec
}

func NewSerializer(codec codec) *serializer {
	return &serializer{
		codec: codec,
	}
}

func (s *serializer) Serialize(writer dataWriter, inputValues []interface{}) error {
	var err error

	for i, value := range inputValues {
		if value == nil {
			return errNilInputValue
		}

		switch value.(type) {
		case MultiValue:
			err = s.serializeMultiValue(writer, value.(MultiValue))
		case InputVariadicValues:
			if i != len(inputValues)-1 {
				return errVariadicMustBeLast
			}

			err = s.serializeInputVariadicValues(writer, value.(InputVariadicValues))
		default:
			err = s.serializeDirectlyEncodableValue(writer, value)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *serializer) Deserialize(reader dataReader, outputValues []interface{}) error {
	var err error

	for i, value := range outputValues {
		if value == nil {
			return errNilOutputValue
		}

		switch value.(type) {
		case *MultiValue:
			err = s.deserializeMultiValue(reader, value.(*MultiValue))
		case *OutputVariadicValues:
			if i != len(outputValues)-1 {
				return errVariadicMustBeLast
			}

			err = s.deserializeOutputVariadicValues(reader, value.(*OutputVariadicValues))
		default:
			err = s.deserializeDirectlyEncodableValue(reader, value)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *serializer) serializeMultiValue(writer dataWriter, value MultiValue) error {
	for _, item := range value.Items {
		writer.GotoNextPart()

		err := s.codec.EncodeTopLevel(writer, item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *serializer) serializeInputVariadicValues(writer dataWriter, value InputVariadicValues) error {
	for _, item := range value.Items {
		err := s.Serialize(writer, []interface{}{item})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *serializer) serializeDirectlyEncodableValue(writer dataWriter, value interface{}) error {
	writer.GotoNextPart()

	err := s.codec.EncodeTopLevel(writer, value)
	if err != nil {
		return err
	}

	return nil
}

func (s *serializer) deserializeMultiValue(reader dataReader, value *MultiValue) error {
	for _, item := range value.Items {
		err := s.deserializeDirectlyEncodableValue(reader, item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *serializer) deserializeOutputVariadicValues(reader dataReader, value *OutputVariadicValues) error {
	if value.ItemCreator == nil {
		return errNilItemCreator
	}

	for reader.HasUnreadData() {
		newItem := value.ItemCreator()

		err := s.Deserialize(reader, []interface{}{newItem})
		if err != nil {
			return err
		}

		value.Items = append(value.Items, newItem)
	}

	return nil
}

func (s *serializer) deserializeDirectlyEncodableValue(reader dataReader, value interface{}) error {
	err := s.codec.DecodeTopLevel(reader, value)
	if err != nil {
		return err
	}

	err = reader.GotoNextPart()
	if err != nil {
		return err
	}

	return nil
}
