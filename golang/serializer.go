package abi

type serializer struct {
	codec codec
}

func NewSerializer(codec codec) *serializer {
	return &serializer{
		codec: codec,
	}
}

func (s *serializer) Serialize(writer dataWriter, values []interface{}) error {
	var err error

	for i, value := range values {
		switch value.(type) {
		case CompositeValue:
			err = s.serializeCompositeValue(writer, value.(CompositeValue))
		case InputVariadicValues:
			if i != len(values)-1 {
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

func (s *serializer) Deserialize(reader dataReader, values []interface{}) error {
	var err error

	for i, value := range values {
		switch value.(type) {
		case *CompositeValue:
			err = s.deserializeCompositeValue(reader, value.(*CompositeValue))
		case *OutputVariadicValues:
			if i != len(values)-1 {
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

func (s *serializer) serializeCompositeValue(writer dataWriter, value CompositeValue) error {
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

func (s *serializer) deserializeCompositeValue(reader dataReader, value *CompositeValue) error {
	for _, item := range value.Items {
		err := s.deserializeDirectlyEncodableValue(reader, item)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *serializer) deserializeOutputVariadicValues(reader dataReader, value *OutputVariadicValues) error {
	for reader.HasNextPart() {
		newItem := deepClone(value.ItemPrototype)

		err := s.Deserialize(reader, []interface{}{newItem})
		if err != nil {
			return err
		}

		value.Items = append(value.Items, newItem)

		if !reader.HasNextPart() {
			break
		}

		err = reader.GotoNextPart()
		if err != nil {
			return err
		}
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
