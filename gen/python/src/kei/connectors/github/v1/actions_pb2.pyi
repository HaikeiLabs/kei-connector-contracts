from kei.contracts.v1 import options_pb2 as _options_pb2
from kei.connectors.github.v1 import pulls_pb2 as _pulls_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ListWorkflowRunsRequest(_message.Message):
    __slots__ = ("owner", "repository", "per_page", "page", "branch", "event")
    OWNER_FIELD_NUMBER: _ClassVar[int]
    REPOSITORY_FIELD_NUMBER: _ClassVar[int]
    PER_PAGE_FIELD_NUMBER: _ClassVar[int]
    PAGE_FIELD_NUMBER: _ClassVar[int]
    BRANCH_FIELD_NUMBER: _ClassVar[int]
    EVENT_FIELD_NUMBER: _ClassVar[int]
    owner: str
    repository: str
    per_page: int
    page: int
    branch: str
    event: str
    def __init__(self, owner: _Optional[str] = ..., repository: _Optional[str] = ..., per_page: _Optional[int] = ..., page: _Optional[int] = ..., branch: _Optional[str] = ..., event: _Optional[str] = ...) -> None: ...

class WorkflowRunSummary(_message.Message):
    __slots__ = ("id", "name", "head_branch", "status", "conclusion", "html_url", "actor", "head_sha")
    ID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    HEAD_BRANCH_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    CONCLUSION_FIELD_NUMBER: _ClassVar[int]
    HTML_URL_FIELD_NUMBER: _ClassVar[int]
    ACTOR_FIELD_NUMBER: _ClassVar[int]
    HEAD_SHA_FIELD_NUMBER: _ClassVar[int]
    id: int
    name: str
    head_branch: str
    status: str
    conclusion: str
    html_url: str
    actor: _pulls_pb2.User
    head_sha: str
    def __init__(self, id: _Optional[int] = ..., name: _Optional[str] = ..., head_branch: _Optional[str] = ..., status: _Optional[str] = ..., conclusion: _Optional[str] = ..., html_url: _Optional[str] = ..., actor: _Optional[_Union[_pulls_pb2.User, _Mapping]] = ..., head_sha: _Optional[str] = ...) -> None: ...

class ListWorkflowRunsResponse(_message.Message):
    __slots__ = ("workflow_runs", "next_page")
    WORKFLOW_RUNS_FIELD_NUMBER: _ClassVar[int]
    NEXT_PAGE_FIELD_NUMBER: _ClassVar[int]
    workflow_runs: _containers.RepeatedCompositeFieldContainer[WorkflowRunSummary]
    next_page: int
    def __init__(self, workflow_runs: _Optional[_Iterable[_Union[WorkflowRunSummary, _Mapping]]] = ..., next_page: _Optional[int] = ...) -> None: ...

class ListWorkflowRunsByWorkflowRequest(_message.Message):
    __slots__ = ("owner", "repository", "workflow_id", "per_page", "page", "branch", "event")
    OWNER_FIELD_NUMBER: _ClassVar[int]
    REPOSITORY_FIELD_NUMBER: _ClassVar[int]
    WORKFLOW_ID_FIELD_NUMBER: _ClassVar[int]
    PER_PAGE_FIELD_NUMBER: _ClassVar[int]
    PAGE_FIELD_NUMBER: _ClassVar[int]
    BRANCH_FIELD_NUMBER: _ClassVar[int]
    EVENT_FIELD_NUMBER: _ClassVar[int]
    owner: str
    repository: str
    workflow_id: str
    per_page: int
    page: int
    branch: str
    event: str
    def __init__(self, owner: _Optional[str] = ..., repository: _Optional[str] = ..., workflow_id: _Optional[str] = ..., per_page: _Optional[int] = ..., page: _Optional[int] = ..., branch: _Optional[str] = ..., event: _Optional[str] = ...) -> None: ...

class ListWorkflowRunsByWorkflowResponse(_message.Message):
    __slots__ = ("workflow_runs", "next_page")
    WORKFLOW_RUNS_FIELD_NUMBER: _ClassVar[int]
    NEXT_PAGE_FIELD_NUMBER: _ClassVar[int]
    workflow_runs: _containers.RepeatedCompositeFieldContainer[WorkflowRunSummary]
    next_page: int
    def __init__(self, workflow_runs: _Optional[_Iterable[_Union[WorkflowRunSummary, _Mapping]]] = ..., next_page: _Optional[int] = ...) -> None: ...

class GetRunAttemptRequest(_message.Message):
    __slots__ = ("owner", "repository", "run_id", "attempt_number")
    OWNER_FIELD_NUMBER: _ClassVar[int]
    REPOSITORY_FIELD_NUMBER: _ClassVar[int]
    RUN_ID_FIELD_NUMBER: _ClassVar[int]
    ATTEMPT_NUMBER_FIELD_NUMBER: _ClassVar[int]
    owner: str
    repository: str
    run_id: int
    attempt_number: int
    def __init__(self, owner: _Optional[str] = ..., repository: _Optional[str] = ..., run_id: _Optional[int] = ..., attempt_number: _Optional[int] = ...) -> None: ...

class GetRunAttemptResponse(_message.Message):
    __slots__ = ("id", "name", "head_branch", "head_sha", "status", "conclusion", "run_attempt", "html_url")
    ID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    HEAD_BRANCH_FIELD_NUMBER: _ClassVar[int]
    HEAD_SHA_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    CONCLUSION_FIELD_NUMBER: _ClassVar[int]
    RUN_ATTEMPT_FIELD_NUMBER: _ClassVar[int]
    HTML_URL_FIELD_NUMBER: _ClassVar[int]
    id: int
    name: str
    head_branch: str
    head_sha: str
    status: str
    conclusion: str
    run_attempt: int
    html_url: str
    def __init__(self, id: _Optional[int] = ..., name: _Optional[str] = ..., head_branch: _Optional[str] = ..., head_sha: _Optional[str] = ..., status: _Optional[str] = ..., conclusion: _Optional[str] = ..., run_attempt: _Optional[int] = ..., html_url: _Optional[str] = ...) -> None: ...

class ListRunJobsRequest(_message.Message):
    __slots__ = ("owner", "repository", "run_id", "per_page", "page")
    OWNER_FIELD_NUMBER: _ClassVar[int]
    REPOSITORY_FIELD_NUMBER: _ClassVar[int]
    RUN_ID_FIELD_NUMBER: _ClassVar[int]
    PER_PAGE_FIELD_NUMBER: _ClassVar[int]
    PAGE_FIELD_NUMBER: _ClassVar[int]
    owner: str
    repository: str
    run_id: int
    per_page: int
    page: int
    def __init__(self, owner: _Optional[str] = ..., repository: _Optional[str] = ..., run_id: _Optional[int] = ..., per_page: _Optional[int] = ..., page: _Optional[int] = ...) -> None: ...

class RunJobSummary(_message.Message):
    __slots__ = ("id", "name", "status", "conclusion", "started_at", "completed_at", "html_url", "runner_name")
    ID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    CONCLUSION_FIELD_NUMBER: _ClassVar[int]
    STARTED_AT_FIELD_NUMBER: _ClassVar[int]
    COMPLETED_AT_FIELD_NUMBER: _ClassVar[int]
    HTML_URL_FIELD_NUMBER: _ClassVar[int]
    RUNNER_NAME_FIELD_NUMBER: _ClassVar[int]
    id: int
    name: str
    status: str
    conclusion: str
    started_at: str
    completed_at: str
    html_url: str
    runner_name: str
    def __init__(self, id: _Optional[int] = ..., name: _Optional[str] = ..., status: _Optional[str] = ..., conclusion: _Optional[str] = ..., started_at: _Optional[str] = ..., completed_at: _Optional[str] = ..., html_url: _Optional[str] = ..., runner_name: _Optional[str] = ...) -> None: ...

class ListRunJobsResponse(_message.Message):
    __slots__ = ("jobs", "next_page")
    JOBS_FIELD_NUMBER: _ClassVar[int]
    NEXT_PAGE_FIELD_NUMBER: _ClassVar[int]
    jobs: _containers.RepeatedCompositeFieldContainer[RunJobSummary]
    next_page: int
    def __init__(self, jobs: _Optional[_Iterable[_Union[RunJobSummary, _Mapping]]] = ..., next_page: _Optional[int] = ...) -> None: ...
