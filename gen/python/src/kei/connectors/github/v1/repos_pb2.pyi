from kei.contracts.v1 import options_pb2 as _options_pb2
from kei.connectors.github.v1 import pulls_pb2 as _pulls_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class GetRepositoryRequest(_message.Message):
    __slots__ = ("owner", "repository")
    OWNER_FIELD_NUMBER: _ClassVar[int]
    REPOSITORY_FIELD_NUMBER: _ClassVar[int]
    owner: str
    repository: str
    def __init__(self, owner: _Optional[str] = ..., repository: _Optional[str] = ...) -> None: ...

class GetRepositoryResponse(_message.Message):
    __slots__ = ("id", "name", "full_name", "html_url", "private", "description", "default_branch", "owner")
    ID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    FULL_NAME_FIELD_NUMBER: _ClassVar[int]
    HTML_URL_FIELD_NUMBER: _ClassVar[int]
    PRIVATE_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    DEFAULT_BRANCH_FIELD_NUMBER: _ClassVar[int]
    OWNER_FIELD_NUMBER: _ClassVar[int]
    id: int
    name: str
    full_name: str
    html_url: str
    private: bool
    description: str
    default_branch: str
    owner: _pulls_pb2.User
    def __init__(self, id: _Optional[int] = ..., name: _Optional[str] = ..., full_name: _Optional[str] = ..., html_url: _Optional[str] = ..., private: _Optional[bool] = ..., description: _Optional[str] = ..., default_branch: _Optional[str] = ..., owner: _Optional[_Union[_pulls_pb2.User, _Mapping]] = ...) -> None: ...
