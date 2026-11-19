# testenv — local database test environment

A [podman compose](https://podman.io/docs) stack for exercising the
connector contracts against real databases. It runs locally with podman and
in CI: the `Database testenv` job in `.github/workflows/ci.yaml` starts the
same stack with docker on every pull request to `main`, on pushes to `main`,
and nightly. It is not a product surface: the root Go module stays
stdlib-only and imports nothing from this directory.

> **LOCAL DEVELOPMENT ONLY.** Every password in this stack is a weak
> placeholder (`testenv`, and `Testenv_local1` for MSSQL, which enforces
> password complexity). Never reuse these credentials outside this directory.

## Services

| Service  | Image                                        | Local port | Database |
| -------- | -------------------------------------------- | ---------- | -------- |
| postgres | `postgres:16.4`                              | 5432       | `testenv` |
| supabase | `postgres:16.4` (Supabase-local runs as Postgres; a plain second Postgres stands in for it) | 5433 | `supabase` |
| mysql    | `mysql:8.4`                                  | 3306       | `testenv` |
| mssql (opt-in)    | `mcr.microsoft.com/mssql/server:2022-CU15-ubuntu-22.04`, `platform: linux/amd64` | 1433 | `master` |
| mssql-seed (opt-in) | same server image (one-shot sidecar that ships `sqlcmd`) | — | — |

MSSQL is behind the compose profile `mssql` and **cannot run locally on
Apple silicon**: the image is amd64-only, and SQL Server 2022 does not
reliably start under amd64 emulation on arm64 Macs. It does not start with a
plain `podman compose up -d`, and to use it locally you need an amd64 host
(Linux x86_64). CI runs on an amd64 runner, so the `Database testenv` job
covers MSSQL on every run.

All ports bind to `127.0.0.1` only.

## Seed

`seed/` holds one hand-written synthetic schema per dialect: `accounts` and
`issues` (FK to `accounts`), 101 accounts and 301 issues. The data includes
NULLs, unicode text, and one id above 2^53 (`9007199254740993`).

Seeds load via each image's init mechanism:

- **postgres / supabase / mysql**: the seed file is mounted into
  `/docker-entrypoint-initdb.d/` and runs automatically on first start (only
  when the data volume is empty).
- **mssql**: the image has no init-script mechanism, so the `mssql-seed`
  sidecar runs `sqlcmd` once the server is healthy.

## Start

```sh
cd testenv
podman compose up -d
```

Including MSSQL (opt-in):

```sh
podman compose --profile mssql up -d
podman compose --profile mssql run --rm mssql-seed   # applies the MSSQL seed
```

Seeding happens automatically on first start for postgres, supabase and
mysql; for MSSQL the `mssql-seed` sidecar applies it once the server is
healthy.

## Verify

```sh
podman compose exec -T postgres psql -U testenv -d testenv -c "SELECT count(*) FROM accounts; SELECT count(*) FROM issues;"
podman compose exec -T supabase psql -U testenv -d supabase -c "SELECT count(*) FROM accounts;"
podman compose exec -T mysql mysql -u testenv -ptestenv testenv -e "SELECT count(*) FROM accounts; SELECT count(*) FROM issues;"
podman compose exec -T mssql /opt/mssql-tools18/bin/sqlcmd -C -S localhost -U sa -P Testenv_local1 -Q "SELECT COUNT(*) FROM accounts; SELECT COUNT(*) FROM issues;"
```

Expected: `101` accounts and `301` issues in each database, including the row
with `id = 9007199254740993`.

## Smoke test

```sh
cd testenv
bash smoke.sh
```

Starts the stack, runs the count verification against each database, and
tears everything down (containers **and** volumes) on exit, even on failure.
Each service gets a bounded health wait (180 s) so the script fails instead
of hanging. MSSQL is only included when opted in:

```sh
bash smoke.sh                  # postgres, supabase, mysql
TESTENV_MSSQL=1 bash smoke.sh  # also mssql (needs an amd64 host)
```

`smoke.sh` is a local convenience; CI runs the contract proof below instead.

## Contract proof: database.query.run

The build-tagged Go helper in `gen/go/conformance/testenv_test.go` proves the
`database.query.run` contract against Postgres, Supabase, MySQL and MSSQL: it
runs three seed queries on each service (one with NULLs, one with a bigint
above 2^53, one with NUMERIC precision), converts each result into the
`RunQueryResponse` JSON shape, and unmarshals it into the generated type. It
shells out to the client inside each container (`psql`, `mysql`, `sqlcmd`),
so the contract modules gain no database driver. The file only builds with
the `testenv` build tag, so the ordinary test jobs leave it out; the
`Database testenv` CI job runs it against all four databases.

Locally, from the repository root (Postgres, Supabase and MySQL; MSSQL needs
an amd64 host, see above):

```sh
podman compose -f testenv/compose.yaml up -d postgres supabase mysql
go -C gen/go test -tags testenv ./conformance -run 'TestEnvDatabaseQueryRun/^(postgres|supabase|mysql)$' -v
podman compose -f testenv/compose.yaml down -v --remove-orphans
```

The `-run` filter matters: the test has no skip for a missing service, so
without it the `mssql` subtests fail when MSSQL is not running. The helper
uses `podman` by default; CI sets `CONTAINER_CLI=docker`.

On an amd64 host, include MSSQL:

```sh
podman compose -f testenv/compose.yaml --profile mssql up -d postgres supabase mysql mssql mssql-seed
go -C gen/go test -tags testenv ./conformance -run TestEnvDatabaseQueryRun -v
podman compose -f testenv/compose.yaml --profile mssql down -v --remove-orphans
```

Expected: 3 subtests pass per service, and the logged JSON shows `null` for
the NULL cell, `"9007199254740993"` as a string, and
`"12345678901234567890.123456789012345678"` as a NUMERIC string.

## Teardown

Full teardown, containers **and** volumes:

```sh
cd testenv
podman compose down -v --remove-orphans
```

To re-seed from scratch, run the teardown and then `podman compose up -d`
again: the init scripts only run against empty volumes.
