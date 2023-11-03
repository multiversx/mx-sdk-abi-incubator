package abi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCodec_EncodeNested_Primitives(t *testing.T) {
	testEncodeNested(t, U8Value{Value: 0x01}, "01")
	testEncodeNested(t, U16Value{Value: 0x4142}, "4142")
	testEncodeNested(t, U32Value{Value: 0x41424344}, "41424344")
	testEncodeNested(t, U64Value{Value: 0x4142434445464748}, "4142434445464748")
}

func TestCodec_EncodeTopLevel_Primitives(t *testing.T) {
	testEncodeTopLevel(t, U8Value{Value: 0x01}, "01")
	testEncodeTopLevel(t, U16Value{Value: 0x0042}, "42")
	testEncodeTopLevel(t, U32Value{Value: 0x00004242}, "4242")
	testEncodeTopLevel(t, U64Value{Value: 0x0042434445464748}, "42434445464748")
}

func TestCodec_DecodeNested_Primitives(t *testing.T) {
	testDecodeNested(t, "01", &U8Value{}, &U8Value{Value: 0x01})
	testDecodeNested(t, "4142", &U16Value{}, &U16Value{Value: 0x4142})
	testDecodeNested(t, "41424344", &U32Value{}, &U32Value{Value: 0x41424344})
	testDecodeNested(t, "4142434445464748", &U64Value{}, &U64Value{Value: 0x4142434445464748})
}

func TestCodec_EncodeNested_Structures(t *testing.T) {
	fooStruct := StructValue{
		Fields: []Field{
			{
				Value: U8Value{Value: 0x01},
			},
			{
				Value: U16Value{Value: 0x4142},
			},
		},
	}

	testEncodeNested(t, fooStruct, "014142")
}

func testEncodeNested(t *testing.T, value interface{}, expected string) {
	codec := NewDefaultCodec()
	writer := NewDefaultDataWriter()
	writer.GotoNextPart()

	err := codec.EncodeNested(writer, value)
	require.NoError(t, err)
	require.Equal(t, expected, writer.String())
}

func testEncodeTopLevel(t *testing.T, value interface{}, expected string) {
	codec := NewDefaultCodec()
	writer := NewDefaultDataWriter()
	writer.GotoNextPart()

	err := codec.EncodeTopLevel(writer, value)
	require.NoError(t, err)
	require.Equal(t, expected, writer.String())
}

func testDecodeNested(t *testing.T, encoded string, destination interface{}, expected interface{}) {
	codec := NewDefaultCodec()
	reader, err := NewDefaultDataReaderFromString(encoded)
	require.NoError(t, err)

	err = codec.DecodeNested(reader, destination)
	require.NoError(t, err)
	require.Equal(t, expected, destination)
}
