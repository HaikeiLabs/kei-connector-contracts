from google.protobuf import struct_pb2 as _struct_pb2
from kei.contracts.v1 import options_pb2 as _options_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ListDatasetsRequest(_message.Message):
    __slots__ = ("group_id",)
    GROUP_ID_FIELD_NUMBER: _ClassVar[int]
    group_id: str
    def __init__(self, group_id: _Optional[str] = ...) -> None: ...

class DatasetSummary(_message.Message):
    __slots__ = ("id", "name", "configured_by", "target_storage_mode")
    ID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    CONFIGURED_BY_FIELD_NUMBER: _ClassVar[int]
    TARGET_STORAGE_MODE_FIELD_NUMBER: _ClassVar[int]
    id: str
    name: str
    configured_by: str
    target_storage_mode: str
    def __init__(self, id: _Optional[str] = ..., name: _Optional[str] = ..., configured_by: _Optional[str] = ..., target_storage_mode: _Optional[str] = ...) -> None: ...

class ListDatasetsResponse(_message.Message):
    __slots__ = ("value",)
    VALUE_FIELD_NUMBER: _ClassVar[int]
    value: _containers.RepeatedCompositeFieldContainer[DatasetSummary]
    def __init__(self, value: _Optional[_Iterable[_Union[DatasetSummary, _Mapping]]] = ...) -> None: ...

class ExecuteQueriesRequest(_message.Message):
    __slots__ = ("group_id", "dataset_id", "queries")
    GROUP_ID_FIELD_NUMBER: _ClassVar[int]
    DATASET_ID_FIELD_NUMBER: _ClassVar[int]
    QUERIES_FIELD_NUMBER: _ClassVar[int]
    group_id: str
    dataset_id: str
    queries: _containers.RepeatedCompositeFieldContainer[QueryDefinition]
    def __init__(self, group_id: _Optional[str] = ..., dataset_id: _Optional[str] = ..., queries: _Optional[_Iterable[_Union[QueryDefinition, _Mapping]]] = ...) -> None: ...

class QueryDefinition(_message.Message):
    __slots__ = ("query",)
    QUERY_FIELD_NUMBER: _ClassVar[int]
    query: str
    def __init__(self, query: _Optional[str] = ...) -> None: ...

class QueryTable(_message.Message):
    __slots__ = ("rows",)
    ROWS_FIELD_NUMBER: _ClassVar[int]
    rows: _containers.RepeatedCompositeFieldContainer[_struct_pb2.Struct]
    def __init__(self, rows: _Optional[_Iterable[_Union[_struct_pb2.Struct, _Mapping]]] = ...) -> None: ...

class QueryResult(_message.Message):
    __slots__ = ("tables",)
    TABLES_FIELD_NUMBER: _ClassVar[int]
    tables: _containers.RepeatedCompositeFieldContainer[QueryTable]
    def __init__(self, tables: _Optional[_Iterable[_Union[QueryTable, _Mapping]]] = ...) -> None: ...

class ExecuteQueriesResponse(_message.Message):
    __slots__ = ("results",)
    RESULTS_FIELD_NUMBER: _ClassVar[int]
    results: _containers.RepeatedCompositeFieldContainer[QueryResult]
    def __init__(self, results: _Optional[_Iterable[_Union[QueryResult, _Mapping]]] = ...) -> None: ...
