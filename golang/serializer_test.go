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
