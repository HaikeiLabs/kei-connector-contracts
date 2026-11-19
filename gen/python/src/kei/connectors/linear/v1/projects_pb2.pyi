from kei.contracts.v1 import options_pb2 as _options_pb2
from kei.connectors.linear.v1 import issues_pb2 as _issues_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Project(_message.Message):
    __slots__ = ("id", "name", "description", "url")
    ID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    URL_FIELD_NUMBER: _ClassVar[int]
    id: str
    name: str
    description: str
    url: str
    def __init__(self, id: _Optional[str] = ..., name: _Optional[str] = ..., description: _Optional[str] = ..., url: _Optional[str] = ...) -> None: ...

class ListProjectsRequest(_message.Message):
    __slots__ = ("first", "after")
    FIRST_FIELD_NUMBER: _ClassVar[int]
    AFTER_FIELD_NUMBER: _ClassVar[int]
    first: int
    after: str
    def __init__(self, first: _Optional[int] = ..., after: _Optional[str] = ...) -> None: ...

class ListProjectsResponse(_message.Message):
    __slots__ = ("nodes", "page_info")
    NODES_FIELD_NUMBER: _ClassVar[int]
    PAGE_INFO_FIELD_NUMBER: _ClassVar[int]
    nodes: _containers.RepeatedCompositeFieldContainer[Project]
    page_info: _issues_pb2.PageInfo
    def __init__(self, nodes: _Optional[_Iterable[_Union[Project, _Mapping]]] = ..., page_info: _Optional[_Union[_issues_pb2.PageInfo, _Mapping]] = ...) -> None: ...
