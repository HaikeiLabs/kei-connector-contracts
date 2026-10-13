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

// Package integration holds the real-Postgres integration test for the
// providers/database package. It runs against a throwaway Postgres container
// (Docker or Podman) and lives in a separate Go module so the parent contract
// module stays stdlib-only. Run it from this directory with a container
// runtime available:
//
//	go test ./...
package integration

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"github.com/HaikeiLabs/kei-connector-contracts/providers/database"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// queryRetrying runs db.Query, retrying while the pool recovers from a
// terminated backend (the first query after pg_terminate_backend can hit the
// dead connection). It returns the last result and error.
func queryRetrying(t *testing.T, db *database.DB, ctx context.Context, query string) ([]map[string]any, error) {
	t.Helper()
	var rows []map[string]any
	var err error
	for i := 0; i < 10; i++ {
		rows, err = db.Query(ctx, query)
		if err == nil {
			return rows, nil
		}
		if strings.Contains(strings.ToLower(err.Error()), "read-only") {
			return rows, err
		}
		time.Sleep(150 * time.Millisecond)
	}
	return rows, err
}

// startPostgres boots a throwaway Postgres container and returns it plus a
// read-write admin *sql.DB (used to force reconnects and inspect catalog
// state) and the DSN for the read-only handle under test.
func startPostgres(t *testing.T) (*sql.DB, string, context.Context) {
	t.Helper()
	ctx := context.Background()
	pg, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.BasicWaitStrategies(),
		postgres.WithUsername("ro"),
		postgres.WithPassword("secret"),
		postgres.WithDatabase("app"),
	)
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}
	t.Cleanup(func() { _ = pg.Terminate(ctx) })

	dsn, err := pg.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("connection string: %v", err)
	}

	admin, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("open admin db: %v", err)
	}
	t.Cleanup(func() { _ = admin.Close() })
	return admin, dsn, ctx
}

func TestReadOnlyRejectsWrite(t *testing.T) {
	ctx := context.Background()
	pg, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.BasicWaitStrategies(),
		postgres.WithUsername("ro"),
		postgres.WithPassword("secret"),
		postgres.WithDatabase("app"),
	)
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}
	t.Cleanup(func() { _ = pg.Terminate(ctx) })

	dsn, err := pg.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("connection string: %v", err)
	}

	db, err := database.Open(ctx, dsn, database.WithAllowInsecure())
	if err != nil {
		t.Fatalf("open read-only db: %v", err)
	}
	defer db.Close()

	if _, err := db.Query(ctx, "SELECT 1 AS one"); err != nil {
		t.Fatalf("read query should succeed: %v", err)
	}

	_, err = db.Query(ctx, "CREATE TABLE t (id int)")
	if err == nil {
		t.Fatal("expected CREATE TABLE to fail in a read-only session")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "read-only") {
		t.Fatalf("expected a read-only error, got: %v", err)
	}
}

// TestReconnectStillReadOnly proves read-only is a property of every new
// connection, not a one-time SET on the first one. It terminates the pool's
// backend from a second admin connection, forcing database/sql to open a new
// one, and checks the replacement is read-only by default and that a write
// still fails.
func TestReconnectStillReadOnly(t *testing.T) {
	admin, dsn, ctx := startPostgres(t)

	db, err := database.Open(ctx, dsn, database.WithAllowInsecure())
	if err != nil {
		t.Fatalf("open read-only db: %v", err)
	}
	defer db.Close()

	if _, err := db.Query(ctx, "SELECT 1 AS one"); err != nil {
		t.Fatalf("read query should succeed: %v", err)
	}

	// Force a reconnect: terminate every backend except the admin's own.
	if _, err := admin.ExecContext(ctx,
		"SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE pid <> pg_backend_pid()"); err != nil {
		t.Fatalf("terminate backends: %v", err)
	}

	// The replacement connection must be read-only by default.
	rows, err := queryRetrying(t, db, ctx, "SHOW default_transaction_read_only")
	if err != nil {
		t.Fatalf("show default_transaction_read_only: %v", err)
	}
	if len(rows) != 1 || rows[0]["default_transaction_read_only"] != "on" {
		t.Fatalf("expected default_transaction_read_only=on after reconnect, got %v", rows)
	}

	// And a write still fails on the replacement connection.
	_, err = queryRetrying(t, db, ctx, "CREATE TABLE t (id int)")
	if err == nil {
		t.Fatal("expected CREATE TABLE to fail after reconnect")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "read-only") {
		t.Fatalf("expected a read-only error after reconnect, got: %v", err)
	}
}

// TestQueryCannotEscapeViaMultiStatement proves a caller cannot escape the
// read-only transaction via a multi-statement simple-protocol query. With
// lib/pq, a query with no args is sent as a single simple-protocol query,
// which Postgres allows to contain multiple statements. The string below
// commits the read-only tx, flips the session read-write, and creates a
// table. The extended protocol (prepared statements) rejects this.
func TestQueryCannotEscapeViaMultiStatement(t *testing.T) {
	admin, dsn, ctx := startPostgres(t)

	db, err := database.Open(ctx, dsn, database.WithAllowInsecure())
	if err != nil {
		t.Fatalf("open read-only db: %v", err)
	}
	defer db.Close()

	_, err = db.Query(ctx,
		"COMMIT; SET default_transaction_read_only = off; CREATE TABLE t2 (id int)")
	if err == nil {
		t.Log("query did not error; the table must still not exist")
	}

	var count int
	if err := admin.QueryRowContext(ctx,
		"SELECT count(*) FROM pg_tables WHERE schemaname='public' AND tablename='t2'").Scan(&count); err != nil {
		t.Fatalf("check pg_tables: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected table 't2' to NOT be created, but it exists (multi-statement escape succeeded)")
	}
}
