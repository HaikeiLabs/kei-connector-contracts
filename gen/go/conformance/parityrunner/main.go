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

// Command parityrunner is the Go runner of the cross-language parity test
// (see scripts/parity/run.sh for the protocol). It calls each operation by
// its generated method name, so a new operation needs no change here; a new
// connector needs one entry in queries.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	"github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connection"
	"github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/github/v1/githubops"
	"github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/http_api/v1/httpapiops"
	"github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/linear/v1/linearops"
	contracts "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/contracts/v1"
)

// queries builds each HTTP connector's generated Queries.
var queries = map[string]func(connection.HTTP) any{
	"github":   func(c connection.HTTP) any { return githubops.New(c) },
	"http_api": func(c connection.HTTP) any { return httpapiops.New(c) },
	"linear":   func(c connection.HTTP) any { return linearops.New(c) },
}

type parityCase struct {
	Op          string          `json:"op"`
	Connector   string          `json:"connector"`
	Method      string          `json:"method"`
	RequestType string          `json:"request_type"`
	Request     json.RawMessage `json:"request"`
}

type result struct {
	Op       string          `json:"op"`
	Response json.RawMessage `json:"response,omitempty"`
	Error    string          `json:"error,omitempty"`
}

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "usage: parityrunner <cases.json> <base url> <token>")
		os.Exit(2)
	}
	b, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	var cases []parityCase
	if err := json.Unmarshal(b, &cases); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	enc := json.NewEncoder(os.Stdout)
	for _, c := range cases {
		resp, err := call(c, os.Args[2], os.Args[3])
		r := result{Op: c.Op}
		if err != nil {
			r.Error = err.Error()
		} else {
			r.Response = resp
		}
		if err := enc.Encode(r); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

func call(c parityCase, base, token string) (json.RawMessage, error) {
	build, ok := queries[c.Connector]
	if !ok {
		return nil, fmt.Errorf("parityrunner: no Go queries for connector %s", c.Connector)
	}
	conn, err := connection.NewHTTP(base+"/"+c.Op, &contracts.Connection{Name: "parity", Kind: contracts.ConnectionKind_CONNECTION_KIND_HTTP_BEARER}, token, http.DefaultClient)
	if err != nil {
		return nil, err
	}
	method := reflect.ValueOf(build(conn)).MethodByName(c.Method)
	if !method.IsValid() {
		return nil, fmt.Errorf("parityrunner: %s has no method %s", c.Connector, c.Method)
	}
	mt, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(c.RequestType))
	if err != nil {
		return nil, err
	}
	req := mt.New().Interface()
	if err := protojson.Unmarshal(c.Request, req); err != nil {
		return nil, fmt.Errorf("parityrunner: request fixture: %w", err)
	}
	out := method.Call([]reflect.Value{reflect.ValueOf(context.Background()), reflect.ValueOf(req)})
	if errV := out[1].Interface(); errV != nil {
		return nil, errV.(error)
	}
	b, err := protojson.Marshal(out[0].Interface().(proto.Message))
	return b, err
}
