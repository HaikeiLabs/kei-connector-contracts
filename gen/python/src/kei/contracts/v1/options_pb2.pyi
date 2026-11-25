from google.protobuf import descriptor_pb2 as _descriptor_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
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

class ConnectionKind(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    CONNECTION_KIND_UNSPECIFIED: _ClassVar[ConnectionKind]
    CONNECTION_KIND_HTTP_BEARER: _ClassVar[ConnectionKind]
    CONNECTION_KIND_HTTP_HEADER: _ClassVar[ConnectionKind]
    CONNECTION_KIND_HTTP_BASIC: _ClassVar[ConnectionKind]
    CONNECTION_KIND_SQL_DSN: _ClassVar[ConnectionKind]
    CONNECTION_KIND_KEY_PAIR_JWT: _ClassVar[ConnectionKind]
    CONNECTION_KIND_TOKEN_EXCHANGE: _ClassVar[ConnectionKind]

class AccountModel(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    ACCOUNT_MODEL_UNSPECIFIED: _ClassVar[AccountModel]
    ACCOUNT_MODEL_PER_USER: _ClassVar[AccountModel]
    ACCOUNT_MODEL_SHARED: _ClassVar[AccountModel]

class ErrorKind(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    ERROR_KIND_UNSPECIFIED: _ClassVar[ErrorKind]
    ERROR_KIND_NOT_FOUND: _ClassVar[ErrorKind]
    ERROR_KIND_UNAUTHORIZED: _ClassVar[ErrorKind]
    ERROR_KIND_FORBIDDEN: _ClassVar[ErrorKind]
    ERROR_KIND_RATE_LIMITED: _ClassVar[ErrorKind]
    ERROR_KIND_CONFLICT: _ClassVar[ErrorKind]
    ERROR_KIND_INVALID: _ClassVar[ErrorKind]
    ERROR_KIND_UNAVAILABLE: _ClassVar[ErrorKind]
    ERROR_KIND_CONTRACT_VIOLATION: _ClassVar[ErrorKind]

class ParamIn(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    PARAM_IN_UNSPECIFIED: _ClassVar[ParamIn]
    PARAM_IN_PATH: _ClassVar[ParamIn]
    PARAM_IN_QUERY: _ClassVar[ParamIn]
    PARAM_IN_BODY: _ClassVar[ParamIn]
    PARAM_IN_HEADER: _ClassVar[ParamIn]

class SemanticAction(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    SEMANTIC_ACTION_UNSPECIFIED: _ClassVar[SemanticAction]
    SEMANTIC_ACTION_READ: _ClassVar[SemanticAction]
    SEMANTIC_ACTION_CREATE: _ClassVar[SemanticAction]
    SEMANTIC_ACTION_UPDATE: _ClassVar[SemanticAction]
    SEMANTIC_ACTION_DELETE: _ClassVar[SemanticAction]
OPERATION_CLASS_UNSPECIFIED: OperationClass
OPERATION_CLASS_READ: OperationClass
OPERATION_CLASS_WRITE: OperationClass
CREDENTIAL_PLACEMENT_UNSPECIFIED: CredentialPlacement
CREDENTIAL_PLACEMENT_HEADER_BEARER: CredentialPlacement
CREDENTIAL_PLACEMENT_HEADER_CUSTOM: CredentialPlacement
CREDENTIAL_PLACEMENT_QUERY: CredentialPlacement
CREDENTIAL_PLACEMENT_BASIC: CredentialPlacement
CREDENTIAL_PLACEMENT_NONE: CredentialPlacement
CONNECTION_KIND_UNSPECIFIED: ConnectionKind
CONNECTION_KIND_HTTP_BEARER: ConnectionKind
CONNECTION_KIND_HTTP_HEADER: ConnectionKind
CONNECTION_KIND_HTTP_BASIC: ConnectionKind
CONNECTION_KIND_SQL_DSN: ConnectionKind
CONNECTION_KIND_KEY_PAIR_JWT: ConnectionKind
CONNECTION_KIND_TOKEN_EXCHANGE: ConnectionKind
ACCOUNT_MODEL_UNSPECIFIED: AccountModel
ACCOUNT_MODEL_PER_USER: AccountModel
ACCOUNT_MODEL_SHARED: AccountModel
ERROR_KIND_UNSPECIFIED: ErrorKind
ERROR_KIND_NOT_FOUND: ErrorKind
ERROR_KIND_UNAUTHORIZED: ErrorKind
ERROR_KIND_FORBIDDEN: ErrorKind
ERROR_KIND_RATE_LIMITED: ErrorKind
ERROR_KIND_CONFLICT: ErrorKind
ERROR_KIND_INVALID: ErrorKind
ERROR_KIND_UNAVAILABLE: ErrorKind
ERROR_KIND_CONTRACT_VIOLATION: ErrorKind
PARAM_IN_UNSPECIFIED: ParamIn
PARAM_IN_PATH: ParamIn
PARAM_IN_QUERY: ParamIn
PARAM_IN_BODY: ParamIn
PARAM_IN_HEADER: ParamIn
SEMANTIC_ACTION_UNSPECIFIED: SemanticAction
SEMANTIC_ACTION_READ: SemanticAction
SEMANTIC_ACTION_CREATE: SemanticAction
SEMANTIC_ACTION_UPDATE: SemanticAction
SEMANTIC_ACTION_DELETE: SemanticAction
OPERATION_FIELD_NUMBER: _ClassVar[int]
operation: _descriptor.FieldDescriptor
CONNECTOR_FIELD_NUMBER: _ClassVar[int]
connector: _descriptor.FieldDescriptor
PARAM_FIELD_NUMBER: _ClassVar[int]
param: _descriptor.FieldDescriptor

class SetupField(_message.Message):
    __slots__ = ("name", "type", "pattern", "secret")
    NAME_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    PATTERN_FIELD_NUMBER: _ClassVar[int]
    SECRET_FIELD_NUMBER: _ClassVar[int]
    name: str
    type: str
    pattern: str
    secret: bool
    def __init__(self, name: _Optional[str] = ..., type: _Optional[str] = ..., pattern: _Optional[str] = ..., secret: _Optional[bool] = ...) -> None: ...

class Connection(_message.Message):
    __slots__ = ("name", "kind", "auth_method", "account_model", "base_url", "base_url_field", "setup_fields", "scopes", "engine", "header_name")
    NAME_FIELD_NUMBER: _ClassVar[int]
    KIND_FIELD_NUMBER: _ClassVar[int]
    AUTH_METHOD_FIELD_NUMBER: _ClassVar[int]
    ACCOUNT_MODEL_FIELD_NUMBER: _ClassVar[int]
    BASE_URL_FIELD_NUMBER: _ClassVar[int]
    BASE_URL_FIELD_FIELD_NUMBER: _ClassVar[int]
    SETUP_FIELDS_FIELD_NUMBER: _ClassVar[int]
    SCOPES_FIELD_NUMBER: _ClassVar[int]
    ENGINE_FIELD_NUMBER: _ClassVar[int]
    HEADER_NAME_FIELD_NUMBER: _ClassVar[int]
    name: str
    kind: ConnectionKind
    auth_method: str
    account_model: AccountModel
    base_url: str
    base_url_field: str
    setup_fields: _containers.RepeatedCompositeFieldContainer[SetupField]
    scopes: _containers.RepeatedScalarFieldContainer[str]
    engine: str
    header_name: str
    def __init__(self, name: _Optional[str] = ..., kind: _Optional[_Union[ConnectionKind, str]] = ..., auth_method: _Optional[str] = ..., account_model: _Optional[_Union[AccountModel, str]] = ..., base_url: _Optional[str] = ..., base_url_field: _Optional[str] = ..., setup_fields: _Optional[_Iterable[_Union[SetupField, _Mapping]]] = ..., scopes: _Optional[_Iterable[str]] = ..., engine: _Optional[str] = ..., header_name: _Optional[str] = ...) -> None: ...

class Connector(_message.Message):
    __slots__ = ("id", "base_url", "connections")
    ID_FIELD_NUMBER: _ClassVar[int]
    BASE_URL_FIELD_NUMBER: _ClassVar[int]
    CONNECTIONS_FIELD_NUMBER: _ClassVar[int]
    id: str
    base_url: str
    connections: _containers.RepeatedCompositeFieldContainer[Connection]
    def __init__(self, id: _Optional[str] = ..., base_url: _Optional[str] = ..., connections: _Optional[_Iterable[_Union[Connection, _Mapping]]] = ...) -> None: ...

class ErrorMapping(_message.Message):
    __slots__ = ("status", "kind", "retryable")
    STATUS_FIELD_NUMBER: _ClassVar[int]
    KIND_FIELD_NUMBER: _ClassVar[int]
    RETRYABLE_FIELD_NUMBER: _ClassVar[int]
    status: int
    kind: ErrorKind
    retryable: bool
    def __init__(self, status: _Optional[int] = ..., kind: _Optional[_Union[ErrorKind, str]] = ..., retryable: _Optional[bool] = ...) -> None: ...

class Param(_message.Message):
    __slots__ = ("name", "multi_segment")
    IN_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    MULTI_SEGMENT_FIELD_NUMBER: _ClassVar[int]
    name: str
    multi_segment: bool
    def __init__(self, name: _Optional[str] = ..., multi_segment: _Optional[bool] = ..., **kwargs) -> None: ...

class Semantics(_message.Message):
    __slots__ = ("action", "resource_type", "version")
    ACTION_FIELD_NUMBER: _ClassVar[int]
    RESOURCE_TYPE_FIELD_NUMBER: _ClassVar[int]
    VERSION_FIELD_NUMBER: _ClassVar[int]
    action: SemanticAction
    resource_type: str
    version: int
    def __init__(self, action: _Optional[_Union[SemanticAction, str]] = ..., resource_type: _Optional[str] = ..., version: _Optional[int] = ...) -> None: ...

class Operation(_message.Message):
    __slots__ = ("id", "connector", "destructive", "http_method", "path_template", "scopes", "credential_placement", "method_name", "connections", "resource_template", "errors", "deprecated", "replaced_by", "graphql_document", "positional_helper", "semantics", "cursor_field", "next_cursor_field", "result_required")
    ID_FIELD_NUMBER: _ClassVar[int]
    CONNECTOR_FIELD_NUMBER: _ClassVar[int]
    CLASS_FIELD_NUMBER: _ClassVar[int]
    DESTRUCTIVE_FIELD_NUMBER: _ClassVar[int]
    HTTP_METHOD_FIELD_NUMBER: _ClassVar[int]
    PATH_TEMPLATE_FIELD_NUMBER: _ClassVar[int]
    SCOPES_FIELD_NUMBER: _ClassVar[int]
    CREDENTIAL_PLACEMENT_FIELD_NUMBER: _ClassVar[int]
    METHOD_NAME_FIELD_NUMBER: _ClassVar[int]
    CONNECTIONS_FIELD_NUMBER: _ClassVar[int]
    RESOURCE_TEMPLATE_FIELD_NUMBER: _ClassVar[int]
    ERRORS_FIELD_NUMBER: _ClassVar[int]
    DEPRECATED_FIELD_NUMBER: _ClassVar[int]
    REPLACED_BY_FIELD_NUMBER: _ClassVar[int]
    GRAPHQL_DOCUMENT_FIELD_NUMBER: _ClassVar[int]
    POSITIONAL_HELPER_FIELD_NUMBER: _ClassVar[int]
    SEMANTICS_FIELD_NUMBER: _ClassVar[int]
    CURSOR_FIELD_FIELD_NUMBER: _ClassVar[int]
    NEXT_CURSOR_FIELD_FIELD_NUMBER: _ClassVar[int]
    RESULT_REQUIRED_FIELD_NUMBER: _ClassVar[int]
    id: str
    connector: str
    destructive: bool
    http_method: str
    path_template: str
    scopes: _containers.RepeatedScalarFieldContainer[str]
    credential_placement: CredentialPlacement
    method_name: str
    connections: _containers.RepeatedScalarFieldContainer[str]
    resource_template: str
    errors: _containers.RepeatedCompositeFieldContainer[ErrorMapping]
    deprecated: bool
    replaced_by: str
    graphql_document: str
    positional_helper: bool
    semantics: Semantics
    cursor_field: str
    next_cursor_field: str
    result_required: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, id: _Optional[str] = ..., connector: _Optional[str] = ..., destructive: _Optional[bool] = ..., http_method: _Optional[str] = ..., path_template: _Optional[str] = ..., scopes: _Optional[_Iterable[str]] = ..., credential_placement: _Optional[_Union[CredentialPlacement, str]] = ..., method_name: _Optional[str] = ..., connections: _Optional[_Iterable[str]] = ..., resource_template: _Optional[str] = ..., errors: _Optional[_Iterable[_Union[ErrorMapping, _Mapping]]] = ..., deprecated: _Optional[bool] = ..., replaced_by: _Optional[str] = ..., graphql_document: _Optional[str] = ..., positional_helper: _Optional[bool] = ..., semantics: _Optional[_Union[Semantics, _Mapping]] = ..., cursor_field: _Optional[str] = ..., next_cursor_field: _Optional[str] = ..., result_required: _Optional[_Iterable[str]] = ..., **kwargs) -> None: ...
