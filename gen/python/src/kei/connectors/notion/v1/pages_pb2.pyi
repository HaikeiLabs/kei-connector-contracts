from google.protobuf import struct_pb2 as _struct_pb2
from kei.contracts.v1 import options_pb2 as _options_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

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

class CreatePageRequest(_message.Message):
    __slots__ = ("parent", "properties", "children")
    PARENT_FIELD_NUMBER: _ClassVar[int]
    PROPERTIES_FIELD_NUMBER: _ClassVar[int]
    CHILDREN_FIELD_NUMBER: _ClassVar[int]
    parent: _struct_pb2.Struct
    properties: _struct_pb2.Struct
    children: _containers.RepeatedCompositeFieldContainer[_struct_pb2.Struct]
    def __init__(self, parent: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ..., properties: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ..., children: _Optional[_Iterable[_Union[_struct_pb2.Struct, _Mapping]]] = ...) -> None: ...

class CreatePageResponse(_message.Message):
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

class UpdatePageRequest(_message.Message):
    __slots__ = ("page_id", "properties", "archived", "in_trash")
    PAGE_ID_FIELD_NUMBER: _ClassVar[int]
    PROPERTIES_FIELD_NUMBER: _ClassVar[int]
    ARCHIVED_FIELD_NUMBER: _ClassVar[int]
    IN_TRASH_FIELD_NUMBER: _ClassVar[int]
    page_id: str
    properties: _struct_pb2.Struct
    archived: bool
    in_trash: bool
    def __init__(self, page_id: _Optional[str] = ..., properties: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ..., archived: _Optional[bool] = ..., in_trash: _Optional[bool] = ...) -> None: ...

class UpdatePageResponse(_message.Message):
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
