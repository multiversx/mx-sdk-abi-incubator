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
		case InputMultiValue:
			err = s.serializeInputMultiValue(writer, value.(InputMultiValue))
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
		case *OutputMultiValue:
			err = s.deserializeOutputMultiValue(reader, value.(*OutputMultiValue))
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

func (s *serializer) serializeInputMultiValue(writer dataWriter, value InputMultiValue) error {
	for _, item := range value.Items {
		err := s.Serialize(writer, []interface{}{item})
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

func (s *serializer) deserializeOutputMultiValue(reader dataReader, value *OutputMultiValue) error {
	for _, item := range value.Items {
		err := s.Deserialize(reader, []interface{}{item})
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

	for !reader.IsEndOfData() {
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
