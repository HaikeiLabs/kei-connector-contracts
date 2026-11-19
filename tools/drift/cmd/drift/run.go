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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	contractspb "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/contracts/v1"
)

// Options configures one drift run.
type Options struct {
	// ManifestPath is the path to connector-manifest.v1.json.
	ManifestPath string
	// RequestsDir holds per-connector request files <id>.json.
	RequestsDir string
	// Only restricts the run to one connector id (empty: all).
	Only string
	// Env is the environment to resolve KEI_SANDBOX_* credentials from.
	Env map[string]string
	// BaseURLs overrides a connector's manifest base_url by id (tests).
	BaseURLs map[string]string
	// Out receives the report lines.
	Out io.Writer
}

// envPrefix is the prefix of the sandbox credential variables:
// KEI_SANDBOX_<CONNECTOR>_<FIELD>, e.g. KEI_SANDBOX_GITHUB_TOKEN.
const envPrefix = "KEI_SANDBOX_"

const skipNoCreds = "no sandbox credentials"

// Manifest is the subset of kei.connector-manifest/v1 the drift check needs.
type manifest struct {
	Connectors []manifestConnector `json:"connectors"`
}

type manifestConnector struct {
	ID           string            `json:"id"`
	BaseURL      string            `json:"base_url"`
	Connections  []manifestConnect `json:"connections"`
	Capabilities []struct {
		ID    string `json:"id"`
		Class string `json:"class"`
	} `json:"capabilities"`
}

type manifestConnect struct {
	Name        string          `json:"name"`
	Kind        string          `json:"kind"`
	SetupFields []manifestField `json:"setup_fields"`
}

type manifestField struct {
	Name   string `json:"name"`
	Secret bool   `json:"secret"`
}

// requestFile is testenv/drift/<connector>.json: the safe READ operations to
// run against the sandbox.
type requestFile struct {
	Requests []requestSpec `json:"requests"`
}

type requestSpec struct {
	Op   string          `json:"op"`
	Args json.RawMessage `json:"args"`
}

// credential is one resolved sandbox connection: everything needed to build
// a connection.HTTP, plus the raw values for redaction. HeaderName is left
// empty: the manifest does not carry a header name, so connection.NewHTTP
// applies its DefaultHeaderName.
type credential struct {
	Kind      contractspb.ConnectionKind
	Value     string
	RawValues []string // every credential component, for redaction only
}

// collectCredentialValues gathers every KEI_SANDBOX_* environment value
// (base-URL overrides excluded) so the reporter can redact them. Over- and
// under-redaction both err on the side of never printing a credential.
func collectCredentialValues(m manifest, env map[string]string) []string {
	var vals []string
	for k, v := range env {
		if strings.HasPrefix(k, envPrefix) && !strings.HasSuffix(k, "_BASE_URL") {
			vals = append(vals, v)
		}
	}
	return vals
}

// Run executes the drift check and returns the process exit code: 0 when
// every run operation passed, 1 when any op drifted or errored.
func Run(ctx context.Context, o Options) int {
	raw, err := os.ReadFile(o.ManifestPath)
	if err != nil {
		fmt.Fprintf(o.Out, "drift: cannot read manifest: %v\n", err)
		return 1
	}
	var m manifest
	if err := json.Unmarshal(raw, &m); err != nil {
		fmt.Fprintf(o.Out, "drift: cannot parse manifest: %v\n", err)
		return 1
	}

	rep := &reporter{out: o.Out}
	rep.creds = collectCredentialValues(m, o.Env)
	for _, c := range m.Connectors {
		if o.Only != "" && c.ID != o.Only {
			continue
		}
		runConnector(ctx, o, c, rep)
	}
	fmt.Fprintf(o.Out, "drift: %d pass, %d drift, %d error, %d skip\n", rep.nPass, rep.nDrift, rep.nErrs, rep.nSkips)
	if rep.failed {
		return 1
	}
	return 0
}

func runConnector(ctx context.Context, o Options, c manifestConnector, rep *reporter) {
	// The request file names the safe READ operations for this connector.
	raw, err := os.ReadFile(o.RequestsDir + "/" + c.ID + ".json")
	if err != nil {
		rep.skip(c.ID, fmt.Sprintf("no request file testenv/drift/%s.json", c.ID))
		return
	}
	var rf requestFile
	if err := json.Unmarshal(raw, &rf); err != nil {
		rep.skip(c.ID, fmt.Sprintf("unparsable request file: %v", err))
		return
	}
	if len(rf.Requests) == 0 {
		rep.skip(c.ID, "empty request file")
		return
	}

	adapter, ok := adapters[c.ID]
	if !ok {
		rep.skip(c.ID, "no drift adapter")
		return
	}

	cred, ok := pickCredential(c, o.Env)
	if !ok {
		rep.skip(c.ID, skipNoCreds)
		return
	}

	base := c.BaseURL
	if override, ok := o.BaseURLs[c.ID]; ok {
		base = override
	}

	capture := &bodyCapture{base: http.DefaultTransport}
	client := &http.Client{Transport: capture}
	env := adapterEnv{
		ctx:      ctx,
		connID:   c.ID,
		base:     base,
		cred:     cred,
		client:   client,
		capture:  capture,
		requests: rf.Requests,
		rep:      rep,
	}
	adapter(env)
}

// pickCredential resolves the first connection of the connector whose secret
// setup fields are all present in the environment. The credential value is
// assembled per connection kind: bearer/header fields are the single secret
// value, basic is "user:password" in setup_fields order.
func pickCredential(c manifestConnector, env map[string]string) (credential, bool) {
	for _, conn := range c.Connections {
		var fields []string
		complete := true
		for _, f := range conn.SetupFields {
			if !f.Secret {
				continue
			}
			v := env[envPrefix+envName(c.ID)+"_"+envName(f.Name)]
			if v == "" {
				complete = false
				break
			}
			fields = append(fields, v)
		}
		if !complete || len(fields) == 0 {
			continue
		}
		var value string
		switch conn.Kind {
		case "http_bearer", "http_header":
			value = fields[0]
		case "http_basic":
			value = strings.Join(fields, ":")
		default:
			continue
		}
		return credential{
			Kind:      kindEnum(conn.Kind),
			Value:     value,
			RawValues: fields,
		}, true
	}
	return credential{}, false
}

// envName maps an id or field name to its environment-variable form:
// letters and digits only, upper-cased, dashes to underscores.
func envName(s string) string {
	return strings.ToUpper(strings.ReplaceAll(s, "-", "_"))
}

func kindEnum(kind string) contractspb.ConnectionKind {
	switch kind {
	case "http_header":
		return contractspb.ConnectionKind_CONNECTION_KIND_HTTP_HEADER
	case "http_basic":
		return contractspb.ConnectionKind_CONNECTION_KIND_HTTP_BASIC
	default:
		return contractspb.ConnectionKind_CONNECTION_KIND_HTTP_BEARER
	}
}

// reporter writes one line per outcome and counts them. Every line goes
// through redact so a credential value can never reach the output.
type reporter struct {
	out                          io.Writer
	creds                        []string
	nPass, nDrift, nErrs, nSkips int
	failed                       bool
}

func (r *reporter) line(format string, a ...any) {
	fmt.Fprintln(r.out, redact(fmt.Sprintf(format, a...), r.creds))
}

func (r *reporter) pass(conn, op string) {
	r.nPass++
	r.line("PASS %s", op)
}

func (r *reporter) drift(conn, op string, err error) {
	r.nDrift++
	r.failed = true
	r.line("DRIFT %s: %v", op, err)
}

func (r *reporter) err(conn, op string, err error) {
	r.nErrs++
	r.failed = true
	r.line("ERROR %s: %v", op, err)
}

func (r *reporter) skip(conn, why string) {
	r.nSkips++
	r.line("SKIP %s: %s", conn, why)
}

// redact replaces every known credential value in s with "[redacted]".
func redact(s string, creds []string) string {
	for _, c := range creds {
		if c != "" {
			s = strings.ReplaceAll(s, c, "[redacted]")
		}
	}
	return s
}
