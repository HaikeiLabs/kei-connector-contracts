from kei.contracts.v1 import options_pb2 as _options_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ListIssuesRequest(_message.Message):
    __slots__ = ("organization_slug", "query", "per_page")
    ORGANIZATION_SLUG_FIELD_NUMBER: _ClassVar[int]
    QUERY_FIELD_NUMBER: _ClassVar[int]
    PER_PAGE_FIELD_NUMBER: _ClassVar[int]
    organization_slug: str
    query: str
    per_page: int
    def __init__(self, organization_slug: _Optional[str] = ..., query: _Optional[str] = ..., per_page: _Optional[int] = ...) -> None: ...

class IssueSummary(_message.Message):
    __slots__ = ("id", "short_id", "title", "status", "level", "share_id", "permalink")
    ID_FIELD_NUMBER: _ClassVar[int]
    SHORT_ID_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    LEVEL_FIELD_NUMBER: _ClassVar[int]
    SHARE_ID_FIELD_NUMBER: _ClassVar[int]
    PERMALINK_FIELD_NUMBER: _ClassVar[int]
    id: str
    short_id: str
    title: str
    status: str
    level: str
    share_id: str
    permalink: str
    def __init__(self, id: _Optional[str] = ..., short_id: _Optional[str] = ..., title: _Optional[str] = ..., status: _Optional[str] = ..., level: _Optional[str] = ..., share_id: _Optional[str] = ..., permalink: _Optional[str] = ...) -> None: ...

class ListIssuesResponse(_message.Message):
    __slots__ = ("issues",)
    ISSUES_FIELD_NUMBER: _ClassVar[int]
    issues: _containers.RepeatedCompositeFieldContainer[IssueSummary]
    def __init__(self, issues: _Optional[_Iterable[_Union[IssueSummary, _Mapping]]] = ...) -> None: ...

class GetIssueRequest(_message.Message):
    __slots__ = ("organization_slug", "issue_id")
    ORGANIZATION_SLUG_FIELD_NUMBER: _ClassVar[int]
    ISSUE_ID_FIELD_NUMBER: _ClassVar[int]
    organization_slug: str
    issue_id: str
    def __init__(self, organization_slug: _Optional[str] = ..., issue_id: _Optional[str] = ...) -> None: ...

class GetIssueResponse(_message.Message):
    __slots__ = ("id", "short_id", "title", "status", "level", "share_id", "permalink")
    ID_FIELD_NUMBER: _ClassVar[int]
    SHORT_ID_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    LEVEL_FIELD_NUMBER: _ClassVar[int]
    SHARE_ID_FIELD_NUMBER: _ClassVar[int]
    PERMALINK_FIELD_NUMBER: _ClassVar[int]
    id: str
    short_id: str
    title: str
    status: str
    level: str
    share_id: str
    permalink: str
    def __init__(self, id: _Optional[str] = ..., short_id: _Optional[str] = ..., title: _Optional[str] = ..., status: _Optional[str] = ..., level: _Optional[str] = ..., share_id: _Optional[str] = ..., permalink: _Optional[str] = ...) -> None: ...
