
from typing import Any


class ErrNilInputValue(Exception):
    def __init__(self):
        super().__init__("nil input value")


class ErrVariadicMustBeLast(Exception):
    def __init__(self):
        super().__init__("variadic must be last")


class ErrUnsupportedType(Exception):
    def __init__(self, when: str, value: Any):
        super().__init__(f"{when}, unsupported type: {type(value)}")
