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
	partsHolder := newEmptyPartsHolder()

	err := s.doSerialize(partsHolder, inputValues)
	if err != nil {
		return "", err
	}

	return partsHolder.encodeToHex(), nil
}

func (s *serializer) doSerialize(partsHolder *partsHolder, inputValues []interface{}) error {
	var err error

	for i, value := range inputValues {
		if value == nil {
			return errors.New("cannot serialize nil value")
		}

		switch value.(type) {
		case InputMultiValue:
			err = s.serializeInputMultiValue(partsHolder, value.(InputMultiValue))
		case InputVariadicValues:
			if i != len(inputValues)-1 {
				return errors.New("variadic values must be last among input values")
			}

			err = s.serializeInputVariadicValues(partsHolder, value.(InputVariadicValues))
		default:
			partsHolder.appendEmptyPart()
			err = s.serializeDirectlyEncodableValue(partsHolder, value)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *serializer) Deserialize(data string, outputValues []interface{}) error {
	partsHolder, err := newPartsHolderFromHex(data)
	if err != nil {
		return err
	}

	err = s.doDeserialize(partsHolder, outputValues)
	if err != nil {
		return err
	}

	return nil
}

func (s *serializer) doDeserialize(partsHolder *partsHolder, outputValues []interface{}) error {
	var err error

	for i, value := range outputValues {
		if value == nil {
			return errors.New("cannot deserialize into nil value")
		}

		switch value.(type) {
		case *OutputMultiValue:
			err = s.deserializeOutputMultiValue(partsHolder, value.(*OutputMultiValue))
		case *OutputVariadicValues:
			if i != len(outputValues)-1 {
				return errors.New("variadic values must be last among output values")
			}

			err = s.deserializeOutputVariadicValues(partsHolder, value.(*OutputVariadicValues))
		default:
			err = s.deserializeDirectlyEncodableValue(partsHolder, value)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *serializer) serializeInputMultiValue(partsHolder *partsHolder, value InputMultiValue) error {
	for _, item := range value.Items {
		err := s.doSerialize(partsHolder, []interface{}{item})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *serializer) serializeInputVariadicValues(partsHolder *partsHolder, value InputVariadicValues) error {
	for _, item := range value.Items {
		err := s.doSerialize(partsHolder, []interface{}{item})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *serializer) serializeDirectlyEncodableValue(partsHolder *partsHolder, value interface{}) error {
	data, err := s.codec.EncodeTopLevel(value)
	if err != nil {
		return err
	}

	return partsHolder.appendToLastPart(data)
}

func (s *serializer) deserializeOutputMultiValue(partsHolder *partsHolder, value *OutputMultiValue) error {
	for _, item := range value.Items {
		err := s.doDeserialize(partsHolder, []interface{}{item})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *serializer) deserializeOutputVariadicValues(partsHolder *partsHolder, value *OutputVariadicValues) error {
	if value.ItemCreator == nil {
		return errors.New("cannot deserialize variadic values: item creator is nil")
	}

	for !partsHolder.isFocusedBeyondLastPart() {
		newItem := value.ItemCreator()

		err := s.doDeserialize(partsHolder, []interface{}{newItem})
		if err != nil {
			return err
		}

		value.Items = append(value.Items, newItem)
	}

	return nil
}

func (s *serializer) deserializeDirectlyEncodableValue(partsHolder *partsHolder, value interface{}) error {
	part, err := partsHolder.readWholePart()
	if err != nil {
		return err
	}

	err = s.codec.DecodeTopLevel(part, value)
	if err != nil {
		return err
	}

	err = partsHolder.focusOnNextPart()
	if err != nil {
		return err
	}

	return nil
}
