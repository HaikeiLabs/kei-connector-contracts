from kei.contracts.v1 import options_pb2 as _options_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ListCheckRunsRequest(_message.Message):
    __slots__ = ("owner", "repository", "ref", "per_page", "page")
    OWNER_FIELD_NUMBER: _ClassVar[int]
    REPOSITORY_FIELD_NUMBER: _ClassVar[int]
    REF_FIELD_NUMBER: _ClassVar[int]
    PER_PAGE_FIELD_NUMBER: _ClassVar[int]
    PAGE_FIELD_NUMBER: _ClassVar[int]
    owner: str
    repository: str
    ref: str
    per_page: int
    page: int
    def __init__(self, owner: _Optional[str] = ..., repository: _Optional[str] = ..., ref: _Optional[str] = ..., per_page: _Optional[int] = ..., page: _Optional[int] = ...) -> None: ...

class CheckRunSummary(_message.Message):
    __slots__ = ("id", "name", "head_sha", "status", "conclusion", "html_url")
    ID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    HEAD_SHA_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    CONCLUSION_FIELD_NUMBER: _ClassVar[int]
    HTML_URL_FIELD_NUMBER: _ClassVar[int]
    id: int
    name: str
    head_sha: str
    status: str
    conclusion: str
    html_url: str
    def __init__(self, id: _Optional[int] = ..., name: _Optional[str] = ..., head_sha: _Optional[str] = ..., status: _Optional[str] = ..., conclusion: _Optional[str] = ..., html_url: _Optional[str] = ...) -> None: ...

class ListCheckRunsResponse(_message.Message):
    __slots__ = ("check_runs", "total_count", "next_page")
    CHECK_RUNS_FIELD_NUMBER: _ClassVar[int]
    TOTAL_COUNT_FIELD_NUMBER: _ClassVar[int]
    NEXT_PAGE_FIELD_NUMBER: _ClassVar[int]
    check_runs: _containers.RepeatedCompositeFieldContainer[CheckRunSummary]
    total_count: int
    next_page: int
    def __init__(self, check_runs: _Optional[_Iterable[_Union[CheckRunSummary, _Mapping]]] = ..., total_count: _Optional[int] = ..., next_page: _Optional[int] = ...) -> None: ...
