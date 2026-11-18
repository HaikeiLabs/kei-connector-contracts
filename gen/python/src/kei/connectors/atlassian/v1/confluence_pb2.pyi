from kei.contracts.v1 import options_pb2 as _options_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class GetConfluencePageRequest(_message.Message):
    __slots__ = ("page_id",)
    PAGE_ID_FIELD_NUMBER: _ClassVar[int]
    page_id: int
    def __init__(self, page_id: _Optional[int] = ...) -> None: ...

class GetConfluencePageResponse(_message.Message):
    __slots__ = ("id", "title", "space_id", "parent_id")
    ID_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    SPACE_ID_FIELD_NUMBER: _ClassVar[int]
    PARENT_ID_FIELD_NUMBER: _ClassVar[int]
    id: int
    title: str
    space_id: str
    parent_id: str
    def __init__(self, id: _Optional[int] = ..., title: _Optional[str] = ..., space_id: _Optional[str] = ..., parent_id: _Optional[str] = ...) -> None: ...

class ListConfluencePagesRequest(_message.Message):
    __slots__ = ("space_id", "limit")
    SPACE_ID_FIELD_NUMBER: _ClassVar[int]
    LIMIT_FIELD_NUMBER: _ClassVar[int]
    space_id: int
    limit: int
    def __init__(self, space_id: _Optional[int] = ..., limit: _Optional[int] = ...) -> None: ...

class ConfluencePageSummary(_message.Message):
    __slots__ = ("id", "title", "space_id")
    ID_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    SPACE_ID_FIELD_NUMBER: _ClassVar[int]
    id: int
    title: str
    space_id: str
    def __init__(self, id: _Optional[int] = ..., title: _Optional[str] = ..., space_id: _Optional[str] = ...) -> None: ...

class ListConfluencePagesResponse(_message.Message):
    __slots__ = ("results",)
    RESULTS_FIELD_NUMBER: _ClassVar[int]
    results: _containers.RepeatedCompositeFieldContainer[ConfluencePageSummary]
    def __init__(self, results: _Optional[_Iterable[_Union[ConfluencePageSummary, _Mapping]]] = ...) -> None: ...
