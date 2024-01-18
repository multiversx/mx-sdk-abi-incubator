package abi

type dataReader interface {
	Read(numBytes int) ([]byte, error)
	ReadWholePart() ([]byte, error)
	GotoNextPart() error
	IsCurrentPartEmpty() bool
	IsEndOfData() bool
}

type dataWriter interface {
	Write(data []byte) (int, error)
	GotoNextPart()
}

type codec interface {
	EncodeNested(value interface{}) ([]byte, error)
	EncodeTopLevel(value interface{}) ([]byte, error)
	DecodeNested(reader dataReader, value interface{}) error
	DecodeTopLevel(reader dataReader, value interface{}) error
}
