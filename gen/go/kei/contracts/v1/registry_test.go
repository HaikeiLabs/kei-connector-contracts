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

package contracts_test

import (
	"reflect"
	"regexp"
	"testing"

	"google.golang.org/protobuf/encoding/protojson"

	githubv1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/github/v1"
	contracts "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/contracts/v1"
)

// The GitHub sample registers its descriptors on init; importing it here is
// what puts the sample operation into protoregistry.GlobalFiles.
var _ = githubv1.GetIssueRequest{}

var dottedIDPattern = regexp.MustCompile(`^[a-z0-9_]+(\.[a-z0-9_]+){2,}$`)

func TestEveryOperationHasUniqueDottedID(t *testing.T) {
	ops := contracts.All()
	if len(ops) == 0 {
		t.Fatal("All() returned no operations; the sample operation did not register")
	}

	seen := make(map[string]bool, len(ops))
	for _, op := range ops {
		if !dottedIDPattern.MatchString(op.GetId()) {
			t.Errorf("operation id %q does not match the dotted-id pattern", op.GetId())
		}
		if seen[op.GetId()] {
			t.Errorf("duplicate operation id %q", op.GetId())
		}
		seen[op.GetId()] = true

		connector := op.GetConnector()
		if connector == "" {
			t.Errorf("operation %q has no connector", op.GetId())
			continue
		}
		id := op.GetId()
		dot := -1
		for i, r := range id {
			if r == '.' {
				dot = i
				break
			}
		}
		if dot < 0 || id[:dot] != connector {
			t.Errorf("operation %q connector %q is not the first id segment", op.GetId(), connector)
		}
	}
}

func TestWriteOperationsDeclareScopes(t *testing.T) {
	for _, op := range contracts.All() {
		if op.GetClass() != contracts.OperationClass_OPERATION_CLASS_WRITE {
			continue
		}
		if len(op.GetScopes()) == 0 {
			t.Errorf("write operation %q declares no scopes", op.GetId())
		}
	}
}

func TestLookupSample(t *testing.T) {
	op, ok := contracts.Lookup("github.issues.get")
	if !ok {
		t.Fatal("Lookup(github.issues.get) not found")
	}
	if op.GetHttpMethod() != "GET" {
		t.Errorf("http_method = %q, want GET", op.GetHttpMethod())
	}
}

func TestSampleUnmarshal(t *testing.T) {
	data := []byte(`{"id": 9007199254740993, "number": 7, "title": "t", "state": "open", "body": null, "html_url": "https://example.invalid/x", "extra_field": 1}`)

	var msg githubv1.GetIssueResponse
	if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(data, &msg); err != nil {
		t.Fatalf("protojson unmarshal failed: %v", err)
	}

	if got := msg.GetId(); got != 9007199254740993 {
		t.Errorf("id = %d, want 9007199254740993 (int64 must not go through float64)", got)
	}
	if got := msg.GetNumber(); got != 7 {
		t.Errorf("number = %d, want 7", got)
	}
	if !reflect.DeepEqual(msg.GetTitle(), "t") || msg.GetState() != "open" {
		t.Errorf("title = %q, state = %q", msg.GetTitle(), msg.GetState())
	}
	if msg.Body != nil {
		t.Errorf("body = %v, want unset (JSON null)", *msg.Body)
	}
	if msg.GetHtmlUrl() != "https://example.invalid/x" {
		t.Errorf("html_url = %q", msg.GetHtmlUrl())
	}
}
