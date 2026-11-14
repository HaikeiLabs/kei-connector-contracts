from google.protobuf import struct_pb2 as _struct_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ToolCall(_message.Message):
    __slots__ = ("server", "tool", "arguments")
    SERVER_FIELD_NUMBER: _ClassVar[int]
    TOOL_FIELD_NUMBER: _ClassVar[int]
    ARGUMENTS_FIELD_NUMBER: _ClassVar[int]
    server: str
    tool: str
    arguments: _struct_pb2.Struct
    def __init__(self, server: _Optional[str] = ..., tool: _Optional[str] = ..., arguments: _Optional[_Union[_struct_pb2.Struct, _Mapping]] = ...) -> None: ...

class ToolResult(_message.Message):
    __slots__ = ("content", "is_error")
    CONTENT_FIELD_NUMBER: _ClassVar[int]
    IS_ERROR_FIELD_NUMBER: _ClassVar[int]
    content: _containers.RepeatedCompositeFieldContainer[_struct_pb2.Struct]
    is_error: bool
    def __init__(self, content: _Optional[_Iterable[_Union[_struct_pb2.Struct, _Mapping]]] = ..., is_error: _Optional[bool] = ...) -> None: ...
