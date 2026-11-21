from google.protobuf import struct_pb2 as _struct_pb2
from kei.contracts.v1 import options_pb2 as _options_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class SearchRequest(_message.Message):
    __slots__ = ("query", "filter", "start_cursor", "page_size")
    QUERY_FIELD_NUMBER: _ClassVar[int]
    FILTER_FIELD_NUMBER: _ClassVar[int]
    START_CURSOR_FIELD_NUMBER: _ClassVar[int]
    PAGE_SIZE_FIELD_NUMBER: _ClassVar[int]
    query: str
    filter: _struct_pb2.Struct
    start_cursor: str
    page_size: int
    def __init__(self, query: _Optional[str] = ..., filter: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ..., start_cursor: _Optional[str] = ..., page_size: _Optional[int] = ...) -> None: ...

class SearchResponse(_message.Message):
    __slots__ = ("results", "next_cursor", "has_more")
    RESULTS_FIELD_NUMBER: _ClassVar[int]
    NEXT_CURSOR_FIELD_NUMBER: _ClassVar[int]
    HAS_MORE_FIELD_NUMBER: _ClassVar[int]
    results: _containers.RepeatedCompositeFieldContainer[_struct_pb2.Struct]
    next_cursor: str
    has_more: bool
    def __init__(self, results: _Optional[_Iterable[_Union[_struct_pb2.Struct, _Mapping]]] = ..., next_cursor: _Optional[str] = ..., has_more: _Optional[bool] = ...) -> None: ...
