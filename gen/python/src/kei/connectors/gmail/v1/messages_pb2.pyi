from kei.contracts.v1 import options_pb2 as _options_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class ListMessagesRequest(_message.Message):
    __slots__ = ("user_id", "q", "page_token", "max_results")
    USER_ID_FIELD_NUMBER: _ClassVar[int]
    Q_FIELD_NUMBER: _ClassVar[int]
    PAGE_TOKEN_FIELD_NUMBER: _ClassVar[int]
    MAX_RESULTS_FIELD_NUMBER: _ClassVar[int]
    user_id: str
    q: str
    page_token: str
    max_results: int
    def __init__(self, user_id: _Optional[str] = ..., q: _Optional[str] = ..., page_token: _Optional[str] = ..., max_results: _Optional[int] = ...) -> None: ...

class MessageSummary(_message.Message):
    __slots__ = ("id", "thread_id")
    ID_FIELD_NUMBER: _ClassVar[int]
    THREAD_ID_FIELD_NUMBER: _ClassVar[int]
    id: str
    thread_id: str
    def __init__(self, id: _Optional[str] = ..., thread_id: _Optional[str] = ...) -> None: ...

class ListMessagesResponse(_message.Message):
    __slots__ = ("messages", "next_page_token")
    MESSAGES_FIELD_NUMBER: _ClassVar[int]
    NEXT_PAGE_TOKEN_FIELD_NUMBER: _ClassVar[int]
    messages: _containers.RepeatedCompositeFieldContainer[MessageSummary]
    next_page_token: str
    def __init__(self, messages: _Optional[_Iterable[_Union[MessageSummary, _Mapping]]] = ..., next_page_token: _Optional[str] = ...) -> None: ...

class GetMessageRequest(_message.Message):
    __slots__ = ("user_id", "message_id")
    USER_ID_FIELD_NUMBER: _ClassVar[int]
    MESSAGE_ID_FIELD_NUMBER: _ClassVar[int]
    user_id: str
    message_id: str
    def __init__(self, user_id: _Optional[str] = ..., message_id: _Optional[str] = ...) -> None: ...

class GetMessageResponse(_message.Message):
    __slots__ = ("id", "thread_id", "snippet", "internal_date")
    ID_FIELD_NUMBER: _ClassVar[int]
    THREAD_ID_FIELD_NUMBER: _ClassVar[int]
    SNIPPET_FIELD_NUMBER: _ClassVar[int]
    INTERNAL_DATE_FIELD_NUMBER: _ClassVar[int]
    id: str
    thread_id: str
    snippet: str
    internal_date: int
    def __init__(self, id: _Optional[str] = ..., thread_id: _Optional[str] = ..., snippet: _Optional[str] = ..., internal_date: _Optional[int] = ...) -> None: ...
