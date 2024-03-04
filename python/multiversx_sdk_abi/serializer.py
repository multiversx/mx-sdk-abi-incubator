from typing import Any, List, Protocol

from python.multiversx_sdk_abi.errors import (ErrNilInputValue,
                                              ErrVariadicMustBeLast)
from python.multiversx_sdk_abi.values import (InputMultiValue,
                                              InputVariadicValues)


class IDataWriter(Protocol):
    def write(self, data: bytes) -> None: ...
    def goto_next_part(self) -> None: ...


class ICodec(Protocol):
    def encode_top_level(self, writer: IDataWriter, value: Any) -> None: ...


class Serializer:
    def __init__(self, codec: ICodec):
        self.codec = codec

    def serialize(self, writer: IDataWriter, inputValues: List[Any]) -> None:
        for i, value in enumerate(inputValues):
            if value is None:
                raise ErrNilInputValue()

            if isinstance(value, InputMultiValue):
                self.serialize_input_multi_value(writer, value)
            elif isinstance(value, InputVariadicValues):
                if i != len(inputValues) - 1:
                    raise ErrVariadicMustBeLast()

                self.serialize_input_variadic_values(writer, value)
            else:
                self.serialize_directly_encodable_value(writer, value)

    def serialize_input_multi_value(self, writer: IDataWriter, value: InputMultiValue) -> None:
        for item in value.items:
            self.serialize(writer, [item])

    def serialize_input_variadic_values(self, writer: IDataWriter, value: InputVariadicValues) -> None:
        for item in value.items:
            self.serialize(writer, [item])

    def serialize_directly_encodable_value(self, writer: IDataWriter, value: Any) -> None:
        writer.goto_next_part()
        self.codec.encode_top_level(writer, value)
