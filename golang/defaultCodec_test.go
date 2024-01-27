package abi

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCodec_EncodeNested(t *testing.T) {
	doTest := func(t *testing.T, value interface{}, expected string) {
		codec := NewDefaultCodec()
		encoded, err := codec.EncodeNested(value)
		require.NoError(t, err)
		require.Equal(t, expected, hex.EncodeToString(encoded))
	}

	t.Run("u8", func(t *testing.T) {
		doTest(t, U8Value{Value: 0x00}, "00")
		doTest(t, U8Value{Value: 0x01}, "01")
		doTest(t, U8Value{Value: 0x42}, "42")
		doTest(t, U8Value{Value: 0xff}, "ff")
	})

	t.Run("u16", func(t *testing.T) {
		doTest(t, U16Value{Value: 0x00}, "0000")
		doTest(t, U16Value{Value: 0x11}, "0011")
		doTest(t, U16Value{Value: 0x1234}, "1234")
		doTest(t, U16Value{Value: 0xffff}, "ffff")
	})

	t.Run("u32", func(t *testing.T) {
		doTest(t, U32Value{Value: 0x00000000}, "00000000")
		doTest(t, U32Value{Value: 0x00000011}, "00000011")
		doTest(t, U32Value{Value: 0x00001122}, "00001122")
		doTest(t, U32Value{Value: 0x00112233}, "00112233")
		doTest(t, U32Value{Value: 0x11223344}, "11223344")
		doTest(t, U32Value{Value: 0xffffffff}, "ffffffff")
	})

	t.Run("u64", func(t *testing.T) {
		doTest(t, U64Value{Value: 0x0000000000000000}, "0000000000000000")
		doTest(t, U64Value{Value: 0x0000000000000011}, "0000000000000011")
		doTest(t, U64Value{Value: 0x0000000000001122}, "0000000000001122")
		doTest(t, U64Value{Value: 0x0000000000112233}, "0000000000112233")
		doTest(t, U64Value{Value: 0x0000000011223344}, "0000000011223344")
		doTest(t, U64Value{Value: 0x0000001122334455}, "0000001122334455")
		doTest(t, U64Value{Value: 0x0000112233445566}, "0000112233445566")
		doTest(t, U64Value{Value: 0x0011223344556677}, "0011223344556677")
		doTest(t, U64Value{Value: 0x1122334455667788}, "1122334455667788")
		doTest(t, U64Value{Value: 0xffffffffffffffff}, "ffffffffffffffff")
	})

	t.Run("struct", func(t *testing.T) {
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

		doTest(t, fooStruct, "014142")
	})

	t.Run("enum (discriminant == 0)", func(t *testing.T) {
		fooEnum := EnumValue{
			Discriminant: 0,
		}

		doTest(t, fooEnum, "00")
	})

	t.Run("enum (discriminant != 0)", func(t *testing.T) {
		fooEnum := EnumValue{
			Discriminant: 42,
		}

		doTest(t, fooEnum, "2a")
	})

	t.Run("enum with fields", func(t *testing.T) {
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

		doTest(t, fooEnum, "2a014142")
	})
}

func TestCodec_EncodeTopLevel(t *testing.T) {
	doTest := func(t *testing.T, value interface{}, expected string) {
		codec := NewDefaultCodec()
		encoded, err := codec.EncodeTopLevel(value)
		require.NoError(t, err)
		require.Equal(t, expected, hex.EncodeToString(encoded))
	}

	t.Run("u8", func(t *testing.T) {
		doTest(t, U8Value{Value: 0x00}, "")
		doTest(t, U8Value{Value: 0x01}, "01")
	})

	t.Run("u16", func(t *testing.T) {
		doTest(t, U16Value{Value: 0x0042}, "42")
	})

	t.Run("u32", func(t *testing.T) {
		doTest(t, U32Value{Value: 0x00004242}, "4242")
	})

	t.Run("u64", func(t *testing.T) {
		doTest(t, U64Value{Value: 0x0042434445464748}, "42434445464748")
	})

	t.Run("struct", func(t *testing.T) {
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

		doTest(t, fooStruct, "014142")
	})

	t.Run("enum (discriminant == 0)", func(t *testing.T) {
		fooEnum := EnumValue{
			Discriminant: 0,
		}

		doTest(t, fooEnum, "")
	})

	t.Run("enum (discriminant != 0)", func(t *testing.T) {
		fooEnum := EnumValue{
			Discriminant: 42,
		}

		doTest(t, fooEnum, "2a")
	})

	t.Run("enum with fields", func(t *testing.T) {
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

		doTest(t, fooEnum, "2a014142")
	})
}

func TestCodec_DecodeNested(t *testing.T) {
	codec := NewDefaultCodec()

	t.Run("u8", func(t *testing.T) {
		data, _ := hex.DecodeString("01")
		destination := &U8Value{}

		err := codec.DecodeNested(data, destination)
		require.NoError(t, err)
		require.Equal(t, &U8Value{Value: 0x01}, destination)
	})

	t.Run("u16", func(t *testing.T) {
		data, _ := hex.DecodeString("4142")
		destination := &U16Value{}

		err := codec.DecodeNested(data, destination)
		require.NoError(t, err)
		require.Equal(t, &U16Value{Value: 0x4142}, destination)
	})

	t.Run("u32", func(t *testing.T) {
		data, _ := hex.DecodeString("41424344")
		destination := &U32Value{}

		err := codec.DecodeNested(data, destination)
		require.NoError(t, err)
		require.Equal(t, &U32Value{Value: 0x41424344}, destination)
	})

	t.Run("u64", func(t *testing.T) {
		data, _ := hex.DecodeString("4142434445464748")
		destination := &U64Value{}

		err := codec.DecodeNested(data, destination)
		require.NoError(t, err)
		require.Equal(t, &U64Value{Value: 0x4142434445464748}, destination)
	})

	t.Run("u16, should err because it cannot read 2 bytes", func(t *testing.T) {
		data, _ := hex.DecodeString("01")
		destination := &U16Value{}

		err := codec.DecodeNested(data, destination)
		require.ErrorContains(t, err, "cannot read 2 bytes")
	})

	t.Run("u32, should err because it cannot read 4 bytes", func(t *testing.T) {
		data, _ := hex.DecodeString("4142")
		destination := &U32Value{}

		err := codec.DecodeNested(data, destination)
		require.ErrorContains(t, err, "cannot read 4 bytes")
	})

	t.Run("u64, should err because it cannot read 8 bytes", func(t *testing.T) {
		data, _ := hex.DecodeString("41424344")
		destination := &U64Value{}

		err := codec.DecodeNested(data, destination)
		require.ErrorContains(t, err, "cannot read 8 bytes")
	})

	t.Run("struct", func(t *testing.T) {
		data, _ := hex.DecodeString("014142")

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

		err := codec.DecodeNested(data, destination)
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
		data, _ := hex.DecodeString("00")
		destination := &EnumValue{}

		err := codec.DecodeNested(data, destination)
		require.NoError(t, err)
		require.Equal(t, &EnumValue{
			Discriminant: 0x00,
		}, destination)
	})

	t.Run("enum (discriminant != 0)", func(t *testing.T) {
		data, _ := hex.DecodeString("01")
		destination := &EnumValue{}

		err := codec.DecodeNested(data, destination)
		require.NoError(t, err)
		require.Equal(t, &EnumValue{
			Discriminant: 0x01,
		}, destination)
	})

	t.Run("enum with fields", func(t *testing.T) {
		data, _ := hex.DecodeString("01014142")

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

		err := codec.DecodeNested(data, destination)
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
		data, _ := hex.DecodeString("01")
		destination := &U8Value{}

		err := codec.DecodeTopLevel(data, destination)
		require.NoError(t, err)
		require.Equal(t, &U8Value{Value: 0x01}, destination)
	})

	t.Run("u16", func(t *testing.T) {
		data, _ := hex.DecodeString("02")
		destination := &U16Value{}

		err := codec.DecodeTopLevel(data, destination)
		require.NoError(t, err)
		require.Equal(t, &U16Value{Value: 0x0002}, destination)
	})

	t.Run("u32", func(t *testing.T) {
		data, _ := hex.DecodeString("03")
		destination := &U32Value{}

		err := codec.DecodeTopLevel(data, destination)
		require.NoError(t, err)
		require.Equal(t, &U32Value{Value: 0x00000003}, destination)
	})

	t.Run("u64", func(t *testing.T) {
		data, _ := hex.DecodeString("04")
		destination := &U64Value{}

		err := codec.DecodeTopLevel(data, destination)
		require.NoError(t, err)
		require.Equal(t, &U64Value{Value: 0x0000000000000004}, destination)
	})

	t.Run("u8, should err because decoded value is too large", func(t *testing.T) {
		data, _ := hex.DecodeString("4142")
		destination := &U8Value{}

		err := codec.DecodeTopLevel(data, destination)
		require.ErrorContains(t, err, "decoded value is too large")
	})

	t.Run("u16, should err because decoded value is too large", func(t *testing.T) {
		data, _ := hex.DecodeString("41424344")
		destination := &U16Value{}

		err := codec.DecodeTopLevel(data, destination)
		require.ErrorContains(t, err, "decoded value is too large")
	})

	t.Run("u32, should err because decoded value is too large", func(t *testing.T) {
		data, _ := hex.DecodeString("4142434445464748")
		destination := &U32Value{}

		err := codec.DecodeTopLevel(data, destination)
		require.ErrorContains(t, err, "decoded value is too large")
	})

	t.Run("u64, should err because decoded value is too large", func(t *testing.T) {
		data, _ := hex.DecodeString("41424344454647489876")
		destination := &U64Value{}

		err := codec.DecodeTopLevel(data, destination)
		require.ErrorContains(t, err, "decoded value is too large")
	})

	t.Run("struct", func(t *testing.T) {
		data, _ := hex.DecodeString("014142")

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

		err := codec.DecodeTopLevel(data, destination)
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
		data, _ := hex.DecodeString("")
		destination := &EnumValue{}

		err := codec.DecodeTopLevel(data, destination)
		require.NoError(t, err)
		require.Equal(t, &EnumValue{
			Discriminant: 0x00,
		}, destination)
	})

	t.Run("enum (discriminant != 0)", func(t *testing.T) {
		data, _ := hex.DecodeString("01")
		destination := &EnumValue{}

		err := codec.DecodeTopLevel(data, destination)
		require.NoError(t, err)
		require.Equal(t, &EnumValue{
			Discriminant: 0x01,
		}, destination)
	})

	t.Run("enum with fields", func(t *testing.T) {
		data, _ := hex.DecodeString("01014142")

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

		err := codec.DecodeTopLevel(data, destination)
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
