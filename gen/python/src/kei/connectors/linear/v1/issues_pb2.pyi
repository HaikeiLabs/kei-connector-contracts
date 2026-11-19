from kei.contracts.v1 import options_pb2 as _options_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Label(_message.Message):
    __slots__ = ("id", "name")
    ID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    id: str
    name: str
    def __init__(self, id: _Optional[str] = ..., name: _Optional[str] = ...) -> None: ...

class Issue(_message.Message):
    __slots__ = ("id", "identifier", "title", "description", "url", "created_at", "updated_at", "labels")
    ID_FIELD_NUMBER: _ClassVar[int]
    IDENTIFIER_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    URL_FIELD_NUMBER: _ClassVar[int]
    CREATED_AT_FIELD_NUMBER: _ClassVar[int]
    UPDATED_AT_FIELD_NUMBER: _ClassVar[int]
    LABELS_FIELD_NUMBER: _ClassVar[int]
    id: str
    identifier: str
    title: str
    description: str
    url: str
    created_at: str
    updated_at: str
    labels: _containers.RepeatedCompositeFieldContainer[Label]
    def __init__(self, id: _Optional[str] = ..., identifier: _Optional[str] = ..., title: _Optional[str] = ..., description: _Optional[str] = ..., url: _Optional[str] = ..., created_at: _Optional[str] = ..., updated_at: _Optional[str] = ..., labels: _Optional[_Iterable[_Union[Label, _Mapping]]] = ...) -> None: ...

class GetIssueRequest(_message.Message):
    __slots__ = ("id",)
    ID_FIELD_NUMBER: _ClassVar[int]
    id: str
    def __init__(self, id: _Optional[str] = ...) -> None: ...

class GetIssueResponse(_message.Message):
    __slots__ = ("issue",)
    ISSUE_FIELD_NUMBER: _ClassVar[int]
    issue: Issue
    def __init__(self, issue: _Optional[_Union[Issue, _Mapping]] = ...) -> None: ...

class CreateIssueRequest(_message.Message):
    __slots__ = ("title", "team_id", "description", "priority")
    TITLE_FIELD_NUMBER: _ClassVar[int]
    TEAM_ID_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    PRIORITY_FIELD_NUMBER: _ClassVar[int]
    title: str
    team_id: str
    description: str
    priority: int
    def __init__(self, title: _Optional[str] = ..., team_id: _Optional[str] = ..., description: _Optional[str] = ..., priority: _Optional[int] = ...) -> None: ...

class CreateIssueResponse(_message.Message):
    __slots__ = ("success", "issue")
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    ISSUE_FIELD_NUMBER: _ClassVar[int]
    success: bool
    issue: Issue
    def __init__(self, success: _Optional[bool] = ..., issue: _Optional[_Union[Issue, _Mapping]] = ...) -> None: ...

class SearchIssuesRequest(_message.Message):
    __slots__ = ("term", "first", "after")
    TERM_FIELD_NUMBER: _ClassVar[int]
    FIRST_FIELD_NUMBER: _ClassVar[int]
    AFTER_FIELD_NUMBER: _ClassVar[int]
    term: str
    first: int
    after: str
    def __init__(self, term: _Optional[str] = ..., first: _Optional[int] = ..., after: _Optional[str] = ...) -> None: ...

class PageInfo(_message.Message):
    __slots__ = ("has_next_page", "end_cursor")
    HAS_NEXT_PAGE_FIELD_NUMBER: _ClassVar[int]
    END_CURSOR_FIELD_NUMBER: _ClassVar[int]
    has_next_page: bool
    end_cursor: str
    def __init__(self, has_next_page: _Optional[bool] = ..., end_cursor: _Optional[str] = ...) -> None: ...

class SearchIssuesResponse(_message.Message):
    __slots__ = ("nodes", "page_info")
    NODES_FIELD_NUMBER: _ClassVar[int]
    PAGE_INFO_FIELD_NUMBER: _ClassVar[int]
    nodes: _containers.RepeatedCompositeFieldContainer[Issue]
    page_info: PageInfo
    def __init__(self, nodes: _Optional[_Iterable[_Union[Issue, _Mapping]]] = ..., page_info: _Optional[_Union[PageInfo, _Mapping]] = ...) -> None: ...

class UpdateIssueRequest(_message.Message):
    __slots__ = ("id", "title", "description")
    ID_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    id: str
    title: str
    description: str
    def __init__(self, id: _Optional[str] = ..., title: _Optional[str] = ..., description: _Optional[str] = ...) -> None: ...

class UpdateIssueResponse(_message.Message):
    __slots__ = ("success", "issue")
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    ISSUE_FIELD_NUMBER: _ClassVar[int]
    success: bool
    issue: Issue
    def __init__(self, success: _Optional[bool] = ..., issue: _Optional[_Union[Issue, _Mapping]] = ...) -> None: ...
