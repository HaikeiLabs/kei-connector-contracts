from google.protobuf import struct_pb2 as _struct_pb2
from kei.contracts.v1 import options_pb2 as _options_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ListRecordsRequest(_message.Message):
    __slots__ = ("collection", "limit", "cursor")
    COLLECTION_FIELD_NUMBER: _ClassVar[int]
    LIMIT_FIELD_NUMBER: _ClassVar[int]
    CURSOR_FIELD_NUMBER: _ClassVar[int]
    collection: str
    limit: int
    cursor: str
    def __init__(self, collection: _Optional[str] = ..., limit: _Optional[int] = ..., cursor: _Optional[str] = ...) -> None: ...

class ListRecordsResponse(_message.Message):
    __slots__ = ("records", "next_cursor")
    RECORDS_FIELD_NUMBER: _ClassVar[int]
    NEXT_CURSOR_FIELD_NUMBER: _ClassVar[int]
    records: _containers.RepeatedCompositeFieldContainer[_struct_pb2.Struct]
    next_cursor: str
    def __init__(self, records: _Optional[_Iterable[_Union[_struct_pb2.Struct, _Mapping]]] = ..., next_cursor: _Optional[str] = ...) -> None: ...

class GetRecordRequest(_message.Message):
    __slots__ = ("collection", "record_id")
    COLLECTION_FIELD_NUMBER: _ClassVar[int]
    RECORD_ID_FIELD_NUMBER: _ClassVar[int]
    collection: str
    record_id: str
    def __init__(self, collection: _Optional[str] = ..., record_id: _Optional[str] = ...) -> None: ...

class GetRecordResponse(_message.Message):
    __slots__ = ("record",)
    RECORD_FIELD_NUMBER: _ClassVar[int]
    record: _struct_pb2.Struct
    def __init__(self, record: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ...) -> None: ...
