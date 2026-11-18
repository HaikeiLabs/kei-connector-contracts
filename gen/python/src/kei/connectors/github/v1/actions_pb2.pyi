from kei.contracts.v1 import options_pb2 as _options_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ListWorkflowRunsRequest(_message.Message):
    __slots__ = ("owner", "repository", "per_page", "page")
    OWNER_FIELD_NUMBER: _ClassVar[int]
    REPOSITORY_FIELD_NUMBER: _ClassVar[int]
    PER_PAGE_FIELD_NUMBER: _ClassVar[int]
    PAGE_FIELD_NUMBER: _ClassVar[int]
    owner: str
    repository: str
    per_page: int
    page: int
    def __init__(self, owner: _Optional[str] = ..., repository: _Optional[str] = ..., per_page: _Optional[int] = ..., page: _Optional[int] = ...) -> None: ...

class WorkflowRunSummary(_message.Message):
    __slots__ = ("id", "name", "head_branch", "status", "conclusion", "html_url")
    ID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    HEAD_BRANCH_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    CONCLUSION_FIELD_NUMBER: _ClassVar[int]
    HTML_URL_FIELD_NUMBER: _ClassVar[int]
    id: int
    name: str
    head_branch: str
    status: str
    conclusion: str
    html_url: str
    def __init__(self, id: _Optional[int] = ..., name: _Optional[str] = ..., head_branch: _Optional[str] = ..., status: _Optional[str] = ..., conclusion: _Optional[str] = ..., html_url: _Optional[str] = ...) -> None: ...

class ListWorkflowRunsResponse(_message.Message):
    __slots__ = ("workflow_runs", "next_page")
    WORKFLOW_RUNS_FIELD_NUMBER: _ClassVar[int]
    NEXT_PAGE_FIELD_NUMBER: _ClassVar[int]
    workflow_runs: _containers.RepeatedCompositeFieldContainer[WorkflowRunSummary]
    next_page: int
    def __init__(self, workflow_runs: _Optional[_Iterable[_Union[WorkflowRunSummary, _Mapping]]] = ..., next_page: _Optional[int] = ...) -> None: ...
