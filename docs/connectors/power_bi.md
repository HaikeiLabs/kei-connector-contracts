# power_bi connector

## Connections

Source: [`connector.proto`](../../proto/kei/connectors/power_bi/v1/connector.proto).

| Name | Kind | Auth method | Account model | Setup fields |
| --- | --- | --- | --- | --- |
| `user_obo` | `CONNECTION_KIND_TOKEN_EXCHANGE` | `kei_token_exchange` | `ACCOUNT_MODEL_PER_USER` | — |

## Operations

| Id | Method | Path or GraphQL | Class | Destructive | Connections | Resource template |
| --- | --- | --- | --- | --- | --- | --- |
| `power_bi.datasets.execute_queries` | `POST` | `/v1.0/myorg/groups/{group_id}/datasets/{dataset_id}/executeQueries` | `read` | false | `user_obo` | `groups/{group_id}/datasets/{dataset_id}` |
| `power_bi.datasets.list` | `GET` | `/v1.0/myorg/groups/{group_id}/datasets` | `read` | false | `user_obo` | `groups/{group_id}` |
