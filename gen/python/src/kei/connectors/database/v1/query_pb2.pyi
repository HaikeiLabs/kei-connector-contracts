from google.protobuf import struct_pb2 as _struct_pb2
from kei.contracts.v1 import options_pb2 as _options_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class Dialect(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    DIALECT_UNSPECIFIED: _ClassVar[Dialect]
    DIALECT_POSTGRES: _ClassVar[Dialect]
    DIALECT_MYSQL: _ClassVar[Dialect]
    DIALECT_MSSQL: _ClassVar[Dialect]
    DIALECT_SNOWFLAKE: _ClassVar[Dialect]
    DIALECT_SUPABASE: _ClassVar[Dialect]
DIALECT_UNSPECIFIED: Dialect
DIALECT_POSTGRES: Dialect
DIALECT_MYSQL: Dialect
DIALECT_MSSQL: Dialect
DIALECT_SNOWFLAKE: Dialect
DIALECT_SUPABASE: Dialect

class RunQueryRequest(_message.Message):
    __slots__ = ("dialect", "sql", "params", "max_rows")
    DIALECT_FIELD_NUMBER: _ClassVar[int]
    SQL_FIELD_NUMBER: _ClassVar[int]
    PARAMS_FIELD_NUMBER: _ClassVar[int]
    MAX_ROWS_FIELD_NUMBER: _ClassVar[int]
    dialect: Dialect
    sql: str
    params: _containers.RepeatedCompositeFieldContainer[_struct_pb2.Value]
    max_rows: int
    def __init__(self, dialect: _Optional[_Union[Dialect, str]] = ..., sql: _Optional[str] = ..., params: _Optional[_Iterable[_Union[_struct_pb2.Value, _Mapping]]] = ..., max_rows: _Optional[int] = ...) -> None: ...

class Column(_message.Message):
    __slots__ = ("name", "type")
    NAME_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    name: str
    type: str
    def __init__(self, name: _Optional[str] = ..., type: _Optional[str] = ...) -> None: ...

class RunQueryResponse(_message.Message):
    __slots__ = ("columns", "rows", "truncated")
    COLUMNS_FIELD_NUMBER: _ClassVar[int]
    ROWS_FIELD_NUMBER: _ClassVar[int]
    TRUNCATED_FIELD_NUMBER: _ClassVar[int]
    columns: _containers.RepeatedCompositeFieldContainer[Column]
    rows: _containers.RepeatedCompositeFieldContainer[_struct_pb2.ListValue]
    truncated: bool
    def __init__(self, columns: _Optional[_Iterable[_Union[Column, _Mapping]]] = ..., rows: _Optional[_Iterable[_Union[_struct_pb2.ListValue, _Mapping]]] = ..., truncated: _Optional[bool] = ...) -> None: ...
