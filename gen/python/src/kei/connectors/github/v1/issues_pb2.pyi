from kei.contracts.v1 import options_pb2 as _options_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class GetIssueRequest(_message.Message):
    __slots__ = ("owner", "repository", "issue_number")
    OWNER_FIELD_NUMBER: _ClassVar[int]
    REPOSITORY_FIELD_NUMBER: _ClassVar[int]
    ISSUE_NUMBER_FIELD_NUMBER: _ClassVar[int]
    owner: str
    repository: str
    issue_number: int
    def __init__(self, owner: _Optional[str] = ..., repository: _Optional[str] = ..., issue_number: _Optional[int] = ...) -> None: ...

class GetIssueResponse(_message.Message):
    __slots__ = ("id", "number", "title", "state", "body", "html_url")
    ID_FIELD_NUMBER: _ClassVar[int]
    NUMBER_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    STATE_FIELD_NUMBER: _ClassVar[int]
    BODY_FIELD_NUMBER: _ClassVar[int]
    HTML_URL_FIELD_NUMBER: _ClassVar[int]
    id: int
    number: int
    title: str
    state: str
    body: str
    html_url: str
    def __init__(self, id: _Optional[int] = ..., number: _Optional[int] = ..., title: _Optional[str] = ..., state: _Optional[str] = ..., body: _Optional[str] = ..., html_url: _Optional[str] = ...) -> None: ...

class ListIssuesRequest(_message.Message):
    __slots__ = ("owner", "repository", "state", "per_page")
    OWNER_FIELD_NUMBER: _ClassVar[int]
    REPOSITORY_FIELD_NUMBER: _ClassVar[int]
    STATE_FIELD_NUMBER: _ClassVar[int]
    PER_PAGE_FIELD_NUMBER: _ClassVar[int]
    owner: str
    repository: str
    state: str
    per_page: int
    def __init__(self, owner: _Optional[str] = ..., repository: _Optional[str] = ..., state: _Optional[str] = ..., per_page: _Optional[int] = ...) -> None: ...

class IssueSummary(_message.Message):
    __slots__ = ("id", "number", "title", "state", "body", "html_url")
    ID_FIELD_NUMBER: _ClassVar[int]
    NUMBER_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    STATE_FIELD_NUMBER: _ClassVar[int]
    BODY_FIELD_NUMBER: _ClassVar[int]
    HTML_URL_FIELD_NUMBER: _ClassVar[int]
    id: int
    number: int
    title: str
    state: str
    body: str
    html_url: str
    def __init__(self, id: _Optional[int] = ..., number: _Optional[int] = ..., title: _Optional[str] = ..., state: _Optional[str] = ..., body: _Optional[str] = ..., html_url: _Optional[str] = ...) -> None: ...

class ListIssuesResponse(_message.Message):
    __slots__ = ("issues",)
    ISSUES_FIELD_NUMBER: _ClassVar[int]
    issues: _containers.RepeatedCompositeFieldContainer[IssueSummary]
    def __init__(self, issues: _Optional[_Iterable[_Union[IssueSummary, _Mapping]]] = ...) -> None: ...
