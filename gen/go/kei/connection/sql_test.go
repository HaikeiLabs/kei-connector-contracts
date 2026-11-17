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

package connection_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"io"
	"reflect"
	"sync"
	"testing"

	"github.com/HaikeiLabs/kei-connector-contracts/gen/go/kei/connection"
)

// fakeDriver serves a fixed result set and records the query and arguments.
type fakeDriver struct {
	mu    sync.Mutex
	query string
	args  []driver.NamedValue
	rows  [][]driver.Value
}

func (d *fakeDriver) Open(string) (driver.Conn, error) { return &fakeConn{d: d}, nil }

type fakeConn struct{ d *fakeDriver }

func (c *fakeConn) Prepare(string) (driver.Stmt, error) { return nil, driver.ErrSkip }
func (c *fakeConn) Close() error                        { return nil }
func (c *fakeConn) Begin() (driver.Tx, error)           { return nil, driver.ErrSkip }
func (c *fakeConn) QueryContext(_ context.Context, q string, args []driver.NamedValue) (driver.Rows, error) {
	c.d.mu.Lock()
	defer c.d.mu.Unlock()
	c.d.query, c.d.args = q, args
	return &fakeRows{rows: c.d.rows}, nil
}

type fakeRows struct {
	rows [][]driver.Value
	i    int
}

func (r *fakeRows) Columns() []string { return []string{"region", "total"} }
func (r *fakeRows) Close() error      { return nil }
func (r *fakeRows) Next(dest []driver.Value) error {
	if r.i >= len(r.rows) {
		return io.EOF
	}
	copy(dest, r.rows[r.i])
	r.i++
	return nil
}

var registerOnce sync.Once
var theDriver = &fakeDriver{}

func openFake(t *testing.T, rows [][]driver.Value) *sql.DB {
	t.Helper()
	registerOnce.Do(func() { sql.Register("kei-fake", theDriver) })
	theDriver.mu.Lock()
	theDriver.rows = rows
	theDriver.mu.Unlock()
	db, err := sql.Open("kei-fake", "")
	if err != nil {
		t.Fatalf("sql.Open: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func TestSQLQuery_ColumnsRowsAndArgs(t *testing.T) {
	db := openFake(t, [][]driver.Value{{[]byte("west"), int64(10)}, {"east", int64(7)}})
	q := connection.NewSQL(db)
	cols, rows, truncated, err := q.Query(context.Background(), "select region, total from t where q = $1", []any{"2026Q3"}, 10)
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	if truncated {
		t.Fatalf("truncated = true, want false")
	}
	if names := []string{cols[0].Name, cols[1].Name}; !reflect.DeepEqual(names, []string{"region", "total"}) {
		t.Fatalf("columns = %v", names)
	}
	// []byte values come back as strings, so a result marshals as text.
	want := [][]any{{"west", int64(10)}, {"east", int64(7)}}
	if !reflect.DeepEqual(rows, want) {
		t.Fatalf("rows = %#v, want %#v", rows, want)
	}
	if theDriver.query != "select region, total from t where q = $1" || len(theDriver.args) != 1 || theDriver.args[0].Value != "2026Q3" {
		t.Fatalf("driver got %q %v", theDriver.query, theDriver.args)
	}
}

func TestSQLQuery_MaxRowsTruncates(t *testing.T) {
	db := openFake(t, [][]driver.Value{{"a", int64(1)}, {"b", int64(2)}, {"c", int64(3)}})
	q := connection.NewSQL(db)
	_, rows, truncated, err := q.Query(context.Background(), "select 1", nil, 2)
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	if len(rows) != 2 || !truncated {
		t.Fatalf("rows=%d truncated=%v, want 2 true", len(rows), truncated)
	}
	_, rows, truncated, err = q.Query(context.Background(), "select 1", nil, 3)
	if err != nil || len(rows) != 3 || truncated {
		t.Fatalf("exact fit: rows=%d truncated=%v err=%v", len(rows), truncated, err)
	}
}

func TestSQLQuery_RejectsNonPositiveMaxRows(t *testing.T) {
	db := openFake(t, nil)
	q := connection.NewSQL(db)
	if _, _, _, err := q.Query(context.Background(), "select 1", nil, 0); err == nil {
		t.Fatalf("maxRows 0: want an error")
	}
}
