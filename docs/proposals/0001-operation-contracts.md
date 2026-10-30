# Proposal 0001: Operation contracts

**Status:** Proposal. Not accepted. For discussion with the maintainers.
Tracking: HAI-915.

Nothing in this document is decided, and nothing in it changes the repository
today. It describes a possible direction, what it would take, and what would
make us stop. This is the first file in `docs/proposals/`; there was no
existing proposal or ADR folder in this repository, so this path is a
suggestion too.

## 1. Motivation

Today this module holds:

- the connector contract in Go: `contract`, `providers`, `capability`, `setup`,
  `access`, `envelope`, `governance`, `audit`, `invoke`;
- the shared policy model (`policy`) and the policy bundle parsers (`bundle`);
- JSON schemas, including `kei.tool-definition/v1` and the tool call template
  profile built on it (`docs/tool-call-templates.md`), with one entry file per
  connector.

Two things motivate looking at a different shape.

**Tool definitions carry more than an API contract.** An entry holds a tool
name, a model-facing description and a parameter list next to the transport
(method, URL template, argument placement). The model-facing part is a property
of the harness and the model, not of the provider's API. Putting it in the
contract limits how an agent SDK can build model-specific tool calls and
reasoning flavors, and it shows up as clutter wherever the catalog is listed
for people.

**There is no single source of truth for an API call across languages.** The Go
`providers` package has hand-written typed clients. The template path
describes calls as data, but only the Go agent SDK has an engine that reads it,
and its result and error handling are not part of the entry. Each additional
language would need its own copy of the rules.

## 2. Proposal

Contracts would hold **API contracts and data contracts for connectors**, and
the unit would be an **operation**, not a tool.

An operation declares:

| Field | Meaning |
| --- | --- |
| `id` | Stable dotted name, for example `github.issues.get` |
| `connector` | The connector it belongs to |
| `request` | Method, path or URL template, parameters with their location (path, query, body, header), body encoding and limits |
| `response` | Typed result |
| `errors` | Status to meaning, and whether it is retryable |
| `scopes` | The provider scope or credential scope required, and where the credential is placed (never its value) |
| `class` | `read` or `write`, plus a destructive flag |

It carries no model-facing text and no model-specific schema. The agent SDK
builds the tool call from operations using its own paradigms for each model
family.

### Why policy fits

The policy package already matches connector targets by dot segment, with `*`
as a segment wildcard: `connector:github.*` matches `github.issues.get`. An
operation id is therefore a policy target as the model stands today:

```text
connector:github.issues.get        one operation
connector:github.issues.*          every issues operation
connector:github.*                 the connector
```

For calls that have no contract operation (a tool discovered from a remote MCP
server, a built-in tool of a harness), the agent SDK could assign ids in the
same space at registration, for example `mcp.<server>.<tool>` and
`builtin.<name>`. Then every call, whether it comes from a bring-your-own
harness author, a remote MCP server or a built-in, is an operation id and a
`read`/`write` class to policy and audit. Whether the id space needs reserved
prefixes or a policy schema change is an open question (section 8).

The same layer is a place to attach later features that need to talk about
operations, such as managing contracts semantically or describing workflows.
This proposal does not design them.

## 3. Code generation

One source of truth, generated per language, requires an IDL. The options:

| | Protobuf with buf | OpenAPI-first | Bespoke generator over JSON Schema |
| --- | --- | --- | --- |
| One source, many languages | Mature for Go, TypeScript, Python, Rust | Good for REST clients; quality varies by language, weakest for Rust | Whatever we write; we own every language's templates |
| Custom semantics (id, scopes, class) | Custom options on services and methods | `x-` extensions that generators ignore; we would parse them ourselves | Total control |
| Mechanical breaking-change check | Strong (`buf breaking`) | Available, less strict | We would have to build it |
| Third-party REST and GraphQL | Awkward: we call APIs we do not own, so their shapes are re-modelled as messages; GraphQL needs the query text carried as an option | Natural for REST, no GraphQL | Natural; we model what we call |
| Dependency impact | Generated Go needs a protobuf runtime library | Generated clients usually need a runtime library | Can be dependency-free |

**Recommendation, with moderate confidence:** protobuf with buf for the
messages and operation definitions, kei-specific custom options, and one small
generator plugin that emits the request builder, typed client and
response/error decoding for each language. The deciding reasons are the
mechanical additive-only check and multi-language support. The fallback, if the
experiment in section 5 shows protobuf fights the third-party REST shapes, is
the bespoke generator over JSON Schema with our own breaking-change check.
OpenAPI would be at most a one-way importer for drafting operations from a
provider's published spec; it would not be the source.

**Why grpc-gateway alone is the wrong direction:** it exposes a gRPC service as
a REST API. This module needs the inverse, a description of REST and GraphQL
calls to APIs we do not run. We would borrow its idea (HTTP annotations on
methods), not use the tool.

### Shipping without breaking the module rules

The root module is stdlib-only and CI fails if `go.mod` gains a `require`
block. Generated protobuf code needs a runtime library, so it cannot live in the
root module. Proposed layout:

- IDL files, conformance fixtures and the generator stay in this repository.
  The root Go module stays dependency-free.
- Generated clients ship as separate packages: a separate Go module path, an npm
  package, a PyPI package, a crate later. They are built in CI and published
  from the same tag. Tags continue to be cut by the maintainers.
- Languages: Go, then TypeScript, then Python, released at the same version
  (lockstep); Rust later. Because the Python and TypeScript engines do not exist
  yet, "lockstep" in the early phases means one version number and shared
  fixtures; a hard "every release ships every language" gate would start once
  Python exists (open question 3).
- Versioning: additive changes within v1; v2 only for breaking changes. The
  breaking-change check is the gate.

## 4. What is not changing

This proposal removes nothing and deprecates nothing in this repository.

- The tool call templates (`docs/tool-call-templates.md`, the
  `kei.tool-definition/v1` schemas, `calltemplate`, `schemas.FS`) stay as they
  are and keep working for their consumers.
- `providers/` stays; whether it is ever retired is a separate question.
- The policy package, bundle parsers and the other schemas are unchanged.
- The runtime stays decision-only: it authorizes, passes credentials through
  and audits. Writes remain actions the harness performs; this proposal only
  describes the calls.

If the direction is accepted, a compatibility window and a removal path for the
tool-definition schemas would be a **later, separate proposal**. A generator
could emit the older tool-definition shape from operations during such a
window; that is a possibility, not a plan. The related plan to deprecate the
`kei-agents` example repository is a separate item and is not covered here.

## 5. Phases

Each phase is useful on its own and can be the last one.

1. **Operation spec, fixtures, IDL experiment.** Define the operation fields,
   write golden request/response fixtures for about ten operations taken from
   the existing GitHub, Linear and `http_api` entries (without model-facing
   text), and run the bake-off below. Useful alone: a tool-free list of the
   things policy can be attached to, and a clear view of which entries are API
   operations and which are harness actions.
2. **Go generator and typed client** on the agent SDK's existing request
   builder, with response decoding and error mapping. One real function in a
   harness moves from the template call to the typed call.
3. **TypeScript** from the same IDL and fixtures.
4. **Python**, after which every release ships all languages.
5. **Uniform operations for MCP and built-in tools** in the agent SDK, so policy
   and audit see one id space.
6. **Rust**, later.

### Smallest experiment

About a day, no change to this repository's code: hand-write typed
`github.get_issue` and `github.file_write` and compare them with the
template-driven call. The known mismatch (`{repo}` in the contracts entry,
`{repository}` in a harness's hand-written transport) is the test case: does
the typed version make that drift impossible or at least a compile error?
Then express three operations (`github.get_issue`; `github.file_write`, a `PUT`
with a base64 body; a Linear GraphQL create) in protobuf with buf, a bespoke
JSON Schema spec and OpenAPI. Score each on: can it carry id, scopes and class;
does it generate Go and TypeScript from one file; is a deliberate breaking edit
caught mechanically; how much template code would we own per language; does the
root module stay dependency-free. Pass bar: one candidate expresses all three
with no escape hatch and catches the breaking edit. If protobuf needs an escape
hatch for two of the three, choose the bespoke spec.

Stop condition: if the typed versions add no safety or clarity over the
template call, do not continue.

## 6. Risks

- Four language engines in lockstep is the dominant cost. Only the Go engine
  exists today.
- Most entries have no declared result or error shape; authoring them is real
  work and is where provider drift will show.
- Protobuf modelling of third-party REST may be a poor fit.
- A lasting toolchain (buf, plugins, CI) in a repository that has been
  dependency-free by design.
- Migration load on every consumer of the current tool-definition path, even
  with a compatibility window.
- Maintainer load: one tag drives all language packages.

## 7. The strongest argument against

The near-term benefit, uniform policy attachment, may need only an operation id
convention that fits the existing dot-segment matching, which costs almost
nothing and needs no IDL. A single IDL with code generation is a platform bet
whose payoff is mostly future features. Building the generator first risks
paying for four language ports before anyone has shown that typed operations
earn their keep. Also, most consumers use tools dynamically (a chat UI, remote
MCP servers) and never call typed functions by hand. Mitigation: Phase 1 ships
the id convention and the operation list without committing to an IDL.

## 8. Open questions for the maintainers

1. **Who is the typed layer for?** Recommendation: harness authors and agent SDK
   users, plus the uniform id space for MCP and built-ins, not every consumer.
2. **Does the tool-era permission vocabulary** (`github.read`, `linear_write`,
   and the fixed enum in the tool-definition schema) give way to operation id,
   class and provider scope in the new layer? Recommendation: yes for operations,
   with the tool-definition schema untouched.
3. **What does "lockstep" mean before Python exists?** Recommendation: same
   version number and shared fixtures from the start, a hard all-languages gate
   once Python lands.
4. **Where does model-facing description live?** Recommendation: in the agent
   SDK, as data it owns, not in contracts.
5. **How are MCP and built-in ids assigned?** Recommendation: the agent SDK
   synthesizes them under reserved prefixes and policy matches them with the
   existing dot-segment rule, no policy schema change.
6. **Which entries are harness actions rather than API operations?** Entries
   with no transport, such as a "create fix pull request" action where the
   server picks the repository and branch, describe a harness action. Recommendation:
   keep them out of operations, and decide in a later proposal where they live.
