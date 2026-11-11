from google.protobuf import descriptor_pb2 as _descriptor_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class OperationClass(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    OPERATION_CLASS_UNSPECIFIED: _ClassVar[OperationClass]
    OPERATION_CLASS_READ: _ClassVar[OperationClass]
    OPERATION_CLASS_WRITE: _ClassVar[OperationClass]

class CredentialPlacement(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    CREDENTIAL_PLACEMENT_UNSPECIFIED: _ClassVar[CredentialPlacement]
    CREDENTIAL_PLACEMENT_HEADER_BEARER: _ClassVar[CredentialPlacement]
    CREDENTIAL_PLACEMENT_HEADER_CUSTOM: _ClassVar[CredentialPlacement]
    CREDENTIAL_PLACEMENT_QUERY: _ClassVar[CredentialPlacement]
    CREDENTIAL_PLACEMENT_BASIC: _ClassVar[CredentialPlacement]
    CREDENTIAL_PLACEMENT_NONE: _ClassVar[CredentialPlacement]
OPERATION_CLASS_UNSPECIFIED: OperationClass
OPERATION_CLASS_READ: OperationClass
OPERATION_CLASS_WRITE: OperationClass
CREDENTIAL_PLACEMENT_UNSPECIFIED: CredentialPlacement
CREDENTIAL_PLACEMENT_HEADER_BEARER: CredentialPlacement
CREDENTIAL_PLACEMENT_HEADER_CUSTOM: CredentialPlacement
CREDENTIAL_PLACEMENT_QUERY: CredentialPlacement
CREDENTIAL_PLACEMENT_BASIC: CredentialPlacement
CREDENTIAL_PLACEMENT_NONE: CredentialPlacement
OPERATION_FIELD_NUMBER: _ClassVar[int]
operation: _descriptor.FieldDescriptor

class Operation(_message.Message):
    __slots__ = ("id", "connector", "destructive", "http_method", "path_template", "scopes", "credential_placement")
    ID_FIELD_NUMBER: _ClassVar[int]
    CONNECTOR_FIELD_NUMBER: _ClassVar[int]
    CLASS_FIELD_NUMBER: _ClassVar[int]
    DESTRUCTIVE_FIELD_NUMBER: _ClassVar[int]
    HTTP_METHOD_FIELD_NUMBER: _ClassVar[int]
    PATH_TEMPLATE_FIELD_NUMBER: _ClassVar[int]
    SCOPES_FIELD_NUMBER: _ClassVar[int]
    CREDENTIAL_PLACEMENT_FIELD_NUMBER: _ClassVar[int]
    id: str
    connector: str
    destructive: bool
    http_method: str
    path_template: str
    scopes: _containers.RepeatedScalarFieldContainer[str]
    credential_placement: CredentialPlacement
    def __init__(self, id: _Optional[str] = ..., connector: _Optional[str] = ..., destructive: _Optional[bool] = ..., http_method: _Optional[str] = ..., path_template: _Optional[str] = ..., scopes: _Optional[_Iterable[str]] = ..., credential_placement: _Optional[_Union[CredentialPlacement, str]] = ..., **kwargs) -> None: ...
