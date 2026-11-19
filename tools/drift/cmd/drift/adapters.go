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

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connection"
	github "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/github/v1"
	"github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/github/v1/githubops"
	linear "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/linear/v1"
	"github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/linear/v1/linearops"
	contractspb "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/contracts/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// adapter runs one connector's request file through the generated ops and
// reports an outcome per operation.
type adapter func(env adapterEnv)

// adapters maps a connector id to its adapter. Connectors without an adapter
// are skipped, not failed: the nightly check grows adapter by adapter.
var adapters = map[string]adapter{
	"github": runGithub,
	"linear": runLinear,
}

// adapterEnv is everything one adapter needs to run its connector's
// operations and report the outcomes.
type adapterEnv struct {
	ctx      context.Context
	connID   string
	base     string
	cred     credential
	client   *http.Client
	capture  *bodyCapture
	requests []requestSpec
	rep      *reporter
}

// newHTTPConn builds the authenticated connection the generated ops run on.
func newHTTPConn(env adapterEnv) (connection.HTTP, error) {
	connProto := &contractspb.Connection{
		Name: env.connID,
		Kind: env.cred.Kind,
	}
	return connection.NewHTTP(env.base, connProto, env.cred.Value, env.client)
}

// opFunc runs one generated operation over the connection and returns the
// call error plus a fresh response message for the strict re-unmarshal.
type opFunc func(ctx context.Context, args json.RawMessage) (proto.Message, error)

// check runs one request spec: invoke the op, then strict-unmarshal the
// captured raw body with DiscardUnknown=false. A call error is reported as
// ERROR (the sandbox answered but not contract-clean), a strict unmarshal
// failure as DRIFT, success as PASS. Only the call error is shown when both
// fire: an error body is not the response schema.
func (env adapterEnv) check(op string, fn opFunc) {
	resp, callErr := fn(env.ctx, argsOf(env.requests, op))
	if callErr != nil {
		env.rep.err(env.connID, op, callErr)
		return
	}
	if resp == nil {
		env.rep.err(env.connID, op, fmt.Errorf("no response message for strict check"))
		return
	}
	if err := strictUnmarshal(env.capture.take(), resp); err != nil {
		env.rep.drift(env.connID, op, err)
		return
	}
	env.rep.pass(env.connID, op)
}

func argsOf(reqs []requestSpec, op string) json.RawMessage {
	for _, r := range reqs {
		if r.Op == op {
			return r.Args
		}
	}
	return nil
}

// strictUnmarshal re-unmarshals the raw provider body into the operation's
// response type with DiscardUnknown=false. The generated ops unmarshal
// leniently (connection.CallHTTP hard-codes DiscardUnknown=true), so this is
// the only place drift becomes visible.
func strictUnmarshal(body []byte, resp proto.Message) error {
	if err := (protojson.UnmarshalOptions{}).Unmarshal(body, resp); err != nil {
		return fmt.Errorf("strict unmarshal: %w", err)
	}
	return nil
}

// unmarshalArgs fills an operation's request from the request file's args.
// Request files are repo-owned; a malformed args value is reported, not
// ignored. Unknown args are dropped: the request files are maintained by
// hand and may carry more context than one op needs.
func unmarshalArgs(args json.RawMessage, req proto.Message) error {
	if len(args) == 0 {
		return nil
	}
	return (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(args, req)
}

// githubOps maps the github READ operation ids to calls over a connection.
func githubOps(conn connection.HTTP) map[string]opFunc {
	q := githubops.New(conn)
	return map[string]opFunc{
		"github.actions.runs.list": func(ctx context.Context, args json.RawMessage) (proto.Message, error) {
			req := new(github.ListWorkflowRunsRequest)
			if err := unmarshalArgs(args, req); err != nil {
				return nil, err
			}
			_, err := q.ListWorkflowRuns(ctx, req)
			return new(github.ListWorkflowRunsResponse), err
		},
		"github.git.refs.get": func(ctx context.Context, args json.RawMessage) (proto.Message, error) {
			req := new(github.GetRefRequest)
			if err := unmarshalArgs(args, req); err != nil {
				return nil, err
			}
			_, err := q.GetRef(ctx, req)
			return new(github.GetRefResponse), err
		},
		"github.issues.get": func(ctx context.Context, args json.RawMessage) (proto.Message, error) {
			req := new(github.GetIssueRequest)
			if err := unmarshalArgs(args, req); err != nil {
				return nil, err
			}
			_, err := q.GetIssue(ctx, req)
			return new(github.GetIssueResponse), err
		},
		"github.issues.list": func(ctx context.Context, args json.RawMessage) (proto.Message, error) {
			req := new(github.ListIssuesRequest)
			if err := unmarshalArgs(args, req); err != nil {
				return nil, err
			}
			_, err := q.ListIssues(ctx, req)
			return new(github.ListIssuesResponse), err
		},
		"github.pulls.get": func(ctx context.Context, args json.RawMessage) (proto.Message, error) {
			req := new(github.GetPullRequestRequest)
			if err := unmarshalArgs(args, req); err != nil {
				return nil, err
			}
			_, err := q.GetPullRequest(ctx, req)
			return new(github.GetPullRequestResponse), err
		},
		"github.pulls.list": func(ctx context.Context, args json.RawMessage) (proto.Message, error) {
			req := new(github.ListPullsRequest)
			if err := unmarshalArgs(args, req); err != nil {
				return nil, err
			}
			_, err := q.ListPulls(ctx, req)
			return new(github.ListPullsResponse), err
		},
		"github.pulls.reviews.list": func(ctx context.Context, args json.RawMessage) (proto.Message, error) {
			req := new(github.ListPullReviewsRequest)
			if err := unmarshalArgs(args, req); err != nil {
				return nil, err
			}
			_, err := q.ListReviews(ctx, req)
			return new(github.ListPullReviewsResponse), err
		},
	}
}

// linearOps maps the linear READ operation ids to calls over a connection.
func linearOps(conn connection.HTTP) map[string]opFunc {
	q := linearops.New(conn)
	return map[string]opFunc{
		"linear.issues.get": func(ctx context.Context, args json.RawMessage) (proto.Message, error) {
			req := new(linear.GetIssueRequest)
			if err := unmarshalArgs(args, req); err != nil {
				return nil, err
			}
			_, err := q.GetIssue(ctx, req)
			return new(linear.GetIssueResponse), err
		},
		"linear.issues.search": func(ctx context.Context, args json.RawMessage) (proto.Message, error) {
			req := new(linear.SearchIssuesRequest)
			if err := unmarshalArgs(args, req); err != nil {
				return nil, err
			}
			_, err := q.SearchIssues(ctx, req)
			return new(linear.SearchIssuesResponse), err
		},
		"linear.teams.get": func(ctx context.Context, args json.RawMessage) (proto.Message, error) {
			req := new(linear.GetTeamRequest)
			if err := unmarshalArgs(args, req); err != nil {
				return nil, err
			}
			_, err := q.GetTeam(ctx, req)
			return new(linear.GetTeamResponse), err
		},
	}
}

// runOps executes the request file against one op table.
func runOps(env adapterEnv, table map[string]opFunc) {
	for _, spec := range env.requests {
		fn, ok := table[spec.Op]
		if !ok {
			env.rep.skip(env.connID, "no drift adapter for op "+spec.Op)
			continue
		}
		env.check(spec.Op, fn)
	}
}

func runGithub(env adapterEnv) {
	conn, err := newHTTPConn(env)
	if err != nil {
		env.rep.err(env.connID, env.connID+".connect", err)
		return
	}
	runOps(env, githubOps(conn))
}

func runLinear(env adapterEnv) {
	conn, err := newHTTPConn(env)
	if err != nil {
		env.rep.err(env.connID, env.connID+".connect", err)
		return
	}
	runOps(env, linearOps(conn))
}

// bodyCapture records the last response body that flowed through the
// connection's http.Client, so the drift check can re-unmarshal the exact
// bytes the provider sent.
type bodyCapture struct {
	base http.RoundTripper
	mu   sync.Mutex
	last []byte
}

func (c *bodyCapture) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := c.base.RoundTrip(req)
	if err != nil {
		return resp, err
	}
	body, rerr := io.ReadAll(resp.Body)
	resp.Body.Close()
	if rerr == nil {
		c.mu.Lock()
		c.last = body
		c.mu.Unlock()
		resp.Body = io.NopCloser(bytes.NewReader(body))
	}
	return resp, rerr
}

func (c *bodyCapture) take() []byte {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.last
}
