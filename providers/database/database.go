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

// Package database is a read-only Postgres helper for harnesses. It opens a
// read-only handle from a DSN, caps result rows, and redacts the DSN and
// password from every error it returns. The harness supplies the DSN (via
// kei-proxy credential pass-through) and registers the Postgres driver; this
// package never reads environment variables or files itself and never imports
// kei-connector-runtime.
//
// The real guarantee is a read-only database role in the DSN: connect as a
// role that has no write grants and Postgres refuses writes at the
// permission layer, no matter what SQL is sent. Everything this package does
// is defense in depth on top of that — read-only connection startup
// parameters (so every pooled connection, including replacements, starts
// read-only) and a read-only transaction per query (so a single statement
// cannot SET the session read-write and write). It is not a substitute for a
// read-only role.
package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Config is a parsed DSN. It exposes the fields the package needs to validate
// and redact a connection; the full parameter map is kept for pass-through.
type Config struct {
	params map[string]string
}

// Parse parses a lib/pq-style key=value DSN (space-separated pairs, values
// optionally single- or double-quoted). It never reads the environment or any
// file; the DSN is the only input.
func Parse(dsn string) (Config, error) {
	if strings.TrimSpace(dsn) == "" {
		return Config{}, errors.New("database: DSN is empty")
	}
	params, err := parseDSN(dsn)
	if err != nil {
		return Config{}, err
	}
	return Config{params: params}, nil
}

// Get returns a DSN parameter value, or "" if absent.
func (c Config) Get(key string) string { return c.params[key] }

// SSLMode returns the DSN sslmode parameter, or "" if absent.
func (c Config) SSLMode() string { return c.Get("sslmode") }

// Password returns the DSN password parameter, or "" if absent.
func (c Config) Password() string { return c.Get("password") }

// parseDSN parses a DSN in either lib/pq key=value form or postgres:// URL
// form, returning the parameter map.
func parseDSN(dsn string) (map[string]string, error) {
	trimmed := strings.TrimSpace(dsn)
	if strings.HasPrefix(trimmed, "postgres://") || strings.HasPrefix(trimmed, "postgresql://") {
		return parseDSNURL(trimmed)
	}
	return parseDSNKeyValue(trimmed)
}

// parseDSNURL parses a postgres:// or postgresql:// URL into a parameter map.
func parseDSNURL(dsn string) (map[string]string, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, errors.New("database: DSN is not a valid URL")
	}
	params := map[string]string{}
	if h := u.Hostname(); h != "" {
		params["host"] = h
	}
	if p := u.Port(); p != "" {
		params["port"] = p
	}
	if u.User != nil {
		params["user"] = u.User.Username()
		if pw, ok := u.User.Password(); ok {
			params["password"] = pw
		}
	}
	if db := strings.TrimPrefix(u.Path, "/"); db != "" {
		params["dbname"] = db
	}
	for k, vs := range u.Query() {
		if len(vs) > 0 {
			params[k] = vs[0]
		}
	}
	return params, nil
}

// parseDSNKeyValue parses a space-separated key=value DSN, with values
// optionally single- or double-quoted.
func parseDSNKeyValue(dsn string) (map[string]string, error) {
	params := map[string]string{}
	i, n := 0, len(dsn)
	for i < n {
		for i < n && dsn[i] == ' ' {
			i++
		}
		if i >= n {
			break
		}
		keyStart := i
		for i < n && dsn[i] != '=' && dsn[i] != ' ' {
			i++
		}
		if i >= n || dsn[i] != '=' {
			return nil, errors.New("database: DSN contains a token without a key=value pair")
		}
		key := dsn[keyStart:i]
		i++
		if key == "" {
			return nil, errors.New("database: DSN contains an empty key")
		}
		var value string
		if i < n && (dsn[i] == '\'' || dsn[i] == '"') {
			quote := dsn[i]
			i++
			valStart := i
			for i < n && dsn[i] != quote {
				i++
			}
			if i >= n {
				return nil, fmt.Errorf("database: DSN has an unterminated quoted value for key %q", key)
			}
			value = dsn[valStart:i]
			i++
		} else {
			valStart := i
			for i < n && dsn[i] != ' ' {
				i++
			}
			value = dsn[valStart:i]
		}
		params[key] = value
	}
	return params, nil
}

// DB is a read-only Postgres handle. All queries run on a single pooled
// connection that was set to a read-only session at Open time.
type DB struct {
	db     *sql.DB
	rowCap int
}

// Option configures Open.
type Option func(*options)

type options struct {
	driver           string
	statementTimeout time.Duration
	rowCap           int
	allowInsecure    bool
}

func defaultOptions() options {
	return options{
		driver:           "postgres",
		statementTimeout: 30 * time.Second,
		rowCap:           1000,
	}
}

// WithDriver names the registered Postgres driver (default "postgres", i.e.
// lib/pq). The harness must import and register the driver.
func WithDriver(name string) Option {
	return func(o *options) { o.driver = name }
}

// WithStatementTimeout sets the per-statement timeout (default 30s).
func WithStatementTimeout(d time.Duration) Option {
	return func(o *options) { o.statementTimeout = d }
}

// WithRowCap caps the number of rows Query returns (default 1000).
func WithRowCap(n int) Option {
	return func(o *options) { o.rowCap = n }
}

// WithAllowInsecure permits a DSN with sslmode=disable. By default such a DSN
// is rejected.
func WithAllowInsecure() Option {
	return func(o *options) { o.allowInsecure = true }
}

// Open parses and validates the DSN and opens a read-only handle. Read-only
// and the statement timeout are applied as connection startup parameters, so
// every connection the pool opens — including replacements after a drop — is
// read-only, not just the first one. A DSN with sslmode=disable is rejected
// unless WithAllowInsecure is set. The returned error never contains the DSN
// or password.
func Open(ctx context.Context, dsn string, opts ...Option) (*DB, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt(&o)
	}
	cfg, err := Parse(dsn)
	if err != nil {
		return nil, err
	}
	if cfg.SSLMode() == "disable" && !o.allowInsecure {
		return nil, errors.New("database: DSN uses sslmode=disable; set AllowInsecure to override")
	}
	db, err := sql.Open(o.driver, withReadOnlyParams(dsn, o.statementTimeout))
	if err != nil {
		return nil, redact(dsn, err)
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, redact(dsn, err)
	}
	return &DB{db: db, rowCap: o.rowCap}, nil
}

// withReadOnlyParams returns dsn augmented with the read-only and
// statement-timeout connection startup parameters, in the DSN's own form
// (key=value or URL query). Every new connection therefore starts read-only.
func withReadOnlyParams(dsn string, timeout time.Duration) string {
	ms := strconv.FormatInt(timeout.Milliseconds(), 10)
	trimmed := strings.TrimSpace(dsn)
	if strings.HasPrefix(trimmed, "postgres://") || strings.HasPrefix(trimmed, "postgresql://") {
		sep := "?"
		if strings.Contains(trimmed, "?") {
			sep = "&"
		}
		return trimmed + sep + "default_transaction_read_only=on&statement_timeout=" + ms
	}
	return trimmed + " default_transaction_read_only=on statement_timeout=" + ms
}

// Query runs a read-only query and returns the rows as maps keyed by column
// name, at most db.rowCap rows. Each query runs in its own read-only
// transaction that is always rolled back, so the caller's SQL cannot lift the
// read-only mode (e.g. SET ... READ WRITE) for the duration of the query. The
// returned error never contains the DSN or password.
func (db *DB) Query(ctx context.Context, query string, args ...any) ([]map[string]any, error) {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	var out []map[string]any
	for rows.Next() {
		if len(out) >= db.rowCap {
			break
		}
		vals := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}
		row := make(map[string]any, len(cols))
		for i, c := range cols {
			row[c] = vals[i]
		}
		out = append(out, row)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// Close closes the underlying connection pool.
func (db *DB) Close() error {
	return db.db.Close()
}

// passwordParamPattern matches a password= parameter in a DSN or error string,
// whether the value is single-quoted, double-quoted, or bare.
var passwordParamPattern = regexp.MustCompile(`(?i)password=(?:'[^']*'|"[^"]*"|[^\s]+)`)

// redact returns err with the DSN's password removed. It replaces the actual
// password value and any password= parameter with a placeholder, so a driver
// error that echoes the DSN cannot leak the credential.
func redact(dsn string, err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	if cfg, perr := Parse(dsn); perr == nil {
		if pw := cfg.Password(); pw != "" {
			msg = strings.ReplaceAll(msg, pw, "***")
		}
	}
	msg = passwordParamPattern.ReplaceAllString(msg, "password=***")
	return errors.New(msg)
}
