package abi

type dataReader interface {
	ReadWholePart() ([]byte, error)
	GotoNextPart() error
	IsEndOfData() bool
}

type dataWriter interface {
	Write(data []byte) (int, error)
	GotoNextPart()
}

type codec interface {
	EncodeNested(value interface{}) ([]byte, error)
	EncodeTopLevel(value interface{}) ([]byte, error)
	DecodeNested(data []byte, value interface{}) error
	DecodeTopLevel(data []byte, value interface{}) error
}
