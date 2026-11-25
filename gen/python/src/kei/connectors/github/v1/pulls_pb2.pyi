from kei.contracts.v1 import options_pb2 as _options_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class GetPullRequestRequest(_message.Message):
    __slots__ = ("owner", "repository", "pull_number")
    OWNER_FIELD_NUMBER: _ClassVar[int]
    REPOSITORY_FIELD_NUMBER: _ClassVar[int]
    PULL_NUMBER_FIELD_NUMBER: _ClassVar[int]
    owner: str
    repository: str
    pull_number: int
    def __init__(self, owner: _Optional[str] = ..., repository: _Optional[str] = ..., pull_number: _Optional[int] = ...) -> None: ...

class GetPullRequestResponse(_message.Message):
    __slots__ = ("id", "number", "title", "state", "body", "html_url", "draft", "merged_at", "head", "base", "author", "labels", "diff_url", "patch_url")
    ID_FIELD_NUMBER: _ClassVar[int]
    NUMBER_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    STATE_FIELD_NUMBER: _ClassVar[int]
    BODY_FIELD_NUMBER: _ClassVar[int]
    HTML_URL_FIELD_NUMBER: _ClassVar[int]
    DRAFT_FIELD_NUMBER: _ClassVar[int]
    MERGED_AT_FIELD_NUMBER: _ClassVar[int]
    HEAD_FIELD_NUMBER: _ClassVar[int]
    BASE_FIELD_NUMBER: _ClassVar[int]
    AUTHOR_FIELD_NUMBER: _ClassVar[int]
    LABELS_FIELD_NUMBER: _ClassVar[int]
    DIFF_URL_FIELD_NUMBER: _ClassVar[int]
    PATCH_URL_FIELD_NUMBER: _ClassVar[int]
    id: int
    number: int
    title: str
    state: str
    body: str
    html_url: str
    draft: bool
    merged_at: str
    head: PullRef
    base: PullRef
    author: User
    labels: _containers.RepeatedCompositeFieldContainer[Label]
    diff_url: str
    patch_url: str
    def __init__(self, id: _Optional[int] = ..., number: _Optional[int] = ..., title: _Optional[str] = ..., state: _Optional[str] = ..., body: _Optional[str] = ..., html_url: _Optional[str] = ..., draft: _Optional[bool] = ..., merged_at: _Optional[str] = ..., head: _Optional[_Union[PullRef, _Mapping]] = ..., base: _Optional[_Union[PullRef, _Mapping]] = ..., author: _Optional[_Union[User, _Mapping]] = ..., labels: _Optional[_Iterable[_Union[Label, _Mapping]]] = ..., diff_url: _Optional[str] = ..., patch_url: _Optional[str] = ...) -> None: ...

class PullRef(_message.Message):
    __slots__ = ("ref",)
    REF_FIELD_NUMBER: _ClassVar[int]
    ref: str
    def __init__(self, ref: _Optional[str] = ...) -> None: ...

class User(_message.Message):
    __slots__ = ("id", "login")
    ID_FIELD_NUMBER: _ClassVar[int]
    LOGIN_FIELD_NUMBER: _ClassVar[int]
    id: int
    login: str
    def __init__(self, id: _Optional[int] = ..., login: _Optional[str] = ...) -> None: ...

class Label(_message.Message):
    __slots__ = ("id", "name")
    ID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    id: int
    name: str
    def __init__(self, id: _Optional[int] = ..., name: _Optional[str] = ...) -> None: ...

class ListPullsRequest(_message.Message):
    __slots__ = ("owner", "repository", "per_page", "page", "state", "head", "base")
    OWNER_FIELD_NUMBER: _ClassVar[int]
    REPOSITORY_FIELD_NUMBER: _ClassVar[int]
    PER_PAGE_FIELD_NUMBER: _ClassVar[int]
    PAGE_FIELD_NUMBER: _ClassVar[int]
    STATE_FIELD_NUMBER: _ClassVar[int]
    HEAD_FIELD_NUMBER: _ClassVar[int]
    BASE_FIELD_NUMBER: _ClassVar[int]
    owner: str
    repository: str
    per_page: int
    page: int
    state: str
    head: str
    base: str
    def __init__(self, owner: _Optional[str] = ..., repository: _Optional[str] = ..., per_page: _Optional[int] = ..., page: _Optional[int] = ..., state: _Optional[str] = ..., head: _Optional[str] = ..., base: _Optional[str] = ...) -> None: ...

class PullRequestSummary(_message.Message):
    __slots__ = ("id", "number", "title", "state", "body", "html_url", "draft", "head", "base", "author", "labels")
    ID_FIELD_NUMBER: _ClassVar[int]
    NUMBER_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    STATE_FIELD_NUMBER: _ClassVar[int]
    BODY_FIELD_NUMBER: _ClassVar[int]
    HTML_URL_FIELD_NUMBER: _ClassVar[int]
    DRAFT_FIELD_NUMBER: _ClassVar[int]
    HEAD_FIELD_NUMBER: _ClassVar[int]
    BASE_FIELD_NUMBER: _ClassVar[int]
    AUTHOR_FIELD_NUMBER: _ClassVar[int]
    LABELS_FIELD_NUMBER: _ClassVar[int]
    id: int
    number: int
    title: str
    state: str
    body: str
    html_url: str
    draft: bool
    head: PullRef
    base: PullRef
    author: User
    labels: _containers.RepeatedCompositeFieldContainer[Label]
    def __init__(self, id: _Optional[int] = ..., number: _Optional[int] = ..., title: _Optional[str] = ..., state: _Optional[str] = ..., body: _Optional[str] = ..., html_url: _Optional[str] = ..., draft: _Optional[bool] = ..., head: _Optional[_Union[PullRef, _Mapping]] = ..., base: _Optional[_Union[PullRef, _Mapping]] = ..., author: _Optional[_Union[User, _Mapping]] = ..., labels: _Optional[_Iterable[_Union[Label, _Mapping]]] = ...) -> None: ...

class ListPullsResponse(_message.Message):
    __slots__ = ("pulls", "next_page")
    PULLS_FIELD_NUMBER: _ClassVar[int]
    NEXT_PAGE_FIELD_NUMBER: _ClassVar[int]
    pulls: _containers.RepeatedCompositeFieldContainer[PullRequestSummary]
    next_page: int
    def __init__(self, pulls: _Optional[_Iterable[_Union[PullRequestSummary, _Mapping]]] = ..., next_page: _Optional[int] = ...) -> None: ...

class CreatePullRequestRequest(_message.Message):
    __slots__ = ("owner", "repository", "title", "head", "base", "body", "draft")
    OWNER_FIELD_NUMBER: _ClassVar[int]
    REPOSITORY_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    HEAD_FIELD_NUMBER: _ClassVar[int]
    BASE_FIELD_NUMBER: _ClassVar[int]
    BODY_FIELD_NUMBER: _ClassVar[int]
    DRAFT_FIELD_NUMBER: _ClassVar[int]
    owner: str
    repository: str
    title: str
    head: str
    base: str
    body: str
    draft: bool
    def __init__(self, owner: _Optional[str] = ..., repository: _Optional[str] = ..., title: _Optional[str] = ..., head: _Optional[str] = ..., base: _Optional[str] = ..., body: _Optional[str] = ..., draft: _Optional[bool] = ...) -> None: ...

class CreatePullRequestResponse(_message.Message):
    __slots__ = ("id", "number", "title", "state", "html_url", "draft", "head", "base", "author", "labels", "diff_url", "patch_url")
    ID_FIELD_NUMBER: _ClassVar[int]
    NUMBER_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    STATE_FIELD_NUMBER: _ClassVar[int]
    HTML_URL_FIELD_NUMBER: _ClassVar[int]
    DRAFT_FIELD_NUMBER: _ClassVar[int]
    HEAD_FIELD_NUMBER: _ClassVar[int]
    BASE_FIELD_NUMBER: _ClassVar[int]
    AUTHOR_FIELD_NUMBER: _ClassVar[int]
    LABELS_FIELD_NUMBER: _ClassVar[int]
    DIFF_URL_FIELD_NUMBER: _ClassVar[int]
    PATCH_URL_FIELD_NUMBER: _ClassVar[int]
    id: int
    number: int
    title: str
    state: str
    html_url: str
    draft: bool
    head: PullRef
    base: PullRef
    author: User
    labels: _containers.RepeatedCompositeFieldContainer[Label]
    diff_url: str
    patch_url: str
    def __init__(self, id: _Optional[int] = ..., number: _Optional[int] = ..., title: _Optional[str] = ..., state: _Optional[str] = ..., html_url: _Optional[str] = ..., draft: _Optional[bool] = ..., head: _Optional[_Union[PullRef, _Mapping]] = ..., base: _Optional[_Union[PullRef, _Mapping]] = ..., author: _Optional[_Union[User, _Mapping]] = ..., labels: _Optional[_Iterable[_Union[Label, _Mapping]]] = ..., diff_url: _Optional[str] = ..., patch_url: _Optional[str] = ...) -> None: ...

class UpdatePullRequestRequest(_message.Message):
    __slots__ = ("owner", "repository", "pull_number", "title", "body", "state", "base")
    OWNER_FIELD_NUMBER: _ClassVar[int]
    REPOSITORY_FIELD_NUMBER: _ClassVar[int]
    PULL_NUMBER_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    BODY_FIELD_NUMBER: _ClassVar[int]
    STATE_FIELD_NUMBER: _ClassVar[int]
    BASE_FIELD_NUMBER: _ClassVar[int]
    owner: str
    repository: str
    pull_number: int
    title: str
    body: str
    state: str
    base: str
    def __init__(self, owner: _Optional[str] = ..., repository: _Optional[str] = ..., pull_number: _Optional[int] = ..., title: _Optional[str] = ..., body: _Optional[str] = ..., state: _Optional[str] = ..., base: _Optional[str] = ...) -> None: ...

class UpdatePullRequestResponse(_message.Message):
    __slots__ = ("id", "number", "title", "state", "body", "html_url", "draft", "merged_at", "head", "base", "author", "labels", "diff_url", "patch_url")
    ID_FIELD_NUMBER: _ClassVar[int]
    NUMBER_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    STATE_FIELD_NUMBER: _ClassVar[int]
    BODY_FIELD_NUMBER: _ClassVar[int]
    HTML_URL_FIELD_NUMBER: _ClassVar[int]
    DRAFT_FIELD_NUMBER: _ClassVar[int]
    MERGED_AT_FIELD_NUMBER: _ClassVar[int]
    HEAD_FIELD_NUMBER: _ClassVar[int]
    BASE_FIELD_NUMBER: _ClassVar[int]
    AUTHOR_FIELD_NUMBER: _ClassVar[int]
    LABELS_FIELD_NUMBER: _ClassVar[int]
    DIFF_URL_FIELD_NUMBER: _ClassVar[int]
    PATCH_URL_FIELD_NUMBER: _ClassVar[int]
    id: int
    number: int
    title: str
    state: str
    body: str
    html_url: str
    draft: bool
    merged_at: str
    head: PullRef
    base: PullRef
    author: User
    labels: _containers.RepeatedCompositeFieldContainer[Label]
    diff_url: str
    patch_url: str
    def __init__(self, id: _Optional[int] = ..., number: _Optional[int] = ..., title: _Optional[str] = ..., state: _Optional[str] = ..., body: _Optional[str] = ..., html_url: _Optional[str] = ..., draft: _Optional[bool] = ..., merged_at: _Optional[str] = ..., head: _Optional[_Union[PullRef, _Mapping]] = ..., base: _Optional[_Union[PullRef, _Mapping]] = ..., author: _Optional[_Union[User, _Mapping]] = ..., labels: _Optional[_Iterable[_Union[Label, _Mapping]]] = ..., diff_url: _Optional[str] = ..., patch_url: _Optional[str] = ...) -> None: ...

class ListPullReviewsRequest(_message.Message):
    __slots__ = ("owner", "repository", "pull_number", "per_page", "page")
    OWNER_FIELD_NUMBER: _ClassVar[int]
    REPOSITORY_FIELD_NUMBER: _ClassVar[int]
    PULL_NUMBER_FIELD_NUMBER: _ClassVar[int]
    PER_PAGE_FIELD_NUMBER: _ClassVar[int]
    PAGE_FIELD_NUMBER: _ClassVar[int]
    owner: str
    repository: str
    pull_number: int
    per_page: int
    page: int
    def __init__(self, owner: _Optional[str] = ..., repository: _Optional[str] = ..., pull_number: _Optional[int] = ..., per_page: _Optional[int] = ..., page: _Optional[int] = ...) -> None: ...

class ReviewSummary(_message.Message):
    __slots__ = ("id", "state", "body", "html_url", "submitted_at", "author")
    ID_FIELD_NUMBER: _ClassVar[int]
    STATE_FIELD_NUMBER: _ClassVar[int]
    BODY_FIELD_NUMBER: _ClassVar[int]
    HTML_URL_FIELD_NUMBER: _ClassVar[int]
    SUBMITTED_AT_FIELD_NUMBER: _ClassVar[int]
    AUTHOR_FIELD_NUMBER: _ClassVar[int]
    id: int
    state: str
    body: str
    html_url: str
    submitted_at: str
    author: User
    def __init__(self, id: _Optional[int] = ..., state: _Optional[str] = ..., body: _Optional[str] = ..., html_url: _Optional[str] = ..., submitted_at: _Optional[str] = ..., author: _Optional[_Union[User, _Mapping]] = ...) -> None: ...

class ListPullReviewsResponse(_message.Message):
    __slots__ = ("reviews", "next_page")
    REVIEWS_FIELD_NUMBER: _ClassVar[int]
    NEXT_PAGE_FIELD_NUMBER: _ClassVar[int]
    reviews: _containers.RepeatedCompositeFieldContainer[ReviewSummary]
    next_page: int
    def __init__(self, reviews: _Optional[_Iterable[_Union[ReviewSummary, _Mapping]]] = ..., next_page: _Optional[int] = ...) -> None: ...

class CreatePullReviewRequest(_message.Message):
    __slots__ = ("owner", "repository", "pull_number", "body", "event")
    OWNER_FIELD_NUMBER: _ClassVar[int]
    REPOSITORY_FIELD_NUMBER: _ClassVar[int]
    PULL_NUMBER_FIELD_NUMBER: _ClassVar[int]
    BODY_FIELD_NUMBER: _ClassVar[int]
    EVENT_FIELD_NUMBER: _ClassVar[int]
    owner: str
    repository: str
    pull_number: int
    body: str
    event: str
    def __init__(self, owner: _Optional[str] = ..., repository: _Optional[str] = ..., pull_number: _Optional[int] = ..., body: _Optional[str] = ..., event: _Optional[str] = ...) -> None: ...

class CreatePullReviewResponse(_message.Message):
    __slots__ = ("id", "state", "body", "html_url", "submitted_at", "author")
    ID_FIELD_NUMBER: _ClassVar[int]
    STATE_FIELD_NUMBER: _ClassVar[int]
    BODY_FIELD_NUMBER: _ClassVar[int]
    HTML_URL_FIELD_NUMBER: _ClassVar[int]
    SUBMITTED_AT_FIELD_NUMBER: _ClassVar[int]
    AUTHOR_FIELD_NUMBER: _ClassVar[int]
    id: int
    state: str
    body: str
    html_url: str
    submitted_at: str
    author: User
    def __init__(self, id: _Optional[int] = ..., state: _Optional[str] = ..., body: _Optional[str] = ..., html_url: _Optional[str] = ..., submitted_at: _Optional[str] = ..., author: _Optional[_Union[User, _Mapping]] = ...) -> None: ...
