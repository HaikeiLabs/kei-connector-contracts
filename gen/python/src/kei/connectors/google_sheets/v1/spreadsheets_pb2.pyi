from google.protobuf import struct_pb2 as _struct_pb2
from kei.contracts.v1 import options_pb2 as _options_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class GetValuesRequest(_message.Message):
    __slots__ = ("spreadsheet_id", "range")
    SPREADSHEET_ID_FIELD_NUMBER: _ClassVar[int]
    RANGE_FIELD_NUMBER: _ClassVar[int]
    spreadsheet_id: str
    range: str
    def __init__(self, spreadsheet_id: _Optional[str] = ..., range: _Optional[str] = ...) -> None: ...

class GetValuesResponse(_message.Message):
    __slots__ = ("range", "values")
    RANGE_FIELD_NUMBER: _ClassVar[int]
    VALUES_FIELD_NUMBER: _ClassVar[int]
    range: str
    values: _containers.RepeatedCompositeFieldContainer[_struct_pb2.ListValue]
    def __init__(self, range: _Optional[str] = ..., values: _Optional[_Iterable[_Union[_struct_pb2.ListValue, _Mapping]]] = ...) -> None: ...

class GetSpreadsheetRequest(_message.Message):
    __slots__ = ("spreadsheet_id",)
    SPREADSHEET_ID_FIELD_NUMBER: _ClassVar[int]
    spreadsheet_id: str
    def __init__(self, spreadsheet_id: _Optional[str] = ...) -> None: ...

class SpreadsheetProperties(_message.Message):
    __slots__ = ("title",)
    TITLE_FIELD_NUMBER: _ClassVar[int]
    title: str
    def __init__(self, title: _Optional[str] = ...) -> None: ...

class GetSpreadsheetResponse(_message.Message):
    __slots__ = ("spreadsheet_id", "properties")
    SPREADSHEET_ID_FIELD_NUMBER: _ClassVar[int]
    PROPERTIES_FIELD_NUMBER: _ClassVar[int]
    spreadsheet_id: str
    properties: SpreadsheetProperties
    def __init__(self, spreadsheet_id: _Optional[str] = ..., properties: _Optional[_Union[SpreadsheetProperties, _Mapping]] = ...) -> None: ...
