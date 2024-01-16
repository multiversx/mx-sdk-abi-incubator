package abi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCodec_EncodeNested(t *testing.T) {
	codec := NewDefaultCodec()

	t.Run("u8", func(t *testing.T) {
		writer := NewDefaultDataWriter()
		writer.GotoNextPart()

		err := codec.EncodeNested(writer, U8Value{Value: 0x01})
		require.NoError(t, err)
		require.Equal(t, "01", writer.String())
	})

	t.Run("u16", func(t *testing.T) {
		writer := NewDefaultDataWriter()
		writer.GotoNextPart()

		err := codec.EncodeNested(writer, U16Value{Value: 0x4142})
		require.NoError(t, err)
		require.Equal(t, "4142", writer.String())
	})

	t.Run("u32", func(t *testing.T) {
		writer := NewDefaultDataWriter()
		writer.GotoNextPart()

		err := codec.EncodeNested(writer, U32Value{Value: 0x41424344})
		require.NoError(t, err)
		require.Equal(t, "41424344", writer.String())
	})

	t.Run("u64", func(t *testing.T) {
		writer := NewDefaultDataWriter()
		writer.GotoNextPart()

		err := codec.EncodeNested(writer, U64Value{Value: 0x4142434445464748})
		require.NoError(t, err)
		require.Equal(t, "4142434445464748", writer.String())
	})

	t.Run("struct", func(t *testing.T) {
		writer := NewDefaultDataWriter()
		writer.GotoNextPart()

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

		err := codec.EncodeNested(writer, fooStruct)
		require.NoError(t, err)
		require.Equal(t, "014142", writer.String())
	})

	t.Run("enum (discriminant == 0)", func(t *testing.T) {
		writer := NewDefaultDataWriter()
		writer.GotoNextPart()

		fooEnum := EnumValue{
			Discriminant: 0,
		}

		err := codec.EncodeNested(writer, fooEnum)
		require.NoError(t, err)
		require.Equal(t, "00", writer.String())
	})

	t.Run("enum (discriminant != 0)", func(t *testing.T) {
		writer := NewDefaultDataWriter()
		writer.GotoNextPart()

		fooEnum := EnumValue{
			Discriminant: 42,
		}

		err := codec.EncodeNested(writer, fooEnum)
		require.NoError(t, err)
		require.Equal(t, "2a", writer.String())
	})

	t.Run("enum with fields", func(t *testing.T) {
		writer := NewDefaultDataWriter()
		writer.GotoNextPart()

		fooEnum := EnumValue{
			Discriminant: 42,
			Fields: []Field{
				{
					Value: U8Value{Value: 0x01},
				},
				{
					Value: U16Value{Value: 0x4142},
				},
			},
		}

		err := codec.EncodeNested(writer, fooEnum)
		require.NoError(t, err)
		require.Equal(t, "2a014142", writer.String())
	})
}

func TestCodec_EncodeTopLevel(t *testing.T) {
	codec := NewDefaultCodec()

	t.Run("u8", func(t *testing.T) {
		writer := NewDefaultDataWriter()
		writer.GotoNextPart()

		err := codec.EncodeTopLevel(writer, U8Value{Value: 0x01})
		require.NoError(t, err)
		require.Equal(t, "01", writer.String())
	})

	t.Run("u8 (zero)", func(t *testing.T) {
		writer := NewDefaultDataWriter()
		writer.GotoNextPart()

		err := codec.EncodeTopLevel(writer, U8Value{Value: 0})
		require.NoError(t, err)
		require.Equal(t, "", writer.String())
	})

	t.Run("u16", func(t *testing.T) {
		writer := NewDefaultDataWriter()
		writer.GotoNextPart()

		err := codec.EncodeTopLevel(writer, U16Value{Value: 0x0042})
		require.NoError(t, err)
		require.Equal(t, "42", writer.String())
	})

	t.Run("u32", func(t *testing.T) {
		writer := NewDefaultDataWriter()
		writer.GotoNextPart()

		err := codec.EncodeTopLevel(writer, U32Value{Value: 0x00004242})
		require.NoError(t, err)
		require.Equal(t, "4242", writer.String())
	})

	t.Run("u64", func(t *testing.T) {
		writer := NewDefaultDataWriter()
		writer.GotoNextPart()

		err := codec.EncodeTopLevel(writer, U64Value{Value: 0x0042434445464748})
		require.NoError(t, err)
		require.Equal(t, "42434445464748", writer.String())
	})

	t.Run("struct", func(t *testing.T) {
		writer := NewDefaultDataWriter()
		writer.GotoNextPart()

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

		err := codec.EncodeTopLevel(writer, fooStruct)
		require.NoError(t, err)
		require.Equal(t, "014142", writer.String())
	})

	t.Run("enum (discriminant == 0)", func(t *testing.T) {
		writer := NewDefaultDataWriter()
		writer.GotoNextPart()

		fooEnum := EnumValue{
			Discriminant: 0,
		}

		err := codec.EncodeTopLevel(writer, fooEnum)
		require.NoError(t, err)
		require.Equal(t, "", writer.String())
	})

	t.Run("enum (discriminant != 0)", func(t *testing.T) {
		writer := NewDefaultDataWriter()
		writer.GotoNextPart()

		fooEnum := EnumValue{
			Discriminant: 42,
		}

		err := codec.EncodeTopLevel(writer, fooEnum)
		require.NoError(t, err)
		require.Equal(t, "2a", writer.String())
	})

	t.Run("enum with fields", func(t *testing.T) {
		writer := NewDefaultDataWriter()
		writer.GotoNextPart()

		fooEnum := EnumValue{
			Discriminant: 42,
			Fields: []Field{
				{
					Value: U8Value{Value: 0x01},
				},
				{
					Value: U16Value{Value: 0x4142},
				},
			},
		}

		err := codec.EncodeTopLevel(writer, fooEnum)
		require.NoError(t, err)
		require.Equal(t, "2a014142", writer.String())
	})
}

func TestCodec_DecodeNested(t *testing.T) {
	codec := NewDefaultCodec()

	t.Run("u8", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("01")
		destination := &U8Value{}

		err := codec.DecodeNested(reader, destination)
		require.NoError(t, err)
		require.Equal(t, &U8Value{Value: 0x01}, destination)
	})

	t.Run("u16", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("4142")
		destination := &U16Value{}

		err := codec.DecodeNested(reader, destination)
		require.NoError(t, err)
		require.Equal(t, &U16Value{Value: 0x4142}, destination)
	})

	t.Run("u32", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("41424344")
		destination := &U32Value{}

		err := codec.DecodeNested(reader, destination)
		require.NoError(t, err)
		require.Equal(t, &U32Value{Value: 0x41424344}, destination)
	})

	t.Run("u64", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("4142434445464748")
		destination := &U64Value{}

		err := codec.DecodeNested(reader, destination)
		require.NoError(t, err)
		require.Equal(t, &U64Value{Value: 0x4142434445464748}, destination)
	})

	t.Run("u16, should err because it cannot read 2 bytes", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("01")
		destination := &U16Value{}

		err := codec.DecodeNested(reader, destination)
		require.ErrorContains(t, err, "cannot read 2 bytes")
	})

	t.Run("u32, should err because it cannot read 4 bytes", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("4142")
		destination := &U32Value{}

		err := codec.DecodeNested(reader, destination)
		require.ErrorContains(t, err, "cannot read 4 bytes")
	})

	t.Run("u64, should err because it cannot read 8 bytes", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("41424344")
		destination := &U64Value{}

		err := codec.DecodeNested(reader, destination)
		require.ErrorContains(t, err, "cannot read 8 bytes")
	})

	t.Run("struct", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("014142")

		destination := &StructValue{
			Fields: []Field{
				{
					Value: &U8Value{},
				},
				{
					Value: &U16Value{},
				},
			},
		}

		err := codec.DecodeNested(reader, destination)
		require.NoError(t, err)
		require.Equal(t, &StructValue{
			Fields: []Field{
				{
					Value: &U8Value{Value: 0x01},
				},
				{
					Value: &U16Value{Value: 0x4142},
				},
			},
		}, destination)
	})

	t.Run("enum (discriminant == 0)", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("00")
		destination := &EnumValue{}

		err := codec.DecodeNested(reader, destination)
		require.NoError(t, err)
		require.Equal(t, &EnumValue{
			Discriminant: 0x00,
		}, destination)
	})

	t.Run("enum (discriminant != 0)", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("01")
		destination := &EnumValue{}

		err := codec.DecodeNested(reader, destination)
		require.NoError(t, err)
		require.Equal(t, &EnumValue{
			Discriminant: 0x01,
		}, destination)
	})

	t.Run("enum with fields", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("01014142")

		destination := &EnumValue{
			Fields: []Field{
				{
					Value: &U8Value{},
				},
				{
					Value: &U16Value{},
				},
			},
		}

		err := codec.DecodeNested(reader, destination)
		require.NoError(t, err)
		require.Equal(t, &EnumValue{
			Discriminant: 0x01,
			Fields: []Field{
				{
					Value: &U8Value{Value: 0x01},
				},
				{
					Value: &U16Value{Value: 0x4142},
				},
			},
		}, destination)
	})
}

func TestCodec_DecodeTopLevel(t *testing.T) {
	codec := NewDefaultCodec()

	t.Run("u8", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("01")
		destination := &U8Value{}

		err := codec.DecodeTopLevel(reader, destination)
		require.NoError(t, err)
		require.Equal(t, &U8Value{Value: 0x01}, destination)
	})

	t.Run("u16", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("02")
		destination := &U16Value{}

		err := codec.DecodeTopLevel(reader, destination)
		require.NoError(t, err)
		require.Equal(t, &U16Value{Value: 0x0002}, destination)
	})

	t.Run("u32", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("03")
		destination := &U32Value{}

		err := codec.DecodeTopLevel(reader, destination)
		require.NoError(t, err)
		require.Equal(t, &U32Value{Value: 0x00000003}, destination)
	})

	t.Run("u64", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("04")
		destination := &U64Value{}

		err := codec.DecodeTopLevel(reader, destination)
		require.NoError(t, err)
		require.Equal(t, &U64Value{Value: 0x0000000000000004}, destination)
	})

	t.Run("u8, should err because decoded value is too large", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("4142")
		destination := &U8Value{}

		err := codec.DecodeTopLevel(reader, destination)
		require.ErrorContains(t, err, "decoded value is too large")
	})

	t.Run("u16, should err because decoded value is too large", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("41424344")
		destination := &U16Value{}

		err := codec.DecodeTopLevel(reader, destination)
		require.ErrorContains(t, err, "decoded value is too large")
	})

	t.Run("u32, should err because decoded value is too large", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("4142434445464748")
		destination := &U32Value{}

		err := codec.DecodeTopLevel(reader, destination)
		require.ErrorContains(t, err, "decoded value is too large")
	})

	t.Run("u64, should err because decoded value is too large", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("41424344454647489876")
		destination := &U64Value{}

		err := codec.DecodeTopLevel(reader, destination)
		require.ErrorContains(t, err, "decoded value is too large")
	})

	t.Run("struct", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("014142")

		destination := &StructValue{
			Fields: []Field{
				{
					Value: &U8Value{},
				},
				{
					Value: &U16Value{},
				},
			},
		}

		err := codec.DecodeTopLevel(reader, destination)
		require.NoError(t, err)
		require.Equal(t, &StructValue{
			Fields: []Field{
				{
					Value: &U8Value{Value: 0x01},
				},
				{
					Value: &U16Value{Value: 0x4142},
				},
			},
		}, destination)
	})

	t.Run("enum (discriminant == 0)", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("")
		destination := &EnumValue{}

		err := codec.DecodeTopLevel(reader, destination)
		require.NoError(t, err)
		require.Equal(t, &EnumValue{
			Discriminant: 0x00,
		}, destination)
	})

	t.Run("enum (discriminant != 0)", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("01")
		destination := &EnumValue{}

		err := codec.DecodeTopLevel(reader, destination)
		require.NoError(t, err)
		require.Equal(t, &EnumValue{
			Discriminant: 0x01,
		}, destination)
	})

	t.Run("enum with fields", func(t *testing.T) {
		reader, _ := NewDefaultDataReaderFromString("01014142")

		destination := &EnumValue{
			Fields: []Field{
				{
					Value: &U8Value{},
				},
				{
					Value: &U16Value{},
				},
			},
		}

		err := codec.DecodeTopLevel(reader, destination)
		require.NoError(t, err)
		require.Equal(t, &EnumValue{
			Discriminant: 0x01,
			Fields: []Field{
				{
					Value: &U8Value{Value: 0x01},
				},
				{
					Value: &U16Value{Value: 0x4142},
				},
			},
		}, destination)
	})
}
