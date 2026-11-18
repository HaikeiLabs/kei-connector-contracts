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

// Package connection is the hand-written runtime that generated connector
// operations run on (proposal 0002, section 5.2).
//
// It supplies the authenticated connection types (HTTP and SQL), the
// operation descriptor, the path renderer and the error mapping. It takes
// credentials as arguments and never fetches them: it has no kei-proxy
// client, no secret-store client, no environment reads and no logging, and it
// reaches the network only through the *http.Client or *sql.DB the caller
// passes in. No error, String or GoString value it produces contains a
// credential or a provider response body.
package connection

// Class is an operation's class. The zero value is neither Read nor Write,
// so an unset descriptor is never mistaken for a read.
type Class int

const (
	// Read operations do not change provider state.
	Read Class = iota + 1
	// Write operations change provider state.
	Write
)

// String returns "read", "write" or "unspecified".
func (c Class) String() string {
	switch c {
	case Read:
		return "read"
	case Write:
		return "write"
	default:
		return "unspecified"
	}
}

// Op describes one generated operation, for authorization and audit. ID is
// the operation id, which is also the capability and the policy target
// (connector:<ID>).
type Op struct {
	ID          string
	Class       Class
	Destructive bool
}
