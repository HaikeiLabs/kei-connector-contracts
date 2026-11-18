from kei.contracts.v1 import options_pb2 as _options_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class QueryDatabaseRequest(_message.Message):
    __slots__ = ("database_id", "start_cursor", "page_size")
    DATABASE_ID_FIELD_NUMBER: _ClassVar[int]
    START_CURSOR_FIELD_NUMBER: _ClassVar[int]
    PAGE_SIZE_FIELD_NUMBER: _ClassVar[int]
    database_id: str
    start_cursor: str
    page_size: int
    def __init__(self, database_id: _Optional[str] = ..., start_cursor: _Optional[str] = ..., page_size: _Optional[int] = ...) -> None: ...

class PageSummary(_message.Message):
    __slots__ = ("id", "created_time", "last_edited_time", "archived", "url", "public_url")
    ID_FIELD_NUMBER: _ClassVar[int]
    CREATED_TIME_FIELD_NUMBER: _ClassVar[int]
    LAST_EDITED_TIME_FIELD_NUMBER: _ClassVar[int]
    ARCHIVED_FIELD_NUMBER: _ClassVar[int]
    URL_FIELD_NUMBER: _ClassVar[int]
    PUBLIC_URL_FIELD_NUMBER: _ClassVar[int]
    id: str
    created_time: str
    last_edited_time: str
    archived: bool
    url: str
    public_url: str
    def __init__(self, id: _Optional[str] = ..., created_time: _Optional[str] = ..., last_edited_time: _Optional[str] = ..., archived: _Optional[bool] = ..., url: _Optional[str] = ..., public_url: _Optional[str] = ...) -> None: ...

class QueryDatabaseResponse(_message.Message):
    __slots__ = ("results", "next_cursor", "has_more")
    RESULTS_FIELD_NUMBER: _ClassVar[int]
    NEXT_CURSOR_FIELD_NUMBER: _ClassVar[int]
    HAS_MORE_FIELD_NUMBER: _ClassVar[int]
    results: _containers.RepeatedCompositeFieldContainer[PageSummary]
    next_cursor: str
    has_more: bool
    def __init__(self, results: _Optional[_Iterable[_Union[PageSummary, _Mapping]]] = ..., next_cursor: _Optional[str] = ..., has_more: _Optional[bool] = ...) -> None: ...
