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
// to the local database stack in testenv/compose.yaml through the SQL client
// that ships in each container (psql, mysql, sqlcmd), so the contract modules
// gain no database driver. Each query's result is converted into the
// RunQueryResponse JSON shape and unmarshaled into the generated types, which
// proves the contract accepts what the real engines return.
//
// Run it with the stack up:
//
//	podman compose -f testenv/compose.yaml --profile mssql up -d
//	go -C gen/go test -tags testenv ./conformance -run TestEnvDatabaseQueryRun -v
//
// The container CLI defaults to podman and can be overridden with the
// CONTAINER_CLI environment variable; CI sets CONTAINER_CLI=docker. The
// mssql service is amd64-only (see compose.yaml), so CI runs it on an amd64
// hosted runner.
package conformance

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"

	databasev1 "github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connectors/database/v1"
)

const (
	// nullSentinel marks NULL across the three engine clients: psql prints it
	// via -P null=..., and the mysql and mssql queries wrap nullable columns
	// in COALESCE(col, '__NULL__') — mysql --batch prints NULL as the literal
	// string NULL and sqlcmd has no NULL-marker option, so both route NULLs
	// through the sentinel in SQL themselves.
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

// checkNulls asserts the nulls case: two rows whose second column is NULL.
func checkNulls(t *testing.T, resp *databasev1.RunQueryResponse) {
	t.Helper()
	if got := len(resp.GetRows()); got != 2 {
		t.Fatalf("len(rows) = %d, want 2", got)
	}
	for i, row := range resp.GetRows() {
		// Type-assert the oneof, not GetNullValue: an unset oneof also reads
		// as the zero enum (NULL_VALUE), so GetNullValue would pass a non-null
		// cell vacuously.
		if _, ok := row.GetValues()[1].GetKind().(*structpb.Value_NullValue); !ok {
			t.Errorf("rows[%d][1] = %T, want null_value", i, row.GetValues()[1].GetKind())
		}
	}
}

// checkBigID asserts the bigint case: one row whose id is above 2^53 and
// therefore must travel as a string, plus the unicode account name.
func checkBigID(t *testing.T, resp *databasev1.RunQueryResponse) {
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
}

// checkDecimal asserts the numeric case: full decimal precision survives as
// strings.
func checkDecimal(t *testing.T, resp *databasev1.RunQueryResponse) {
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
}

func databaseQueryCases(engine string) []queryCase {
	switch engine {
	case "postgres":
		return []queryCase{
			{
				name: "nulls",
				// The nullable column goes through COALESCE so every engine
				// routes NULLs through the same sentinel (mysql batch prints
				// NULL as the literal string NULL; see runMySQLTSV).
				sql: "SELECT id, COALESCE(handle, '__NULL__') AS handle FROM accounts WHERE handle IS NULL ORDER BY id LIMIT 2",
				columns: []columnSpec{
					{"id", "int8"},
					{"handle", "text"},
				},
				check: checkNulls,
			},
			{
				name: "bigint_above_2_53",
				sql:  "SELECT id, name FROM accounts WHERE id > 9007199254740992 ORDER BY id",
				columns: []columnSpec{
					{"id", "int8"},
					{"name", "text"},
				},
				check: checkBigID,
			},
			{
				name: "numeric_precision",
				sql:  "SELECT 12345678901234567890.123456789012345678::numeric AS balance, 0.000000000000000001::numeric AS tiny",
				columns: []columnSpec{
					{"balance", "numeric"},
					{"tiny", "numeric"},
				},
				check: checkDecimal,
			},
		}
	case "mysql":
		return []queryCase{
			{
				name: "nulls",
				// Same COALESCE sentinel as the other engines (mysql batch
				// prints NULL as the literal string NULL).
				sql: "SELECT id, COALESCE(handle, '__NULL__') AS handle FROM accounts WHERE handle IS NULL ORDER BY id LIMIT 2",
				columns: []columnSpec{
					{"id", "bigint"},
					{"handle", "varchar"},
				},
				check: checkNulls,
			},
			{
				name: "bigint_above_2_53",
				sql:  "SELECT id, name FROM accounts WHERE id > 9007199254740992 ORDER BY id",
				columns: []columnSpec{
					{"id", "bigint"},
					{"name", "varchar"},
				},
				check: checkBigID,
			},
			{
				name: "numeric_precision",
				// MySQL DECIMAL tops out at precision 65; 20 integer digits plus
				// 18 fraction digits round-trip exactly through the cast.
				sql: "SELECT CAST('12345678901234567890.123456789012345678' AS DECIMAL(65,18)) AS balance, CAST('0.000000000000000001' AS DECIMAL(18,18)) AS tiny",
				columns: []columnSpec{
					{"balance", "decimal"},
					{"tiny", "decimal"},
				},
				check: checkDecimal,
			},
		}
	case "mssql":
		return []queryCase{
			{
				name: "nulls",
				// sqlcmd prints NULL as the literal string NULL with no marker
				// option, so the nullable column goes through COALESCE and
				// arrives as the nullSentinel literal.
				sql: "SELECT TOP 2 id, COALESCE(handle, '__NULL__') AS handle FROM accounts WHERE handle IS NULL ORDER BY id",
				columns: []columnSpec{
					{"id", "bigint"},
					{"handle", "nvarchar"},
				},
				check: checkNulls,
			},
			{
				name: "bigint_above_2_53",
				sql:  "SELECT id, name FROM accounts WHERE id > 9007199254740992 ORDER BY id",
				columns: []columnSpec{
					{"id", "bigint"},
					{"name", "nvarchar"},
				},
				check: checkBigID,
			},
			{
				name: "numeric_precision",
				// SQL Server DECIMAL tops out at precision 38: exactly the 20
				// integer plus 18 fraction digits of the balance literal.
				sql: "SELECT CAST('12345678901234567890.123456789012345678' AS DECIMAL(38,18)) AS balance, CAST('0.000000000000000001' AS DECIMAL(18,18)) AS tiny",
				columns: []columnSpec{
					{"balance", "decimal"},
					{"tiny", "decimal"},
				},
				check: checkDecimal,
			},
		}
	default:
		return nil
	}
}

// containerCLI returns the container CLI that drives compose. It defaults to
// podman (local development) and can be overridden with CONTAINER_CLI; CI
// exports CONTAINER_CLI=docker.
func containerCLI(t *testing.T) string {
	t.Helper()
	cli := os.Getenv("CONTAINER_CLI")
	if cli == "" {
		cli = "podman"
	}
	if _, err := exec.LookPath(cli); err != nil {
		t.Skipf("%s not installed; testenv proof skipped", cli)
	}
	return cli
}

// runInContainer executes argv inside a compose service and returns stdout.
func runInContainer(t *testing.T, cli, service string, argv ...string) string {
	t.Helper()
	composeAbs, err := filepath.Abs(composeFile)
	if err != nil {
		t.Fatalf("resolving compose file: %v", err)
	}
	args := append([]string{"compose", "-f", composeAbs, "exec", "-T", service}, argv...)
	cmd := exec.Command(cli, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("querying %s (is the stack up? %s compose -f testenv/compose.yaml --profile mssql up -d): %v: %s",
			service, cli, err, stderr.String())
	}
	return stdout.String()
}

// runPSQLCSV executes one query inside a postgres-family compose service and
// returns its CSV rows, header first.
func runPSQLCSV(t *testing.T, cli, service, database, sql string) [][]string {
	t.Helper()
	out := runInContainer(t, cli, service,
		"psql", "-U", "testenv", "-d", database, "--csv",
		"-P", "null="+nullSentinel, "-v", "ON_ERROR_STOP=1", "-c", sql)
	return parseCSV(t, out, ',')
}

// runMySQLTSV executes one query inside the mysql compose service and returns
// its rows, header first. It connects over TCP: during the image's first-boot
// init phase a temporary server answers on the socket, so a socket connection
// can reach a server that is about to be replaced by the real one (the
// healthcheck guards against the same race; see compose.yaml). mysql --batch
// prints tab-separated values and NULL as the literal string NULL (the \N
// marker is a mysqldump convention, not a client one), so any nullable column
// must arrive through COALESCE(col, nullSentinel) in the SQL itself.
func runMySQLTSV(t *testing.T, cli, service, database, sql string) [][]string {
	t.Helper()
	out := runInContainer(t, cli, service,
		"mysql", "-h", "127.0.0.1", "--protocol=TCP", "-u", "testenv", "-ptestenv",
		"--default-character-set=utf8mb4", database, "--batch", "--raw", "-e", sql)
	return parseCSV(t, out, '\t')
}

// runMSSQLCSV executes one query inside the mssql compose service and returns
// its rows without a header (-h -1). SET NOCOUNT ON keeps the "(N rows
// affected)" footer out of the output and -f 65001 forces UTF-8 output so
// unicode values survive; sqlcmd has no NULL-marker option, so the mssql
// queries wrap nullable columns in COALESCE(col, '__NULL__') themselves.
func runMSSQLCSV(t *testing.T, cli, service, sql string) [][]string {
	t.Helper()
	out := runInContainer(t, cli, service,
		"/opt/mssql-tools18/bin/sqlcmd", "-C", "-S", "localhost", "-U", "sa",
		"-P", "Testenv_local1", "-I", "-b", "-h", "-1", "-W", "-s", ",", "-f", "65001", "-Q",
		"SET NOCOUNT ON; "+sql)
	return parseCSV(t, out, ',')
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
	case "int2", "int4", "int8", "bigint":
		n, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			t.Fatalf("parsing %q as %s: %v", raw, typ, err)
		}
		if n > 1<<53 || n < -(1<<53) {
			return raw
		}
		return n
	case "numeric", "decimal":
		return raw
	case "bool":
		return raw == "t"
	default:
		return raw
	}
}

// toResponse validates the CSV header (psql and mysql print one) and converts
// the data rows into the RunQueryResponse JSON shape.
func toResponse(t *testing.T, records [][]string, columns []columnSpec) *databasev1.RunQueryResponse {
	t.Helper()
	if got, want := records[0], len(columns); len(got) != want {
		t.Fatalf("CSV header has %d columns, want %d: %v", len(got), want, got)
	}
	return toResponseRows(t, records[1:], columns)
}

// toResponseRows converts headerless rows (sqlcmd -h -1 prints no header)
// into the RunQueryResponse JSON shape and unmarshals them into the generated
// message.
func toResponseRows(t *testing.T, rows [][]string, columns []columnSpec) *databasev1.RunQueryResponse {
	t.Helper()
	shape := responseJSON{
		Columns:   make([]columnJSON, len(columns)),
		Rows:      make([][]any, 0, len(rows)),
		Truncated: false,
	}
	for i, col := range columns {
		shape.Columns[i] = columnJSON{Name: col.name, Type: col.typ}
	}
	for _, record := range rows {
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

// TestEnvDatabaseQueryRun proves database.query.run against the real database
// services from testenv/compose.yaml — postgres, supabase (Postgres), mysql
// and mssql — with the same three queries everywhere: NULLs, a bigint above
// 2^53 that must travel as a string, and a high-precision decimal that must
// travel as a string.
func TestEnvDatabaseQueryRun(t *testing.T) {
	cli := containerCLI(t)
	targets := []struct {
		service  string
		engine   string
		database string
	}{
		{"postgres", "postgres", "testenv"},
		{"supabase", "postgres", "supabase"},
		{"mysql", "mysql", "testenv"},
		{"mssql", "mssql", ""}, // seed sidecar applies to master; no default database
	}
	for _, target := range targets {
		t.Run(target.service, func(t *testing.T) {
			for _, tc := range databaseQueryCases(target.engine) {
				t.Run(tc.name, func(t *testing.T) {
					var resp *databasev1.RunQueryResponse
					switch target.engine {
					case "postgres":
						resp = toResponse(t, runPSQLCSV(t, cli, target.service, target.database, tc.sql), tc.columns)
					case "mysql":
						resp = toResponse(t, runMySQLTSV(t, cli, target.service, target.database, tc.sql), tc.columns)
					case "mssql":
						// -h -1 prints no header row, so the rows convert
						// directly.
						resp = toResponseRows(t, runMSSQLCSV(t, cli, target.service, tc.sql), tc.columns)
					default:
						t.Fatalf("unknown engine %q", target.engine)
					}
					if got := len(resp.GetColumns()); got != len(tc.columns) {
						t.Fatalf("len(columns) = %d, want %d", got, len(tc.columns))
					}
					tc.check(t, resp)
				})
			}
		})
	}
}
