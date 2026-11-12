from kei.contracts.v1 import options_pb2 as _options_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class PutContentsRequest(_message.Message):
    __slots__ = ("owner", "repository", "path", "message", "content", "branch", "sha")
    OWNER_FIELD_NUMBER: _ClassVar[int]
    REPOSITORY_FIELD_NUMBER: _ClassVar[int]
    PATH_FIELD_NUMBER: _ClassVar[int]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    CONTENT_FIELD_NUMBER: _ClassVar[int]
    BRANCH_FIELD_NUMBER: _ClassVar[int]
    SHA_FIELD_NUMBER: _ClassVar[int]
    owner: str
    repository: str
    path: str
    message: str
    content: str
    branch: str
    sha: str
    def __init__(self, owner: _Optional[str] = ..., repository: _Optional[str] = ..., path: _Optional[str] = ..., message: _Optional[str] = ..., content: _Optional[str] = ..., branch: _Optional[str] = ..., sha: _Optional[str] = ...) -> None: ...

class ContentFile(_message.Message):
    __slots__ = ("name", "path", "sha", "html_url", "download_url")
    NAME_FIELD_NUMBER: _ClassVar[int]
    PATH_FIELD_NUMBER: _ClassVar[int]
    SHA_FIELD_NUMBER: _ClassVar[int]
    HTML_URL_FIELD_NUMBER: _ClassVar[int]
    DOWNLOAD_URL_FIELD_NUMBER: _ClassVar[int]
    name: str
    path: str
    sha: str
    html_url: str
    download_url: str
    def __init__(self, name: _Optional[str] = ..., path: _Optional[str] = ..., sha: _Optional[str] = ..., html_url: _Optional[str] = ..., download_url: _Optional[str] = ...) -> None: ...

class CommitSummary(_message.Message):
    __slots__ = ("sha", "html_url")
    SHA_FIELD_NUMBER: _ClassVar[int]
    HTML_URL_FIELD_NUMBER: _ClassVar[int]
    sha: str
    html_url: str
    def __init__(self, sha: _Optional[str] = ..., html_url: _Optional[str] = ...) -> None: ...

class PutContentsResponse(_message.Message):
    __slots__ = ("content", "commit")
    CONTENT_FIELD_NUMBER: _ClassVar[int]
    COMMIT_FIELD_NUMBER: _ClassVar[int]
    content: ContentFile
    commit: CommitSummary
    def __init__(self, content: _Optional[_Union[ContentFile, _Mapping]] = ..., commit: _Optional[_Union[CommitSummary, _Mapping]] = ...) -> None: ...
