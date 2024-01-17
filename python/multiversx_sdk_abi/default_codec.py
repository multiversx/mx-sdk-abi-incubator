from typing import Any, Protocol

from python.multiversx_sdk_abi.errors import ErrUnsupportedType
from python.multiversx_sdk_abi.values import (U8Value, U16Value, U32Value,
                                              U64Value)


class IDataWriter(Protocol):
    def write(self, data: bytes) -> None: ...


class DefaultCodec:
    def __init__(self):
        pass

    def encode_nested(self, writer: IDataWriter, value: Any):
        if isinstance(value, U8Value):
            return
        elif isinstance(value, U16Value):
            return
        elif isinstance(value, U32Value):
            return
        elif isinstance(value, U64Value):
            return
        else:
            raise ErrUnsupportedType("DefaultCodec.EncodeNested()", value)

    def encode_top_level(self, writer: IDataWriter, value: Any):
        if isinstance(value, U8Value):
            return
        elif isinstance(value, U16Value):
            return
        elif isinstance(value, U32Value):
            return
        elif isinstance(value, U64Value):
            return
        else:
            raise ErrUnsupportedType("DefaultCodec.EncodeTopLevel()", value)
