from kei.contracts.v1 import options_pb2 as _options_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class GetJiraIssueRequest(_message.Message):
    __slots__ = ("issue_key",)
    ISSUE_KEY_FIELD_NUMBER: _ClassVar[int]
    issue_key: str
    def __init__(self, issue_key: _Optional[str] = ...) -> None: ...

class JiraIssueFields(_message.Message):
    __slots__ = ("summary", "description")
    SUMMARY_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    summary: str
    description: str
    def __init__(self, summary: _Optional[str] = ..., description: _Optional[str] = ...) -> None: ...

class GetJiraIssueResponse(_message.Message):
    __slots__ = ("id", "key", "fields")
    ID_FIELD_NUMBER: _ClassVar[int]
    KEY_FIELD_NUMBER: _ClassVar[int]
    FIELDS_FIELD_NUMBER: _ClassVar[int]
    id: str
    key: str
    fields: JiraIssueFields
    def __init__(self, id: _Optional[str] = ..., key: _Optional[str] = ..., fields: _Optional[_Union[JiraIssueFields, _Mapping]] = ...) -> None: ...

class SearchJiraIssuesRequest(_message.Message):
    __slots__ = ("jql", "next_page_token", "max_results")
    JQL_FIELD_NUMBER: _ClassVar[int]
    NEXT_PAGE_TOKEN_FIELD_NUMBER: _ClassVar[int]
    MAX_RESULTS_FIELD_NUMBER: _ClassVar[int]
    jql: str
    next_page_token: str
    max_results: int
    def __init__(self, jql: _Optional[str] = ..., next_page_token: _Optional[str] = ..., max_results: _Optional[int] = ...) -> None: ...

class JiraIssueSummary(_message.Message):
    __slots__ = ("id", "key", "fields")
    ID_FIELD_NUMBER: _ClassVar[int]
    KEY_FIELD_NUMBER: _ClassVar[int]
    FIELDS_FIELD_NUMBER: _ClassVar[int]
    id: str
    key: str
    fields: JiraIssueFields
    def __init__(self, id: _Optional[str] = ..., key: _Optional[str] = ..., fields: _Optional[_Union[JiraIssueFields, _Mapping]] = ...) -> None: ...

class SearchJiraIssuesResponse(_message.Message):
    __slots__ = ("issues", "next_page_token")
    ISSUES_FIELD_NUMBER: _ClassVar[int]
    NEXT_PAGE_TOKEN_FIELD_NUMBER: _ClassVar[int]
    issues: _containers.RepeatedCompositeFieldContainer[JiraIssueSummary]
    next_page_token: str
    def __init__(self, issues: _Optional[_Iterable[_Union[JiraIssueSummary, _Mapping]]] = ..., next_page_token: _Optional[str] = ...) -> None: ...
