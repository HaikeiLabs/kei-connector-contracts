from google.protobuf import struct_pb2 as _struct_pb2
from kei.contracts.v1 import options_pb2 as _options_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class GetContactRequest(_message.Message):
    __slots__ = ("contact_id",)
    CONTACT_ID_FIELD_NUMBER: _ClassVar[int]
    contact_id: str
    def __init__(self, contact_id: _Optional[str] = ...) -> None: ...

class GetContactResponse(_message.Message):
    __slots__ = ("id", "properties", "created_at", "updated_at", "archived")
    ID_FIELD_NUMBER: _ClassVar[int]
    PROPERTIES_FIELD_NUMBER: _ClassVar[int]
    CREATED_AT_FIELD_NUMBER: _ClassVar[int]
    UPDATED_AT_FIELD_NUMBER: _ClassVar[int]
    ARCHIVED_FIELD_NUMBER: _ClassVar[int]
    id: str
    properties: _struct_pb2.Struct
    created_at: str
    updated_at: str
    archived: bool
    def __init__(self, id: _Optional[str] = ..., properties: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ..., created_at: _Optional[str] = ..., updated_at: _Optional[str] = ..., archived: _Optional[bool] = ...) -> None: ...

class ListContactsRequest(_message.Message):
    __slots__ = ("after", "limit")
    AFTER_FIELD_NUMBER: _ClassVar[int]
    LIMIT_FIELD_NUMBER: _ClassVar[int]
    after: str
    limit: int
    def __init__(self, after: _Optional[str] = ..., limit: _Optional[int] = ...) -> None: ...

class ContactSummary(_message.Message):
    __slots__ = ("id", "properties", "created_at", "updated_at", "archived")
    ID_FIELD_NUMBER: _ClassVar[int]
    PROPERTIES_FIELD_NUMBER: _ClassVar[int]
    CREATED_AT_FIELD_NUMBER: _ClassVar[int]
    UPDATED_AT_FIELD_NUMBER: _ClassVar[int]
    ARCHIVED_FIELD_NUMBER: _ClassVar[int]
    id: str
    properties: _struct_pb2.Struct
    created_at: str
    updated_at: str
    archived: bool
    def __init__(self, id: _Optional[str] = ..., properties: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ..., created_at: _Optional[str] = ..., updated_at: _Optional[str] = ..., archived: _Optional[bool] = ...) -> None: ...

class ListContactsResponse(_message.Message):
    __slots__ = ("results", "next_after")
    RESULTS_FIELD_NUMBER: _ClassVar[int]
    NEXT_AFTER_FIELD_NUMBER: _ClassVar[int]
    results: _containers.RepeatedCompositeFieldContainer[ContactSummary]
    next_after: str
    def __init__(self, results: _Optional[_Iterable[_Union[ContactSummary, _Mapping]]] = ..., next_after: _Optional[str] = ...) -> None: ...
