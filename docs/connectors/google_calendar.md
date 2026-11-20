# google_calendar connector

## Connections

Source: [`connector.proto`](../../proto/kei/connectors/google_calendar/v1/connector.proto).

| Name | Kind | Auth method | Account model | Setup fields |
| --- | --- | --- | --- | --- |
| `user_oauth` | `CONNECTION_KIND_HTTP_BEARER` | `oauth` | `ACCOUNT_MODEL_PER_USER` | — |

## Operations

| Id | Method | Path or GraphQL | Class | Destructive | Connections | Resource template |
| --- | --- | --- | --- | --- | --- | --- |
| `google_calendar.calendars.list` | `GET` | `/calendar/v3/users/me/calendarList` | `read` | false | `user_oauth` | `users/me/calendarList` |
| `google_calendar.events.list` | `GET` | `/calendar/v3/calendars/{calendar_id}/events` | `read` | false | `user_oauth` | `calendars/{calendar_id}` |
| `google_calendar.events.get` | `GET` | `/calendar/v3/calendars/{calendar_id}/events/{event_id}` | `read` | false | `user_oauth` | `calendars/{calendar_id}` |
| `google_calendar.events.create` | `POST` | `/calendar/v3/calendars/{calendar_id}/events` | `write` | false | `user_oauth` | `calendars/{calendar_id}` |
| `google_calendar.events.update` | `PATCH` | `/calendar/v3/calendars/{calendar_id}/events/{event_id}` | `write` | false | `user_oauth` | `calendars/{calendar_id}` |
