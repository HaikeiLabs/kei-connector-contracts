# Connector contracts for Python

Install with `pip install haikeilabs-connector-contracts`.

```python
from google.protobuf import json_format
from kei.connectors.github.v1 import issues_pb2

response = json_format.Parse(
    '{"number": 1347}', issues_pb2.GetIssueResponse(),
    ignore_unknown_fields=True,
)
print(response.number)
```
