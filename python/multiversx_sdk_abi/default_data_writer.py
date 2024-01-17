
from typing import List

from python.multiversx_sdk_abi.constants import PARTS_SEPARATOR
from python.multiversx_sdk_abi.errors import \
    ErrWriterCannotWriteSinceThereIsNoPart


class DefaultDataWriter:
    def __init__(self):
        """
        Creates a new DefaultDataWriter.
        A newly-created writer has no parts.
        Parts are created by calling goto_next_part().
        """
        self.parts: List[bytes] = []

    def write(self, data: bytes):
        if len(self.parts) == 0:
            raise ErrWriterCannotWriteSinceThereIsNoPart()

        part_index = len(self.parts) - 1
        self.parts[part_index] += data

    def goto_next_part(self):
        self.parts.append(b"")

    def get_parts(self) -> List[bytes]:
        return self.parts

    def __str__(self):
        parts_hex = [part.hex() for part in self.parts]
        return PARTS_SEPARATOR.join(parts_hex)
