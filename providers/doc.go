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

// Package providers is a stable, harness-consumable library of governed
// connector clients (ADR-036 section 9.2). A harness — the agentware SDK, a
// chat harness, or a custom SDK agent — imports this package to read and
// write provider data, supplying its own credential and HTTP client.
//
// Each provider constructor takes a credential (a CredentialResolver) and an
// *http.Client and returns a client bound to the provider's fixed base
// endpoint. The client resolves no credential on its own, and this package
// never imports kei-connector-runtime: credential resolution and policy are
// the harness's, and the import boundary is enforced by
// TestProvidersImportBoundary. Breaking changes follow the module's semantic
// versioning.
//
// Every client is a read-only seam over the contract: the backend interface
// is the integration point, and a client exposes only the capabilities the
// contract declares for its provider. New operations must be added to the
// contract before a client may expose them.
package providers
