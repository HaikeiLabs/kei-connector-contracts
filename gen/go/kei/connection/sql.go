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
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Column describes one result column.
type Column struct {
	Name string
	// Type is the driver's database type name, for example "TEXT"; empty
	// when the driver does not report one.
	Type string
}

// SQL is an authenticated SQL connection. Like sqlc, it wraps a *sql.DB that
// the caller opened from a released DSN; this package never sees the DSN.
type SQL interface {
	// Query runs one statement with bind parameters and returns at most
	// maxRows rows (maxRows must be positive). truncated reports whether more
	// rows were available. []byte values are returned as strings.
	Query(ctx context.Context, sql string, args []any, maxRows int64) (columns []Column, rows [][]any, truncated bool, err error)
}

// NewSQL wraps db.
func NewSQL(db *sql.DB) SQL { return sqlConn{db: db} }

type sqlConn struct{ db *sql.DB }

func (s sqlConn) String() string   { return "connection.SQL{}" }
func (s sqlConn) GoString() string { return s.String() }

func (s sqlConn) Query(ctx context.Context, query string, args []any, maxRows int64) ([]Column, [][]any, bool, error) {
	if s.db == nil {
		return nil, nil, false, errors.New("connection: Query: nil *sql.DB")
	}
	if maxRows <= 0 {
		return nil, nil, false, fmt.Errorf("connection: Query: maxRows must be positive, got %d", maxRows)
	}
	rs, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, nil, false, fmt.Errorf("connection: Query: %w", err)
	}
	defer rs.Close()
	types, err := rs.ColumnTypes()
	if err != nil {
		return nil, nil, false, fmt.Errorf("connection: Query: columns: %w", err)
	}
	cols := make([]Column, len(types))
	for i, t := range types {
		cols[i] = Column{Name: t.Name(), Type: t.DatabaseTypeName()}
	}
	var rows [][]any
	truncated := false
	for rs.Next() {
		if int64(len(rows)) == maxRows {
			truncated = true
			break
		}
		vals := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rs.Scan(ptrs...); err != nil {
			return nil, nil, false, fmt.Errorf("connection: Query: scan: %w", err)
		}
		for i, v := range vals {
			if b, ok := v.([]byte); ok {
				vals[i] = string(b)
			}
		}
		rows = append(rows, vals)
	}
	if err := rs.Err(); err != nil {
		return nil, nil, false, fmt.Errorf("connection: Query: %w", err)
	}
	return cols, rows, truncated, nil
}
