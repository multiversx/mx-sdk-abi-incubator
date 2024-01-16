package abi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSerializer_Serialize_DirectlyEncodableValues(t *testing.T) {
	testSerialize(t, []interface{}{
		U8Value{Value: 0x42},
	}, "42")

	testSerialize(t, []interface{}{
		U16Value{Value: 0x4243},
	}, "4243")

	testSerialize(t, []interface{}{
		U8Value{Value: 0x42},
		U16Value{Value: 0x4243},
	}, "42@4243")
}

func TestSerializer_SerializeMultiValue(t *testing.T) {
	testSerialize(t, []interface{}{
		MultiValue{
			Items: []interface{}{
				U8Value{Value: 0x42},
				U16Value{Value: 0x4243},
				U32Value{Value: 0x42434445},
			},
		},
	}, "42@4243@42434445")

	testSerialize(t, []interface{}{
		U8Value{Value: 0x42},
		MultiValue{
			Items: []interface{}{
				U8Value{Value: 0x42},
				U16Value{Value: 0x4243},
				U32Value{Value: 0x42434445},
			},
		},
	}, "42@42@4243@42434445")
}

func TestSerializer_SerializeMultiValues(t *testing.T) {
	testSerialize(t, []interface{}{
		MultiValue{
			Items: []interface{}{},
		},
	}, "")

	testSerialize(t, []interface{}{
		MultiValue{
			Items: []interface{}{
				U8Value{Value: 0x42},
				U8Value{Value: 0x43},
				U8Value{Value: 0x44},
			},
		},
	}, "42@43@44")

	testSerialize(t, []interface{}{
		MultiValue{
			Items: []interface{}{
				MultiValue{
					Items: []interface{}{
						U8Value{Value: 0x42},
						U16Value{Value: 0x4243},
					},
				},
				MultiValue{
					Items: []interface{}{
						U8Value{Value: 0x44},
						U16Value{Value: 0x4445},
					},
				},
			},
		},
	}, "42@4243@44@4445")
}

func TestSerializer_Serialize_WithErrors(t *testing.T) {
	serializer := NewSerializer(NewDefaultCodec())

	t.Run("multi-value items of different types (1)", func(t *testing.T) {
	writer := NewDefaultDataWriter()

		err := serializer.Serialize(writer, []interface{}{
			MultiValue{
				Items: []interface{}{
					U8Value{Value: 0x42},
					U16Value{Value: 0x4243},
				},
			},
		})

		require.Nil(t, err)
		require.Equal(t, "42@4243", writer.String())
	})

	t.Run("multi-values of different types (2)", func(t *testing.T) {
		writer := NewDefaultDataWriter()

		err := serializer.Serialize(writer, []interface{}{
			MultiValue{
				Items: []interface{}{
					MultiValue{
						Items: []interface{}{
							U8Value{Value: 0x42},
						},
					},
					MultiValue{
						Items: []interface{}{
							U16Value{Value: 0x43},
						},
					},
				},
			},
		})

		require.Nil(t, err)
		require.Equal(t, "42@43", writer.String())
	})
}

func TestSerializer_Deserialize_DirectlyEncodableValues(t *testing.T) {
	testDeserialize(t, "42",
		[]interface{}{
			&U8Value{},
		},
		[]interface{}{
			&U8Value{Value: 0x42},
		},
	)

	testDeserialize(t, "4243",
		[]interface{}{
			&U16Value{},
		},
		[]interface{}{
			&U16Value{Value: 0x4243},
		},
	)

	testDeserialize(t, "42@4243",
		[]interface{}{
			&U8Value{},
			&U16Value{},
		},
		[]interface{}{
			&U8Value{Value: 0x42},
			&U16Value{Value: 0x4243},
		},
	)
}

func TestSerializer_DeserializeMultiValue(t *testing.T) {
	testDeserialize(t, "42@4243@42434445",
		[]interface{}{
			&MultiValue{
				Items: []interface{}{
					&U8Value{},
					&U16Value{},
					&U32Value{},
				},
			},
		},
		[]interface{}{
			&MultiValue{
				Items: []interface{}{
					&U8Value{Value: 0x42},
					&U16Value{Value: 0x4243},
					&U32Value{Value: 0x42434445},
				},
			},
		},
	)

	testDeserialize(t, "42@42@4243@42434445",
		[]interface{}{
			&U8Value{},
			&MultiValue{
				Items: []interface{}{
					&U8Value{},
					&U16Value{},
					&U32Value{},
				},
			},
		},
		[]interface{}{
			&U8Value{Value: 0x42},
			&MultiValue{
				Items: []interface{}{
					&U8Value{Value: 0x42},
					&U16Value{Value: 0x4243},
					&U32Value{Value: 0x42434445},
				},
			},
		},
	)
}

func TestSerializer_DeserializeOutputVariadicValues(t *testing.T) {
	t.Run("nil destination", func(t *testing.T) {
		serializer, reader := setupDeserializeTest(t, "")

		err := serializer.Deserialize(reader, []interface{}{nil})
		require.ErrorIs(t, err, errNilOutputValue)
	})

	t.Run("nil item creator", func(t *testing.T) {
		serializer, reader := setupDeserializeTest(t, "")
		destination := &OutputVariadicValues{
			Items: []interface{}{},
		}

		err := serializer.Deserialize(reader, []interface{}{destination})
		require.ErrorIs(t, err, errNilItemCreator)
	})

	t.Run("empty", func(t *testing.T) {
		serializer, reader := setupDeserializeTest(t, "")
		destination := &OutputVariadicValues{
			Items:       []interface{}{},
			ItemCreator: func() interface{} { return struct{}{} },
		}

		err := serializer.Deserialize(reader, []interface{}{destination})
		require.NoError(t, err)
	})

	t.Run("variadic primitives (1)", func(t *testing.T) {
		serializer, reader := setupDeserializeTest(t, "2A@2B@2C")
		destination := &OutputVariadicValues{
			Items:       []interface{}{},
			ItemCreator: func() interface{} { return &U8Value{} },
		}

		err := serializer.Deserialize(reader, []interface{}{destination})
		require.NoError(t, err)

		require.Equal(t, []interface{}{
			&U8Value{Value: 42},
			&U8Value{Value: 43},
			&U8Value{Value: 44},
		}, destination.Items)
	})

	t.Run("variadic primitives (2)", func(t *testing.T) {
		serializer, reader := setupDeserializeTest(t, "@01@")
		destination := &OutputVariadicValues{
			Items:       []interface{}{},
			ItemCreator: func() interface{} { return &U8Value{} },
		}

		err := serializer.Deserialize(reader, []interface{}{destination})
		require.NoError(t, err)

		require.Equal(t, []interface{}{
			&U8Value{Value: 0},
			&U8Value{Value: 1},
			&U8Value{Value: 0},
		}, destination.Items)
	})

	t.Run("variadic primitives (3)", func(t *testing.T) {
		serializer, reader := setupDeserializeTest(t, "AABBCCDD@DDCCBBAA")
		destination := &OutputVariadicValues{
			Items:       []interface{}{},
			ItemCreator: func() interface{} { return &U32Value{} },
		}

		err := serializer.Deserialize(reader, []interface{}{destination})
		require.NoError(t, err)

		require.Equal(t, []interface{}{
			&U32Value{Value: 0xAABBCCDD},
			&U32Value{Value: 0xDDCCBBAA},
		}, destination.Items)
	})

	t.Run("variadic primitives (4)", func(t *testing.T) {
		serializer, reader := setupDeserializeTest(t, "0100")
		destination := &OutputVariadicValues{
			Items:       []interface{}{},
			ItemCreator: func() interface{} { return &U8Value{} },
		}

		err := serializer.Deserialize(reader, []interface{}{destination})
		require.ErrorContains(t, err, "cannot decode u8, because of: decoded value is too large: 256 > 255")
	})
}

func testSerialize(t *testing.T, values []interface{}, expected string) {
	serializer := NewSerializer(NewDefaultCodec())
	writer := NewDefaultDataWriter()

	err := serializer.Serialize(writer, values)
	require.NoError(t, err)
	require.Equal(t, expected, writer.String())
}

func testDeserialize(t *testing.T, encoded string, destination []interface{}, expected []interface{}) {
	serializer := NewSerializer(NewDefaultCodec())
	reader, err := NewDefaultDataReaderFromString(encoded)
	require.NoError(t, err)

	err = serializer.Deserialize(reader, destination)
	require.NoError(t, err)
	require.Equal(t, expected, destination)
}

func setupDeserializeTest(t *testing.T, serializedInput string) (*serializer, *defaultDataReader) {
	serializer := NewSerializer(NewDefaultCodec())
	reader, err := NewDefaultDataReaderFromString(serializedInput)
	require.NoError(t, err)

	return serializer, reader
}
