from typing import Any, Callable, List


class U8Value:
    def __init__(self, value: int = 0):
        self.value = value


class U16Value:
    def __init__(self, value: int = 0):
        self.value = value


class U32Value:
    def __init__(self, value: int = 0):
        self.value = value


class U64Value:
    def __init__(self, value: int = 0):
        self.value = value


class I8Value:
    def __init__(self, value: int = 0):
        self.value = value


class I16Value:
    def __init__(self, value: int = 0):
        self.value = value


class I32Value:
    def __init__(self, value: int = 0):
        self.value = value


class I64Value:
    def __init__(self, value: int = 0):
        self.value = value


class BytesValue:
    def __init__(self, value: bytes = b''):
        self.value = value


class StringValue:
    def __init__(self, value: str = ''):
        self.value = value


class BoolValue:
    def __init__(self, value: bool = False):
        self.value = value


class OptionValue:
    def __init__(self, value: Any = None):
        self.value = value


class Field:
    def __init__(self, name: str = "", value: Any = None):
        self.name = name
        self.value = value


class StructValue:
    def __init__(self, fields: List[Field] = []):
        self.fields = fields


class TupleValue:
    def __init__(self, fields: List[Any] = []):
        self.fields = fields


class EnumValue:
    def __init__(self, discriminant: int = 0, fields: List[Field] = []):
        self.discriminant = discriminant
        self.fields = fields


class InputListValue:
    def __init__(self, items: List[Any] = []):
        self.values = items


class OutputListValue:
    def __init__(self, item_creator: Callable[[], Any]):
        self.items = []
        self.item_creator = item_creator


class InputMultiValue:
    def __init__(self, items: List[Any] = []):
        self.items = items


class InputVariadicValues:
    def __init__(self, items: List[Any] = []):
        self.items = items


class OutputMultiValue:
    def __init__(self, item_creator: Callable[[], Any]):
        self.items = []
        self.item_creator = item_creator


class OutputVariadicValues:
    def __init__(self, item_creator: Callable[[], Any]):
        self.items = []
        self.item_creator = item_creator


class OptionalValue:
    def __init__(self, value: Any = None):
        self.value = value
        self.is_set = value is not None
