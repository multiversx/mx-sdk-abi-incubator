package abi

type U8Value struct {
	Value uint8
}

type U16Value struct {
	Value uint16
}

type U32Value struct {
	Value uint32
}

type U64Value struct {
	Value uint64
}

type I8Value struct {
	Value int8
}

type I16Value struct {
	Value int16
}

type I32Value struct {
	Value int32
}

type I64Value struct {
	Value int64
}

type BytesValue struct {
	Value []byte
}

type StringValue struct {
	Value string
}

type BoolValue struct {
	Value bool
}

type OptionValue struct {
	Value interface{}
}

type Field struct {
	Name  string
	Value interface{}
}

type StructValue struct {
	Fields []Field
}

type TupleValue struct {
	Fields []Field
}

type InputListValue struct {
	Items []interface{}
}

type OutputListValue struct {
	Items       []interface{}
	ItemCreator func() interface{}
}

type MultiValue struct {
	Items []interface{}
}

type OutputVariadicValues struct {
	Items       []interface{}
	ItemCreator func() interface{}
}

type OptionalValue struct {
	Value interface{}
	IsSet bool
}
