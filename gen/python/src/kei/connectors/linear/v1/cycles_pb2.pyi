from kei.contracts.v1 import options_pb2 as _options_pb2
from kei.connectors.linear.v1 import issues_pb2 as _issues_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Cycle(_message.Message):
    __slots__ = ("id", "number", "starts_at", "ends_at", "is_active")
    ID_FIELD_NUMBER: _ClassVar[int]
    NUMBER_FIELD_NUMBER: _ClassVar[int]
    STARTS_AT_FIELD_NUMBER: _ClassVar[int]
    ENDS_AT_FIELD_NUMBER: _ClassVar[int]
    IS_ACTIVE_FIELD_NUMBER: _ClassVar[int]
    id: str
    number: int
    starts_at: str
    ends_at: str
    is_active: bool
    def __init__(self, id: _Optional[str] = ..., number: _Optional[int] = ..., starts_at: _Optional[str] = ..., ends_at: _Optional[str] = ..., is_active: _Optional[bool] = ...) -> None: ...

class ListCyclesRequest(_message.Message):
    __slots__ = ("team_id", "first", "after")
    TEAM_ID_FIELD_NUMBER: _ClassVar[int]
    FIRST_FIELD_NUMBER: _ClassVar[int]
    AFTER_FIELD_NUMBER: _ClassVar[int]
    team_id: str
    first: int
    after: str
    def __init__(self, team_id: _Optional[str] = ..., first: _Optional[int] = ..., after: _Optional[str] = ...) -> None: ...

class ListCyclesResponse(_message.Message):
    __slots__ = ("nodes", "page_info")
    NODES_FIELD_NUMBER: _ClassVar[int]
    PAGE_INFO_FIELD_NUMBER: _ClassVar[int]
    nodes: _containers.RepeatedCompositeFieldContainer[Cycle]
    page_info: _issues_pb2.PageInfo
    def __init__(self, nodes: _Optional[_Iterable[_Union[Cycle, _Mapping]]] = ..., page_info: _Optional[_Union[_issues_pb2.PageInfo, _Mapping]] = ...) -> None: ...
