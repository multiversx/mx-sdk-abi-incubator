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

func TestSerializer_SerializeCompositeValue(t *testing.T) {
	testSerialize(t, []interface{}{
		CompositeValue{
			Items: []interface{}{
				U8Value{Value: 0x42},
				U16Value{Value: 0x4243},
				U32Value{Value: 0x42434445},
			},
		},
	}, "42@4243@42434445")

	testSerialize(t, []interface{}{
		U8Value{Value: 0x42},
		CompositeValue{
			Items: []interface{}{
				U8Value{Value: 0x42},
				U16Value{Value: 0x4243},
				U32Value{Value: 0x42434445},
			},
		},
	}, "42@42@4243@42434445")
}

func TestSerializer_SerializeInputVariadicValues(t *testing.T) {
	testSerialize(t, []interface{}{
		InputVariadicValues{
			Items: []interface{}{},
		},
	}, "")

	testSerialize(t, []interface{}{
		InputVariadicValues{
			Items: []interface{}{
				U8Value{Value: 0x42},
				U8Value{Value: 0x43},
				U8Value{Value: 0x44},
			},
		},
	}, "42@43@44")

	testSerialize(t, []interface{}{
		InputVariadicValues{
			Items: []interface{}{
				CompositeValue{
					Items: []interface{}{
						U8Value{Value: 0x42},
						U16Value{Value: 0x4243},
					},
				},
				CompositeValue{
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
	writer := NewDefaultDataWriter()

	t.Run("variadic items of different types (1)", func(t *testing.T) {
		err := serializer.Serialize(writer, []interface{}{
			InputVariadicValues{
				Items: []interface{}{
					U8Value{Value: 0x42},
					U16Value{Value: 0x4243},
				},
			},
		})
		// For now, the serializer does not perform such a strict type check.
		// Although doable, it would be slightly complex and, if done, might be even dropped in the future
		// (with respect to the decoder that is embedded in Rust-based smart contracts).
		require.Nil(t, err)
	})

	t.Run("variadic items of different types (2)", func(t *testing.T) {
		err := serializer.Serialize(writer, []interface{}{
			InputVariadicValues{
				Items: []interface{}{
					CompositeValue{
						Items: []interface{}{
							U8Value{Value: 0x42},
						},
					},
					CompositeValue{
						Items: []interface{}{
							U16Value{Value: 0x43},
						},
					},
				},
			},
		})
		// For now, the serializer does not perform such a strict type check.
		// Although doable, it would be slightly complex and, if done, might be even dropped in the future
		// (with respect to the decoder that is embedded in Rust-based smart contracts).
		require.Nil(t, err)
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

func TestSerializer_DeserializeCompositeValue(t *testing.T) {
	testDeserialize(t, "42@4243@42434445",
		[]interface{}{
			&CompositeValue{
				Items: []interface{}{
					&U8Value{},
					&U16Value{},
					&U32Value{},
				},
			},
		},
		[]interface{}{
			&CompositeValue{
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
			&CompositeValue{
				Items: []interface{}{
					&U8Value{},
					&U16Value{},
					&U32Value{},
				},
			},
		},
		[]interface{}{
			&U8Value{Value: 0x42},
			&CompositeValue{
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
