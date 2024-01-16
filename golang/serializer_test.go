package abi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSerializer_Serialize(t *testing.T) {
	serializer := NewSerializer(NewDefaultCodec())

	t.Run("u8", func(t *testing.T) {
		writer := NewDefaultDataWriter()
		err := serializer.Serialize(writer, []interface{}{
			U8Value{Value: 0x42},
		})

		require.NoError(t, err)
		require.Equal(t, "42", writer.String())
	})

	t.Run("u16", func(t *testing.T) {
		writer := NewDefaultDataWriter()
		err := serializer.Serialize(writer, []interface{}{
			U16Value{Value: 0x4243},
		})

		require.NoError(t, err)
		require.Equal(t, "4243", writer.String())
	})

	t.Run("u8, u16", func(t *testing.T) {
		writer := NewDefaultDataWriter()

		err := serializer.Serialize(writer, []interface{}{
			U8Value{Value: 0x42},
			U16Value{Value: 0x4243},
		})

		require.NoError(t, err)
		require.Equal(t, "42@4243", writer.String())
	})

	t.Run("multi<u8, u16, u32>", func(t *testing.T) {
		writer := NewDefaultDataWriter()

		err := serializer.Serialize(writer, []interface{}{
			InputMultiValue{
				Items: []interface{}{
					U8Value{Value: 0x42},
					U16Value{Value: 0x4243},
					U32Value{Value: 0x42434445},
				},
			},
		})

		require.NoError(t, err)
		require.Equal(t, "42@4243@42434445", writer.String())
	})

	t.Run("u8, multi<u8, u16, u32>", func(t *testing.T) {
		writer := NewDefaultDataWriter()

		err := serializer.Serialize(writer, []interface{}{
			U8Value{Value: 0x42},
			InputMultiValue{
				Items: []interface{}{
					U8Value{Value: 0x42},
					U16Value{Value: 0x4243},
					U32Value{Value: 0x42434445},
				},
			},
		})

		require.NoError(t, err)
		require.Equal(t, "42@42@4243@42434445", writer.String())
	})

	t.Run("multi<multi<u8, u16>, multi<u8, u16>>", func(t *testing.T) {
		writer := NewDefaultDataWriter()

		err := serializer.Serialize(writer, []interface{}{
			InputMultiValue{
				Items: []interface{}{
					InputMultiValue{
						Items: []interface{}{
							U8Value{Value: 0x42},
							U16Value{Value: 0x4243},
						},
					},
					InputMultiValue{
						Items: []interface{}{
							U8Value{Value: 0x44},
							U16Value{Value: 0x4445},
						},
					},
				},
			},
		})

		require.NoError(t, err)
		require.Equal(t, "42@4243@44@4445", writer.String())
	})

	t.Run("variadic, of different types", func(t *testing.T) {
		writer := NewDefaultDataWriter()

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
		require.Equal(t, "42@4243", writer.String())
	})

	t.Run("variadic<u8>, u8: should err because variadic must be last", func(t *testing.T) {
		writer := NewDefaultDataWriter()

		err := serializer.Serialize(writer, []interface{}{
			InputVariadicValues{
				Items: []interface{}{
					U8Value{Value: 0x42},
					U8Value{Value: 0x43},
				},
			},
			U8Value{Value: 0x44},
		})

		require.ErrorIs(t, err, errVariadicMustBeLast)
	})

	t.Run("u8, variadic<u8>", func(t *testing.T) {
		writer := NewDefaultDataWriter()

		err := serializer.Serialize(writer, []interface{}{
			U8Value{Value: 0x41},
			InputVariadicValues{
				Items: []interface{}{
					U8Value{Value: 0x42},
					U8Value{Value: 0x43},
				},
			},
		})

		require.Nil(t, err)
		require.Equal(t, "41@42@43", writer.String())
	})
}

func TestSerializer_Deserialize(t *testing.T) {
	serializer := NewSerializer(NewDefaultCodec())

	t.Run("nil destination", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("")

		err := serializer.Deserialize(reader, []interface{}{nil})
		require.ErrorIs(t, err, errNilOutputValue)
	})

	t.Run("u8", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("42")

		outputValues := []interface{}{
			&U8Value{},
		}

		err := serializer.Deserialize(reader, outputValues)

		require.Nil(t, err)
		require.Equal(t, []interface{}{
			&U8Value{Value: 0x42},
		}, outputValues)
	})

	t.Run("u16", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("4243")

		outputValues := []interface{}{
			&U16Value{},
		}

		err := serializer.Deserialize(reader, outputValues)

		require.Nil(t, err)
		require.Equal(t, []interface{}{
			&U16Value{Value: 0x4243},
		}, outputValues)
	})

	t.Run("u8, u16", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("42@4243")

		outputValues := []interface{}{
			&U8Value{},
			&U16Value{},
		}

		err := serializer.Deserialize(reader, outputValues)

		require.Nil(t, err)
		require.Equal(t, []interface{}{
			&U8Value{Value: 0x42},
			&U16Value{Value: 0x4243},
		}, outputValues)
	})

	t.Run("multi<u8, u16, u32>", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("42@4243@42434445")

		outputValues := []interface{}{
			&OutputMultiValue{
				Items: []interface{}{
					&U8Value{},
					&U16Value{},
					&U32Value{},
				},
			},
		}

		err := serializer.Deserialize(reader, outputValues)

		require.Nil(t, err)
		require.Equal(t, []interface{}{
			&OutputMultiValue{
				Items: []interface{}{
					&U8Value{Value: 0x42},
					&U16Value{Value: 0x4243},
					&U32Value{Value: 0x42434445},
				},
			},
		}, outputValues)
	})

	t.Run("u8, multi<u8, u16, u32>", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("42@42@4243@42434445")

		outputValues := []interface{}{
			&U8Value{},
			&OutputMultiValue{
				Items: []interface{}{
					&U8Value{},
					&U16Value{},
					&U32Value{},
				},
			},
		}

		err := serializer.Deserialize(reader, outputValues)

		require.Nil(t, err)
		require.Equal(t, []interface{}{
			&U8Value{Value: 0x42},
			&OutputMultiValue{
				Items: []interface{}{
					&U8Value{Value: 0x42},
					&U16Value{Value: 0x4243},
					&U32Value{Value: 0x42434445},
				},
			},
		}, outputValues)
	})

	t.Run("variadic, should err because of nil item creator", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("")
		destination := &OutputVariadicValues{
			Items: []interface{}{},
		}

		err := serializer.Deserialize(reader, []interface{}{destination})
		require.ErrorIs(t, err, errNilItemCreator)
	})

	t.Run("empty: u8", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("")

		destination := &OutputVariadicValues{
			Items:       []interface{}{},
			ItemCreator: func() interface{} { return &U8Value{} },
		}

		err := serializer.Deserialize(reader, []interface{}{destination})
		require.NoError(t, err)
		require.Equal(t, []interface{}{&U8Value{Value: 0}}, destination.Items)
	})

	t.Run("variadic<u8>", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("2A@2B@2C")

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

	t.Run("varidic<u8>, with empty items", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("@01@")

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

	t.Run("varidic<u32>", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("AABBCCDD@DDCCBBAA")

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

	t.Run("varidic<u8>, should err because decoded value is too large", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("0100")

		destination := &OutputVariadicValues{
			Items:       []interface{}{},
			ItemCreator: func() interface{} { return &U8Value{} },
		}

		err := serializer.Deserialize(reader, []interface{}{destination})
		require.ErrorContains(t, err, "cannot decode u8, because of: decoded value is too large: 256 > 255")
	})
}
