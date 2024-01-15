package abi

type dataReader interface {
	Read(numBytes int) ([]byte, error)
	ReadWholePart() ([]byte, error)
	HasUnreadData() bool
	GotoNextPart() error
}

type dataWriter interface {
	Write(data []byte) error
	GotoNextPart()
}

type codec interface {
	EncodeNested(writer dataWriter, value interface{}) error
	EncodeTopLevel(writer dataWriter, value interface{}) error
	DecodeNested(reader dataReader, value interface{}) error
	DecodeTopLevel(reader dataReader, value interface{}) error
}
