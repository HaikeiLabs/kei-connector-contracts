from kei.contracts.v1 import options_pb2 as _options_pb2
from kei.connectors.github.v1 import contents_pb2 as _contents_pb2
from kei.connectors.github.v1 import pulls_pb2 as _pulls_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ListCommitsRequest(_message.Message):
    __slots__ = ("owner", "repository", "sha", "per_page", "page")
    OWNER_FIELD_NUMBER: _ClassVar[int]
    REPOSITORY_FIELD_NUMBER: _ClassVar[int]
    SHA_FIELD_NUMBER: _ClassVar[int]
    PER_PAGE_FIELD_NUMBER: _ClassVar[int]
    PAGE_FIELD_NUMBER: _ClassVar[int]
    owner: str
    repository: str
    sha: str
    per_page: int
    page: int
    def __init__(self, owner: _Optional[str] = ..., repository: _Optional[str] = ..., sha: _Optional[str] = ..., per_page: _Optional[int] = ..., page: _Optional[int] = ...) -> None: ...

class ListCommitsResponse(_message.Message):
    __slots__ = ("commits", "next_page")
    COMMITS_FIELD_NUMBER: _ClassVar[int]
    NEXT_PAGE_FIELD_NUMBER: _ClassVar[int]
    commits: _containers.RepeatedCompositeFieldContainer[_contents_pb2.CommitSummary]
    next_page: int
    def __init__(self, commits: _Optional[_Iterable[_Union[_contents_pb2.CommitSummary, _Mapping]]] = ..., next_page: _Optional[int] = ...) -> None: ...

class GetCommitRequest(_message.Message):
    __slots__ = ("owner", "repository", "ref")
    OWNER_FIELD_NUMBER: _ClassVar[int]
    REPOSITORY_FIELD_NUMBER: _ClassVar[int]
    REF_FIELD_NUMBER: _ClassVar[int]
    owner: str
    repository: str
    ref: str
    def __init__(self, owner: _Optional[str] = ..., repository: _Optional[str] = ..., ref: _Optional[str] = ...) -> None: ...

class CommitDetail(_message.Message):
    __slots__ = ("message",)
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    message: str
    def __init__(self, message: _Optional[str] = ...) -> None: ...

class GetCommitResponse(_message.Message):
    __slots__ = ("sha", "html_url", "commit", "author")
    SHA_FIELD_NUMBER: _ClassVar[int]
    HTML_URL_FIELD_NUMBER: _ClassVar[int]
    COMMIT_FIELD_NUMBER: _ClassVar[int]
    AUTHOR_FIELD_NUMBER: _ClassVar[int]
    sha: str
    html_url: str
    commit: CommitDetail
    author: _pulls_pb2.User
    def __init__(self, sha: _Optional[str] = ..., html_url: _Optional[str] = ..., commit: _Optional[_Union[CommitDetail, _Mapping]] = ..., author: _Optional[_Union[_pulls_pb2.User, _Mapping]] = ...) -> None: ...
