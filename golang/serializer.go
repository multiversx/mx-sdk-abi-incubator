package abi

import (
	"errors"
	"io"
)

type serializer struct {
	codec codec
}

func NewSerializer(codec codec) *serializer {
	return &serializer{
		codec: codec,
	}
}

func (s *serializer) Serialize(inputValues []interface{}) (string, error) {
	writer := newDataWriter()

	err := s.doSerialize(writer, inputValues)
	if err != nil {
		return "", err
	}

	return writer.String(), nil
}

func (s *serializer) doSerialize(writer *dataWriter, inputValues []interface{}) error {
	var err error

	for i, value := range inputValues {
		if value == nil {
			return errors.New("cannot serialize nil value")
		}

		switch value.(type) {
		case InputMultiValue:
			err = s.serializeInputMultiValue(writer, value.(InputMultiValue))
		case InputVariadicValues:
			if i != len(inputValues)-1 {
				return errors.New("variadic values must be last among input values")
			}

			err = s.serializeInputVariadicValues(writer, value.(InputVariadicValues))
		default:
			writer.GotoNextPart()
			err = s.serializeDirectlyEncodableValue(writer, value)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *serializer) Deserialize(data string, outputValues []interface{}) error {
	reader, err := newDataReaderFromString(data)
	if err != nil {
		return err
	}

	err = s.doDeserialize(reader, outputValues)
	if err != nil {
		return err
	}

	return nil
}

func (s *serializer) doDeserialize(reader *dataReader, outputValues []interface{}) error {
	var err error

	for i, value := range outputValues {
		if value == nil {
			return errors.New("cannot deserialize into nil value")
		}

		switch value.(type) {
		case *OutputMultiValue:
			err = s.deserializeOutputMultiValue(reader, value.(*OutputMultiValue))
		case *OutputVariadicValues:
			if i != len(outputValues)-1 {
				return errors.New("variadic values must be last among output values")
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

func (s *serializer) serializeInputMultiValue(writer *dataWriter, value InputMultiValue) error {
	for _, item := range value.Items {
		err := s.doSerialize(writer, []interface{}{item})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *serializer) serializeInputVariadicValues(writer *dataWriter, value InputVariadicValues) error {
	for _, item := range value.Items {
		err := s.doSerialize(writer, []interface{}{item})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *serializer) serializeDirectlyEncodableValue(writer io.Writer, value interface{}) error {
	data, err := s.codec.EncodeTopLevel(value)
	if err != nil {
		return err
	}

	_, err = writer.Write(data)
	return err
}

func (s *serializer) deserializeOutputMultiValue(reader *dataReader, value *OutputMultiValue) error {
	for _, item := range value.Items {
		err := s.doDeserialize(reader, []interface{}{item})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *serializer) deserializeOutputVariadicValues(reader *dataReader, value *OutputVariadicValues) error {
	if value.ItemCreator == nil {
		return errors.New("cannot deserialize variadic values: item creator is nil")
	}

	for !reader.IsEndOfData() {
		newItem := value.ItemCreator()

		err := s.doDeserialize(reader, []interface{}{newItem})
		if err != nil {
			return err
		}

		value.Items = append(value.Items, newItem)
	}

	return nil
}

func (s *serializer) deserializeDirectlyEncodableValue(reader *dataReader, value interface{}) error {
	part, err := reader.ReadWholePart()
	if err != nil {
		return err
	}

	err = s.codec.DecodeTopLevel(part, value)
	if err != nil {
		return err
	}

	err = reader.GotoNextPart()
	if err != nil {
		return err
	}

	return nil
}
