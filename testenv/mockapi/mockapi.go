// Copyright 2026 Haikei Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package mockapi serves a provider-shaped HTTP API generated from the
// repository's contract metadata. It exists so integration tests can point
// generated connection runtimes at a faithful stand-in of the real
// providers: the routes, credential placements and response shapes all come
// from the same contract declarations the generators read.
//
// What it serves:
//
//   - One route per manifest operation. HTTP operations match their method
//     and path template; GraphQL operations match POST /graphql by the
//     document's root field. SQL operations (database.query.run) have no
//     HTTP route: they run over a DSN, not HTTP.
//   - The operation's response fixture (fixtures/operations/<connector>/…),
//     verbatim for HTTP and wrapped in {"data": …} for GraphQL.
//   - Seeded synthetic pages for operations with cursor fields: the first
//     page serves the fixture with the next cursor set; a follow-up request
//     carrying that cursor serves a varied list of a different length and
//     ends the pagination. Page tokens are deterministic, so tests can
//     assert on them.
//   - Failure injection via the X-Mock-Fail request header, and request
//     validation (credential placement, path escaping) as 400s with a
//     reason.
package mockapi
