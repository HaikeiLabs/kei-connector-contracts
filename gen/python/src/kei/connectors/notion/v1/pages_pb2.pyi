from kei.contracts.v1 import options_pb2 as _options_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class GetPageRequest(_message.Message):
    __slots__ = ("page_id",)
    PAGE_ID_FIELD_NUMBER: _ClassVar[int]
    page_id: str
    def __init__(self, page_id: _Optional[str] = ...) -> None: ...

class GetPageResponse(_message.Message):
    __slots__ = ("id", "created_time", "last_edited_time", "archived", "url", "public_url")
    ID_FIELD_NUMBER: _ClassVar[int]
    CREATED_TIME_FIELD_NUMBER: _ClassVar[int]
    LAST_EDITED_TIME_FIELD_NUMBER: _ClassVar[int]
    ARCHIVED_FIELD_NUMBER: _ClassVar[int]
    URL_FIELD_NUMBER: _ClassVar[int]
    PUBLIC_URL_FIELD_NUMBER: _ClassVar[int]
    id: str
    created_time: str
    last_edited_time: str
    archived: bool
    url: str
    public_url: str
    def __init__(self, id: _Optional[str] = ..., created_time: _Optional[str] = ..., last_edited_time: _Optional[str] = ..., archived: _Optional[bool] = ..., url: _Optional[str] = ..., public_url: _Optional[str] = ...) -> None: ...
