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

package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"
)

// fakeDriver is a stdlib database/sql driver used to exercise the package
// without a real Postgres. It returns controlled rows and can fail on Open.
type fakeDriver struct {
	cols    []string
	rows    [][]any
	openErr error
	dsn     *string
}

func (d *fakeDriver) Open(name string) (driver.Conn, error) {
	if d.dsn != nil {
		*d.dsn = name
	}
	if d.openErr != nil {
		return nil, d.openErr
	}
	return &fakeConn{d: d}, nil
}

type fakeConn struct{ d *fakeDriver }

func (c *fakeConn) Prepare(query string) (driver.Stmt, error) { return &fakeStmt{c: c}, nil }
func (c *fakeConn) Close() error                              { return nil }
func (c *fakeConn) Begin() (driver.Tx, error)                 { return fakeTx{}, nil }

func (c *fakeConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	return fakeTx{}, nil
}

func (c *fakeConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	return driver.RowsAffected(0), nil
}

func (c *fakeConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	return &fakeRows{cols: c.d.cols, rows: c.d.rows}, nil
}

type fakeTx struct{}

func (fakeTx) Commit() error   { return nil }
func (fakeTx) Rollback() error { return nil }

type fakeStmt struct{ c *fakeConn }

func (s *fakeStmt) Close() error  { return nil }
func (s *fakeStmt) NumInput() int { return 0 }
func (s *fakeStmt) Exec(args []driver.Value) (driver.Result, error) {
	return driver.RowsAffected(0), nil
}
func (s *fakeStmt) Query(args []driver.Value) (driver.Rows, error) {
	return &fakeRows{cols: s.c.d.cols, rows: s.c.d.rows}, nil
}

type fakeRows struct {
	cols []string
	rows [][]any
	i    int
}

func (r *fakeRows) Columns() []string { return r.cols }
func (r *fakeRows) Close() error      { return nil }
func (r *fakeRows) Next(dest []driver.Value) error {
	if r.i >= len(r.rows) {
		return io.EOF
	}
	row := r.rows[r.i]
	for i := range dest {
		if i < len(row) {
			dest[i] = row[i]
		}
	}
	r.i++
	return nil
}

var fakeDriverCounter int

func registerFake(t *testing.T, d *fakeDriver) string {
	t.Helper()
	fakeDriverCounter++
	name := fmt.Sprintf("fake-driver-%d", fakeDriverCounter)
	sql.Register(name, d)
	return name
}

func TestOpenRejectsSSLModeDisable(t *testing.T) {
	dsn := "host=localhost port=5432 user=ro password=topsecret dbname=app sslmode=disable"
	_, err := Open(context.Background(), dsn)
	if err == nil {
		t.Fatal("expected error for sslmode=disable, got nil")
	}
	if !strings.Contains(err.Error(), "sslmode=disable") {
		t.Fatalf("error should mention sslmode=disable, got: %v", err)
	}
	if strings.Contains(err.Error(), "topsecret") {
		t.Fatalf("error leaked the password: %v", err)
	}
	if strings.Contains(err.Error(), "host=localhost") {
		t.Fatalf("error leaked the DSN: %v", err)
	}
}

func TestOpenAllowsInsecure(t *testing.T) {
	name := registerFake(t, &fakeDriver{cols: []string{"id"}, rows: [][]any{{"1"}}})
	dsn := "host=localhost sslmode=disable"
	db, err := Open(context.Background(), dsn, WithDriver(name), WithAllowInsecure())
	if err != nil {
		t.Fatalf("expected success with AllowInsecure, got: %v", err)
	}
	db.Close()
}

func TestErrorRedaction(t *testing.T) {
	dsn := "host=localhost password=supersecret dbname=app"
	d := &fakeDriver{openErr: fmt.Errorf("pq: could not connect to host=localhost user=ro password=supersecret")}
	name := registerFake(t, d)
	_, err := Open(context.Background(), dsn, WithDriver(name))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if strings.Contains(err.Error(), "supersecret") {
		t.Fatalf("error leaked the password: %v", err)
	}
	if !strings.Contains(err.Error(), "***") {
		t.Fatalf("expected redacted placeholder in error, got: %v", err)
	}
}

func TestRowCap(t *testing.T) {
	rows := make([][]any, 0, 5)
	for i := 0; i < 5; i++ {
		rows = append(rows, []any{fmt.Sprintf("id%d", i), fmt.Sprintf("name%d", i)})
	}
	name := registerFake(t, &fakeDriver{cols: []string{"id", "name"}, rows: rows})
	db, err := Open(context.Background(), "host=localhost", WithDriver(name), WithRowCap(3))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()
	got, err := db.Query(context.Background(), "SELECT id, name FROM t")
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 rows (cap), got %d", len(got))
	}
	if got[0]["id"] != "id0" || got[2]["name"] != "name2" {
		t.Fatalf("unexpected rows: %v", got)
	}
}

func TestParseURLDSN(t *testing.T) {
	cfg, err := Parse("postgres://ro:secret@localhost:5432/app?sslmode=disable")
	if err != nil {
		t.Fatalf("parse URL DSN: %v", err)
	}
	if cfg.Get("host") != "localhost" {
		t.Fatalf("host = %q, want localhost", cfg.Get("host"))
	}
	if cfg.Get("port") != "5432" {
		t.Fatalf("port = %q, want 5432", cfg.Get("port"))
	}
	if cfg.Get("user") != "ro" {
		t.Fatalf("user = %q, want ro", cfg.Get("user"))
	}
	if cfg.Password() != "secret" {
		t.Fatalf("password = %q, want secret", cfg.Password())
	}
	if cfg.Get("dbname") != "app" {
		t.Fatalf("dbname = %q, want app", cfg.Get("dbname"))
	}
	if cfg.SSLMode() != "disable" {
		t.Fatalf("sslmode = %q, want disable", cfg.SSLMode())
	}
}

func TestOpenAppliesReadOnlyStartupParams(t *testing.T) {
	var dsn string
	d := &fakeDriver{cols: []string{"id"}, rows: [][]any{{"1"}}, dsn: &dsn}
	name := registerFake(t, d)
	db, err := Open(context.Background(), "host=localhost", WithDriver(name), WithStatementTimeout(15*time.Second))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()
	if !strings.Contains(dsn, "default_transaction_read_only=on") {
		t.Fatalf("DSN should include default_transaction_read_only=on, got: %s", dsn)
	}
	if !strings.Contains(dsn, "statement_timeout=15000") {
		t.Fatalf("DSN should include statement_timeout=15000, got: %s", dsn)
	}
}

func TestOpenAppliesReadOnlyStartupParamsURL(t *testing.T) {
	var dsn string
	d := &fakeDriver{cols: []string{"id"}, rows: [][]any{{"1"}}, dsn: &dsn}
	name := registerFake(t, d)
	db, err := Open(context.Background(), "postgres://ro@localhost/app?sslmode=disable",
		WithDriver(name), WithAllowInsecure())
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()
	if !strings.Contains(dsn, "default_transaction_read_only=on") {
		t.Fatalf("URL DSN should include default_transaction_read_only=on, got: %s", dsn)
	}
	if !strings.Contains(dsn, "statement_timeout=30000") {
		t.Fatalf("URL DSN should include the default statement_timeout=30000, got: %s", dsn)
	}
}
