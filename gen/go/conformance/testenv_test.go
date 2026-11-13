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

//go:build testenv

// Testenv-only proof for database.query.run.
//
// This file is behind the "testenv" build tag so normal CI skips it. It talks
// to the local database stack in testenv/compose.yaml through the psql that
// ships in each container, so the contract modules gain no database driver.
// Each query's CSV result is converted into the RunQueryResponse JSON shape
// and unmarshaled into the generated types, which proves the contract accepts
// what the real engine returns.
//
// Run it with the stack up:
//
//	podman compose -f testenv/compose.yaml up -d postgres supabase
//	go -C gen/go test -tags testenv ./conformance -run TestEnvDatabaseQueryRun -v
package conformance

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"

	databasev1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/database/v1"
)

const (
	// psql prints this instead of an empty field, so NULL survives CSV.
	nullSentinel = "__NULL__"
	// composeFile is relative to this package directory (gen/go/conformance).
	composeFile = "../../../testenv/compose.yaml"
)

// columnSpec mirrors RunQueryResponse.Column: the engine's column name and
// type name, used to decide how the cell enters the JSON shape.
type columnSpec struct {
	name string
	typ  string
}

type queryCase struct {
	name    string
	sql     string
	columns []columnSpec
	check   func(t *testing.T, resp *databasev1.RunQueryResponse)
}

type responseJSON struct {
	Columns   []columnJSON `json:"columns"`
	Rows      [][]any      `json:"rows"`
	Truncated bool         `json:"truncated"`
}

type columnJSON struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func databaseQueryCases() []queryCase {
	return []queryCase{
		{
			name: "nulls",
			sql:  "SELECT id, handle FROM accounts WHERE handle IS NULL ORDER BY id LIMIT 2",
			columns: []columnSpec{
				{"id", "int8"},
				{"handle", "text"},
			},
			check: func(t *testing.T, resp *databasev1.RunQueryResponse) {
				t.Helper()
				if got := len(resp.GetRows()); got != 2 {
					t.Fatalf("len(rows) = %d, want 2", got)
				}
				for i, row := range resp.GetRows() {
					if got := row.GetValues()[1].GetNullValue(); got != structpb.NullValue_NULL_VALUE {
						t.Errorf("rows[%d][1] = %v, want NULL", i, got)
					}
				}
			},
		},
		{
			name: "bigint_above_2_53",
			sql:  "SELECT id, name FROM accounts WHERE id > 9007199254740992 ORDER BY id",
			columns: []columnSpec{
				{"id", "int8"},
				{"name", "text"},
			},
			check: func(t *testing.T, resp *databasev1.RunQueryResponse) {
				t.Helper()
				if got := len(resp.GetRows()); got != 1 {
					t.Fatalf("len(rows) = %d, want 1", got)
				}
				row := resp.GetRows()[0].GetValues()
				if got := row[0].GetStringValue(); got != "9007199254740993" {
					t.Errorf("rows[0][0] = %q, want 9007199254740993 as a string", got)
				}
				if got := row[1].GetStringValue(); got != "Big-ID Account — ünïcodé ✓" {
					t.Errorf("rows[0][1] = %q, want the unicode account name", got)
				}
			},
		},
		{
			name: "numeric_precision",
			sql:  "SELECT 12345678901234567890.123456789012345678::numeric AS balance, 0.000000000000000001::numeric AS tiny",
			columns: []columnSpec{
				{"balance", "numeric"},
				{"tiny", "numeric"},
			},
			check: func(t *testing.T, resp *databasev1.RunQueryResponse) {
				t.Helper()
				if got := len(resp.GetRows()); got != 1 {
					t.Fatalf("len(rows) = %d, want 1", got)
				}
				row := resp.GetRows()[0].GetValues()
				if got := row[0].GetStringValue(); got != "12345678901234567890.123456789012345678" {
					t.Errorf("rows[0][0] = %q, want NUMERIC precision preserved", got)
				}
				if got := row[1].GetStringValue(); got != "0.000000000000000001" {
					t.Errorf("rows[0][1] = %q, want the fraction preserved", got)
				}
			},
		},
	}
}

// runPSQLCSV executes one query inside the compose service and returns its
// CSV rows, header first.
func runPSQLCSV(t *testing.T, service, database, sql string) [][]string {
	t.Helper()
	if _, err := exec.LookPath("podman"); err != nil {
		t.Skip("podman not installed; testenv proof skipped")
	}
	composeAbs, err := filepath.Abs(composeFile)
	if err != nil {
		t.Fatalf("resolving compose file: %v", err)
	}
	cmd := exec.Command("podman", "compose", "-f", composeAbs, "exec", "-T", service,
		"psql", "-U", "testenv", "-d", database, "--csv",
		"-P", "null="+nullSentinel, "-v", "ON_ERROR_STOP=1", "-c", sql)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("querying %s (is the stack up? podman compose -f testenv/compose.yaml up -d postgres supabase): %v: %s", service, err, stderr.String())
	}
	records, err := csv.NewReader(strings.NewReader(stdout.String())).ReadAll()
	if err != nil {
		t.Fatalf("parsing CSV from %s: %v", service, err)
	}
	if len(records) == 0 {
		t.Fatalf("no CSV output from %s", service)
	}
	return records
}

// convertCell models how a runtime fills a google.protobuf.Value: null stays
// null, int64 above 2^53 becomes a string (JSON numbers are float64), and
// NUMERIC becomes a string so precision and scale survive.
func convertCell(t *testing.T, raw, typ string) any {
	t.Helper()
	if raw == nullSentinel {
		return nil
	}
	switch typ {
	case "int2", "int4", "int8":
		n, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			t.Fatalf("parsing %q as %s: %v", raw, typ, err)
		}
		if n > 1<<53 || n < -(1<<53) {
			return raw
		}
		return n
	case "numeric":
		return raw
	case "bool":
		return raw == "t"
	default:
		return raw
	}
}

// toResponse converts CSV records into the RunQueryResponse JSON shape and
// unmarshals them into the generated message.
func toResponse(t *testing.T, records [][]string, columns []columnSpec) *databasev1.RunQueryResponse {
	t.Helper()
	if got, want := records[0], len(columns); len(got) != want {
		t.Fatalf("CSV header has %d columns, want %d: %v", len(got), want, got)
	}
	shape := responseJSON{
		Columns:   make([]columnJSON, len(columns)),
		Rows:      make([][]any, 0, len(records)-1),
		Truncated: false,
	}
	for i, col := range columns {
		shape.Columns[i] = columnJSON{Name: col.name, Type: col.typ}
	}
	for _, record := range records[1:] {
		row := make([]any, len(columns))
		for i, col := range columns {
			row[i] = convertCell(t, record[i], col.typ)
		}
		shape.Rows = append(shape.Rows, row)
	}

	payload, err := json.Marshal(shape)
	if err != nil {
		t.Fatalf("marshaling response shape: %v", err)
	}
	t.Logf("response JSON: %s", payload)

	resp := &databasev1.RunQueryResponse{}
	if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(payload, resp); err != nil {
		t.Fatalf("generated type did not accept the converted response: %v", err)
	}
	return resp
}

// TestEnvDatabaseQueryRun proves database.query.run against the real Postgres
// and Supabase (Postgres) services from testenv/compose.yaml: one query with
// NULLs, one with a bigint above 2^53, one with NUMERIC precision.
func TestEnvDatabaseQueryRun(t *testing.T) {
	targets := []struct {
		service  string
		database string
	}{
		{"postgres", "testenv"},
		{"supabase", "supabase"},
	}
	for _, target := range targets {
		t.Run(target.service, func(t *testing.T) {
			for _, tc := range databaseQueryCases() {
				t.Run(tc.name, func(t *testing.T) {
					records := runPSQLCSV(t, target.service, target.database, tc.sql)
					resp := toResponse(t, records, tc.columns)
					if got := len(resp.GetColumns()); got != len(tc.columns) {
						t.Fatalf("len(columns) = %d, want %d", got, len(tc.columns))
					}
					tc.check(t, resp)
				})
			}
		})
	}
}
