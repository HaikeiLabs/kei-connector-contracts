from kei.contracts.v1 import options_pb2 as _options_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class EventTime(_message.Message):
    __slots__ = ("date_time", "date", "time_zone")
    DATE_TIME_FIELD_NUMBER: _ClassVar[int]
    DATE_FIELD_NUMBER: _ClassVar[int]
    TIME_ZONE_FIELD_NUMBER: _ClassVar[int]
    date_time: str
    date: str
    time_zone: str
    def __init__(self, date_time: _Optional[str] = ..., date: _Optional[str] = ..., time_zone: _Optional[str] = ...) -> None: ...

class Event(_message.Message):
    __slots__ = ("id", "summary", "description", "location", "start", "end", "status")
    ID_FIELD_NUMBER: _ClassVar[int]
    SUMMARY_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    LOCATION_FIELD_NUMBER: _ClassVar[int]
    START_FIELD_NUMBER: _ClassVar[int]
    END_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    id: str
    summary: str
    description: str
    location: str
    start: EventTime
    end: EventTime
    status: str
    def __init__(self, id: _Optional[str] = ..., summary: _Optional[str] = ..., description: _Optional[str] = ..., location: _Optional[str] = ..., start: _Optional[_Union[EventTime, _Mapping]] = ..., end: _Optional[_Union[EventTime, _Mapping]] = ..., status: _Optional[str] = ...) -> None: ...

class ListEventsRequest(_message.Message):
    __slots__ = ("calendar_id", "time_min", "time_max", "q", "max_results", "page_token")
    CALENDAR_ID_FIELD_NUMBER: _ClassVar[int]
    TIME_MIN_FIELD_NUMBER: _ClassVar[int]
    TIME_MAX_FIELD_NUMBER: _ClassVar[int]
    Q_FIELD_NUMBER: _ClassVar[int]
    MAX_RESULTS_FIELD_NUMBER: _ClassVar[int]
    PAGE_TOKEN_FIELD_NUMBER: _ClassVar[int]
    calendar_id: str
    time_min: str
    time_max: str
    q: str
    max_results: int
    page_token: str
    def __init__(self, calendar_id: _Optional[str] = ..., time_min: _Optional[str] = ..., time_max: _Optional[str] = ..., q: _Optional[str] = ..., max_results: _Optional[int] = ..., page_token: _Optional[str] = ...) -> None: ...

class ListEventsResponse(_message.Message):
    __slots__ = ("items", "next_page_token")
    ITEMS_FIELD_NUMBER: _ClassVar[int]
    NEXT_PAGE_TOKEN_FIELD_NUMBER: _ClassVar[int]
    items: _containers.RepeatedCompositeFieldContainer[Event]
    next_page_token: str
    def __init__(self, items: _Optional[_Iterable[_Union[Event, _Mapping]]] = ..., next_page_token: _Optional[str] = ...) -> None: ...

class GetEventRequest(_message.Message):
    __slots__ = ("calendar_id", "event_id")
    CALENDAR_ID_FIELD_NUMBER: _ClassVar[int]
    EVENT_ID_FIELD_NUMBER: _ClassVar[int]
    calendar_id: str
    event_id: str
    def __init__(self, calendar_id: _Optional[str] = ..., event_id: _Optional[str] = ...) -> None: ...

class GetEventResponse(_message.Message):
    __slots__ = ("id", "summary", "description", "location", "start", "end", "status")
    ID_FIELD_NUMBER: _ClassVar[int]
    SUMMARY_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    LOCATION_FIELD_NUMBER: _ClassVar[int]
    START_FIELD_NUMBER: _ClassVar[int]
    END_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    id: str
    summary: str
    description: str
    location: str
    start: EventTime
    end: EventTime
    status: str
    def __init__(self, id: _Optional[str] = ..., summary: _Optional[str] = ..., description: _Optional[str] = ..., location: _Optional[str] = ..., start: _Optional[_Union[EventTime, _Mapping]] = ..., end: _Optional[_Union[EventTime, _Mapping]] = ..., status: _Optional[str] = ...) -> None: ...

class CreateEventRequest(_message.Message):
    __slots__ = ("calendar_id", "summary", "description", "location", "start", "end")
    CALENDAR_ID_FIELD_NUMBER: _ClassVar[int]
    SUMMARY_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    LOCATION_FIELD_NUMBER: _ClassVar[int]
    START_FIELD_NUMBER: _ClassVar[int]
    END_FIELD_NUMBER: _ClassVar[int]
    calendar_id: str
    summary: str
    description: str
    location: str
    start: EventTime
    end: EventTime
    def __init__(self, calendar_id: _Optional[str] = ..., summary: _Optional[str] = ..., description: _Optional[str] = ..., location: _Optional[str] = ..., start: _Optional[_Union[EventTime, _Mapping]] = ..., end: _Optional[_Union[EventTime, _Mapping]] = ...) -> None: ...

class CreateEventResponse(_message.Message):
    __slots__ = ("id", "summary", "description", "location", "start", "end", "status")
    ID_FIELD_NUMBER: _ClassVar[int]
    SUMMARY_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    LOCATION_FIELD_NUMBER: _ClassVar[int]
    START_FIELD_NUMBER: _ClassVar[int]
    END_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    id: str
    summary: str
    description: str
    location: str
    start: EventTime
    end: EventTime
    status: str
    def __init__(self, id: _Optional[str] = ..., summary: _Optional[str] = ..., description: _Optional[str] = ..., location: _Optional[str] = ..., start: _Optional[_Union[EventTime, _Mapping]] = ..., end: _Optional[_Union[EventTime, _Mapping]] = ..., status: _Optional[str] = ...) -> None: ...

class UpdateEventRequest(_message.Message):
    __slots__ = ("calendar_id", "event_id", "summary", "description", "location", "start", "end")
    CALENDAR_ID_FIELD_NUMBER: _ClassVar[int]
    EVENT_ID_FIELD_NUMBER: _ClassVar[int]
    SUMMARY_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    LOCATION_FIELD_NUMBER: _ClassVar[int]
    START_FIELD_NUMBER: _ClassVar[int]
    END_FIELD_NUMBER: _ClassVar[int]
    calendar_id: str
    event_id: str
    summary: str
    description: str
    location: str
    start: EventTime
    end: EventTime
    def __init__(self, calendar_id: _Optional[str] = ..., event_id: _Optional[str] = ..., summary: _Optional[str] = ..., description: _Optional[str] = ..., location: _Optional[str] = ..., start: _Optional[_Union[EventTime, _Mapping]] = ..., end: _Optional[_Union[EventTime, _Mapping]] = ...) -> None: ...

class UpdateEventResponse(_message.Message):
    __slots__ = ("id", "summary", "description", "location", "start", "end", "status")
    ID_FIELD_NUMBER: _ClassVar[int]
    SUMMARY_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    LOCATION_FIELD_NUMBER: _ClassVar[int]
    START_FIELD_NUMBER: _ClassVar[int]
    END_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    id: str
    summary: str
    description: str
    location: str
    start: EventTime
    end: EventTime
    status: str
    def __init__(self, id: _Optional[str] = ..., summary: _Optional[str] = ..., description: _Optional[str] = ..., location: _Optional[str] = ..., start: _Optional[_Union[EventTime, _Mapping]] = ..., end: _Optional[_Union[EventTime, _Mapping]] = ..., status: _Optional[str] = ...) -> None: ...
