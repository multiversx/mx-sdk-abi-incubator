package abi

type codec interface {
	EncodeNested(value interface{}) ([]byte, error)
	EncodeTopLevel(value interface{}) ([]byte, error)
	DecodeNested(data []byte, value interface{}) error
	DecodeTopLevel(data []byte, value interface{}) error
}
