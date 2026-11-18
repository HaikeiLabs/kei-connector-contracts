from google.protobuf import struct_pb2 as _struct_pb2
from kei.contracts.v1 import options_pb2 as _options_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ListRecordsRequest(_message.Message):
    __slots__ = ("base_id", "table_id_or_name", "offset", "page_size")
    BASE_ID_FIELD_NUMBER: _ClassVar[int]
    TABLE_ID_OR_NAME_FIELD_NUMBER: _ClassVar[int]
    OFFSET_FIELD_NUMBER: _ClassVar[int]
    PAGE_SIZE_FIELD_NUMBER: _ClassVar[int]
    base_id: str
    table_id_or_name: str
    offset: str
    page_size: int
    def __init__(self, base_id: _Optional[str] = ..., table_id_or_name: _Optional[str] = ..., offset: _Optional[str] = ..., page_size: _Optional[int] = ...) -> None: ...

class RecordSummary(_message.Message):
    __slots__ = ("id", "fields", "created_time")
    ID_FIELD_NUMBER: _ClassVar[int]
    FIELDS_FIELD_NUMBER: _ClassVar[int]
    CREATED_TIME_FIELD_NUMBER: _ClassVar[int]
    id: str
    fields: _struct_pb2.Struct
    created_time: str
    def __init__(self, id: _Optional[str] = ..., fields: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ..., created_time: _Optional[str] = ...) -> None: ...

class ListRecordsResponse(_message.Message):
    __slots__ = ("records", "offset")
    RECORDS_FIELD_NUMBER: _ClassVar[int]
    OFFSET_FIELD_NUMBER: _ClassVar[int]
    records: _containers.RepeatedCompositeFieldContainer[RecordSummary]
    offset: str
    def __init__(self, records: _Optional[_Iterable[_Union[RecordSummary, _Mapping]]] = ..., offset: _Optional[str] = ...) -> None: ...

class GetRecordRequest(_message.Message):
    __slots__ = ("base_id", "table_id_or_name", "record_id")
    BASE_ID_FIELD_NUMBER: _ClassVar[int]
    TABLE_ID_OR_NAME_FIELD_NUMBER: _ClassVar[int]
    RECORD_ID_FIELD_NUMBER: _ClassVar[int]
    base_id: str
    table_id_or_name: str
    record_id: str
    def __init__(self, base_id: _Optional[str] = ..., table_id_or_name: _Optional[str] = ..., record_id: _Optional[str] = ...) -> None: ...

class GetRecordResponse(_message.Message):
    __slots__ = ("id", "fields", "created_time")
    ID_FIELD_NUMBER: _ClassVar[int]
    FIELDS_FIELD_NUMBER: _ClassVar[int]
    CREATED_TIME_FIELD_NUMBER: _ClassVar[int]
    id: str
    fields: _struct_pb2.Struct
    created_time: str
    def __init__(self, id: _Optional[str] = ..., fields: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ..., created_time: _Optional[str] = ...) -> None: ...

class CreateRecordRequest(_message.Message):
    __slots__ = ("base_id", "table_id_or_name", "fields")
    BASE_ID_FIELD_NUMBER: _ClassVar[int]
    TABLE_ID_OR_NAME_FIELD_NUMBER: _ClassVar[int]
    FIELDS_FIELD_NUMBER: _ClassVar[int]
    base_id: str
    table_id_or_name: str
    fields: _struct_pb2.Struct
    def __init__(self, base_id: _Optional[str] = ..., table_id_or_name: _Optional[str] = ..., fields: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ...) -> None: ...

class CreateRecordResponse(_message.Message):
    __slots__ = ("id", "fields", "created_time")
    ID_FIELD_NUMBER: _ClassVar[int]
    FIELDS_FIELD_NUMBER: _ClassVar[int]
    CREATED_TIME_FIELD_NUMBER: _ClassVar[int]
    id: str
    fields: _struct_pb2.Struct
    created_time: str
    def __init__(self, id: _Optional[str] = ..., fields: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ..., created_time: _Optional[str] = ...) -> None: ...

class UpdateRecordRequest(_message.Message):
    __slots__ = ("base_id", "table_id_or_name", "record_id", "fields")
    BASE_ID_FIELD_NUMBER: _ClassVar[int]
    TABLE_ID_OR_NAME_FIELD_NUMBER: _ClassVar[int]
    RECORD_ID_FIELD_NUMBER: _ClassVar[int]
    FIELDS_FIELD_NUMBER: _ClassVar[int]
    base_id: str
    table_id_or_name: str
    record_id: str
    fields: _struct_pb2.Struct
    def __init__(self, base_id: _Optional[str] = ..., table_id_or_name: _Optional[str] = ..., record_id: _Optional[str] = ..., fields: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ...) -> None: ...

class UpdateRecordResponse(_message.Message):
    __slots__ = ("id", "fields", "created_time")
    ID_FIELD_NUMBER: _ClassVar[int]
    FIELDS_FIELD_NUMBER: _ClassVar[int]
    CREATED_TIME_FIELD_NUMBER: _ClassVar[int]
    id: str
    fields: _struct_pb2.Struct
    created_time: str
    def __init__(self, id: _Optional[str] = ..., fields: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ..., created_time: _Optional[str] = ...) -> None: ...
