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

package connection

import (
	"fmt"
	"strings"

	contracts "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/contracts/v1"
)

// Error is a provider error returned by a generated operation. It carries
// only the operation id, the status and the mapped kind: never the response
// body, a header or a credential.
type Error struct {
	Kind      contracts.ErrorKind
	Status    int
	Retryable bool
	OpID      string
}

// Error returns, for example,
// "connection: github.issues.get: not_found (status 404)".
func (e *Error) Error() string {
	retry := ""
	if e.Retryable {
		retry = ", retryable"
	}
	return fmt.Sprintf("connection: %s: %s (status %d%s)", e.OpID, kindName(e.Kind), e.Status, retry)
}

// MapError maps a provider HTTP status to an error. A 2xx status returns
// nil. Otherwise the operation's declared table decides; a status the table
// does not list falls back to the defaults: 400 and 422 invalid, 401
// unauthorized, 403 forbidden, 404 not found, 409 conflict, 429 rate limited
// (retryable), 5xx unavailable (502, 503 and 504 retryable), anything else
// unspecified.
func MapError(op Op, status int, table []*contracts.ErrorMapping) error {
	if status >= 200 && status < 300 {
		return nil
	}
	e := &Error{Status: status, OpID: op.ID}
	for _, m := range table {
		if int(m.GetStatus()) == status {
			e.Kind, e.Retryable = m.GetKind(), m.GetRetryable()
			return e
		}
	}
	e.Kind, e.Retryable = defaultKind(status)
	return e
}

func defaultKind(status int) (contracts.ErrorKind, bool) {
	switch {
	case status == 400 || status == 422:
		return contracts.ErrorKind_ERROR_KIND_INVALID, false
	case status == 401:
		return contracts.ErrorKind_ERROR_KIND_UNAUTHORIZED, false
	case status == 403:
		return contracts.ErrorKind_ERROR_KIND_FORBIDDEN, false
	case status == 404:
		return contracts.ErrorKind_ERROR_KIND_NOT_FOUND, false
	case status == 409:
		return contracts.ErrorKind_ERROR_KIND_CONFLICT, false
	case status == 429:
		return contracts.ErrorKind_ERROR_KIND_RATE_LIMITED, true
	case status == 502 || status == 503 || status == 504:
		return contracts.ErrorKind_ERROR_KIND_UNAVAILABLE, true
	case status >= 500 && status < 600:
		return contracts.ErrorKind_ERROR_KIND_UNAVAILABLE, false
	default:
		return contracts.ErrorKind_ERROR_KIND_UNSPECIFIED, false
	}
}

// kindName turns ERROR_KIND_NOT_FOUND into not_found.
func kindName(k contracts.ErrorKind) string {
	return strings.ToLower(strings.TrimPrefix(k.String(), "ERROR_KIND_"))
}
