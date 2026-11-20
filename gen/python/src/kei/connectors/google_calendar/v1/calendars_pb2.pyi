from kei.contracts.v1 import options_pb2 as _options_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ListCalendarsRequest(_message.Message):
    __slots__ = ("max_results", "page_token")
    MAX_RESULTS_FIELD_NUMBER: _ClassVar[int]
    PAGE_TOKEN_FIELD_NUMBER: _ClassVar[int]
    max_results: int
    page_token: str
    def __init__(self, max_results: _Optional[int] = ..., page_token: _Optional[str] = ...) -> None: ...

class CalendarListEntry(_message.Message):
    __slots__ = ("id", "summary", "description", "time_zone", "primary")
    ID_FIELD_NUMBER: _ClassVar[int]
    SUMMARY_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    TIME_ZONE_FIELD_NUMBER: _ClassVar[int]
    PRIMARY_FIELD_NUMBER: _ClassVar[int]
    id: str
    summary: str
    description: str
    time_zone: str
    primary: bool
    def __init__(self, id: _Optional[str] = ..., summary: _Optional[str] = ..., description: _Optional[str] = ..., time_zone: _Optional[str] = ..., primary: _Optional[bool] = ...) -> None: ...

class ListCalendarsResponse(_message.Message):
    __slots__ = ("items", "next_page_token")
    ITEMS_FIELD_NUMBER: _ClassVar[int]
    NEXT_PAGE_TOKEN_FIELD_NUMBER: _ClassVar[int]
    items: _containers.RepeatedCompositeFieldContainer[CalendarListEntry]
    next_page_token: str
    def __init__(self, items: _Optional[_Iterable[_Union[CalendarListEntry, _Mapping]]] = ..., next_page_token: _Optional[str] = ...) -> None: ...
