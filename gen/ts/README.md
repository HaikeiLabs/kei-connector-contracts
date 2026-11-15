# Connector contracts for TypeScript

Install with `npm install @haikeilabs/connector-contracts`.

```ts
import { fromJson } from "@bufbuild/protobuf";
import { GetIssueResponseSchema } from "@haikeilabs/connector-contracts/kei/connectors/github/v1/issues_pb";

const response = fromJson(GetIssueResponseSchema, { number: 1347 }, {
  ignoreUnknownFields: true,
});
console.log(response.number);
```

Provider JSON may contain int64 values as JSON numbers larger than JavaScript's
safe integer range. `JSON.parse` rounds those values before `fromJson` sees them.
Parse the original response text losslessly, converting unsafe integers to
decimal strings, which protobuf-es accepts for int64 fields. For example, after
installing `lossless-json`:

```ts
import { isInteger, parse } from "lossless-json";
import type { JsonValue } from "@bufbuild/protobuf";

const json = parse(responseText, undefined, {
  parseNumber: (value) => {
    if (isInteger(value)) {
      const integer = BigInt(value);
      if (integer > BigInt(Number.MAX_SAFE_INTEGER) ||
          integer < BigInt(Number.MIN_SAFE_INTEGER)) return value;
    }
    return Number(value);
  },
}) as JsonValue;
const issue = fromJson(GetIssueResponseSchema, json, { ignoreUnknownFields: true });
```
