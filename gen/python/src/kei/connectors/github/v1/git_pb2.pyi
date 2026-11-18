from kei.contracts.v1 import options_pb2 as _options_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class GetRefRequest(_message.Message):
    __slots__ = ("owner", "repository", "ref")
    OWNER_FIELD_NUMBER: _ClassVar[int]
    REPOSITORY_FIELD_NUMBER: _ClassVar[int]
    REF_FIELD_NUMBER: _ClassVar[int]
    owner: str
    repository: str
    ref: str
    def __init__(self, owner: _Optional[str] = ..., repository: _Optional[str] = ..., ref: _Optional[str] = ...) -> None: ...

class GitObject(_message.Message):
    __slots__ = ("sha", "type")
    SHA_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    sha: str
    type: str
    def __init__(self, sha: _Optional[str] = ..., type: _Optional[str] = ...) -> None: ...

class GetRefResponse(_message.Message):
    __slots__ = ("ref", "object")
    REF_FIELD_NUMBER: _ClassVar[int]
    OBJECT_FIELD_NUMBER: _ClassVar[int]
    ref: str
    object: GitObject
    def __init__(self, ref: _Optional[str] = ..., object: _Optional[_Union[GitObject, _Mapping]] = ...) -> None: ...

class CreateRefRequest(_message.Message):
    __slots__ = ("owner", "repository", "ref", "sha")
    OWNER_FIELD_NUMBER: _ClassVar[int]
    REPOSITORY_FIELD_NUMBER: _ClassVar[int]
    REF_FIELD_NUMBER: _ClassVar[int]
    SHA_FIELD_NUMBER: _ClassVar[int]
    owner: str
    repository: str
    ref: str
    sha: str
    def __init__(self, owner: _Optional[str] = ..., repository: _Optional[str] = ..., ref: _Optional[str] = ..., sha: _Optional[str] = ...) -> None: ...

class CreateRefResponse(_message.Message):
    __slots__ = ("ref", "object")
    REF_FIELD_NUMBER: _ClassVar[int]
    OBJECT_FIELD_NUMBER: _ClassVar[int]
    ref: str
    object: GitObject
    def __init__(self, ref: _Optional[str] = ..., object: _Optional[_Union[GitObject, _Mapping]] = ...) -> None: ...

class UpdateRefRequest(_message.Message):
    __slots__ = ("owner", "repository", "ref", "sha", "force")
    OWNER_FIELD_NUMBER: _ClassVar[int]
    REPOSITORY_FIELD_NUMBER: _ClassVar[int]
    REF_FIELD_NUMBER: _ClassVar[int]
    SHA_FIELD_NUMBER: _ClassVar[int]
    FORCE_FIELD_NUMBER: _ClassVar[int]
    owner: str
    repository: str
    ref: str
    sha: str
    force: bool
    def __init__(self, owner: _Optional[str] = ..., repository: _Optional[str] = ..., ref: _Optional[str] = ..., sha: _Optional[str] = ..., force: _Optional[bool] = ...) -> None: ...

class UpdateRefResponse(_message.Message):
    __slots__ = ("ref", "object")
    REF_FIELD_NUMBER: _ClassVar[int]
    OBJECT_FIELD_NUMBER: _ClassVar[int]
    ref: str
    object: GitObject
    def __init__(self, ref: _Optional[str] = ..., object: _Optional[_Union[GitObject, _Mapping]] = ...) -> None: ...

class CreateBlobRequest(_message.Message):
    __slots__ = ("owner", "repository", "content", "encoding")
    OWNER_FIELD_NUMBER: _ClassVar[int]
    REPOSITORY_FIELD_NUMBER: _ClassVar[int]
    CONTENT_FIELD_NUMBER: _ClassVar[int]
    ENCODING_FIELD_NUMBER: _ClassVar[int]
    owner: str
    repository: str
    content: str
    encoding: str
    def __init__(self, owner: _Optional[str] = ..., repository: _Optional[str] = ..., content: _Optional[str] = ..., encoding: _Optional[str] = ...) -> None: ...

class CreateBlobResponse(_message.Message):
    __slots__ = ("sha", "url")
    SHA_FIELD_NUMBER: _ClassVar[int]
    URL_FIELD_NUMBER: _ClassVar[int]
    sha: str
    url: str
    def __init__(self, sha: _Optional[str] = ..., url: _Optional[str] = ...) -> None: ...

class TreeEntry(_message.Message):
    __slots__ = ("path", "mode", "type", "sha")
    PATH_FIELD_NUMBER: _ClassVar[int]
    MODE_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    SHA_FIELD_NUMBER: _ClassVar[int]
    path: str
    mode: str
    type: str
    sha: str
    def __init__(self, path: _Optional[str] = ..., mode: _Optional[str] = ..., type: _Optional[str] = ..., sha: _Optional[str] = ...) -> None: ...

class CreateTreeRequest(_message.Message):
    __slots__ = ("owner", "repository", "tree", "base_tree")
    OWNER_FIELD_NUMBER: _ClassVar[int]
    REPOSITORY_FIELD_NUMBER: _ClassVar[int]
    TREE_FIELD_NUMBER: _ClassVar[int]
    BASE_TREE_FIELD_NUMBER: _ClassVar[int]
    owner: str
    repository: str
    tree: _containers.RepeatedCompositeFieldContainer[TreeEntry]
    base_tree: str
    def __init__(self, owner: _Optional[str] = ..., repository: _Optional[str] = ..., tree: _Optional[_Iterable[_Union[TreeEntry, _Mapping]]] = ..., base_tree: _Optional[str] = ...) -> None: ...

class CreateTreeResponse(_message.Message):
    __slots__ = ("sha", "url", "tree")
    SHA_FIELD_NUMBER: _ClassVar[int]
    URL_FIELD_NUMBER: _ClassVar[int]
    TREE_FIELD_NUMBER: _ClassVar[int]
    sha: str
    url: str
    tree: _containers.RepeatedCompositeFieldContainer[TreeEntry]
    def __init__(self, sha: _Optional[str] = ..., url: _Optional[str] = ..., tree: _Optional[_Iterable[_Union[TreeEntry, _Mapping]]] = ...) -> None: ...

class CreateCommitRequest(_message.Message):
    __slots__ = ("owner", "repository", "message", "tree", "parents")
    OWNER_FIELD_NUMBER: _ClassVar[int]
    REPOSITORY_FIELD_NUMBER: _ClassVar[int]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    TREE_FIELD_NUMBER: _ClassVar[int]
    PARENTS_FIELD_NUMBER: _ClassVar[int]
    owner: str
    repository: str
    message: str
    tree: str
    parents: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, owner: _Optional[str] = ..., repository: _Optional[str] = ..., message: _Optional[str] = ..., tree: _Optional[str] = ..., parents: _Optional[_Iterable[str]] = ...) -> None: ...

class CreateCommitResponse(_message.Message):
    __slots__ = ("sha", "html_url", "message", "tree", "parents")
    SHA_FIELD_NUMBER: _ClassVar[int]
    HTML_URL_FIELD_NUMBER: _ClassVar[int]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    TREE_FIELD_NUMBER: _ClassVar[int]
    PARENTS_FIELD_NUMBER: _ClassVar[int]
    sha: str
    html_url: str
    message: str
    tree: GitObject
    parents: _containers.RepeatedCompositeFieldContainer[GitObject]
    def __init__(self, sha: _Optional[str] = ..., html_url: _Optional[str] = ..., message: _Optional[str] = ..., tree: _Optional[_Union[GitObject, _Mapping]] = ..., parents: _Optional[_Iterable[_Union[GitObject, _Mapping]]] = ...) -> None: ...
