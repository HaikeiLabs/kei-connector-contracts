#!/usr/bin/env bash
# Local-only smoke test for the kei-connector-contracts testenv stack (HAI-998).
# NOT part of CI. Starts the stack, verifies the seed in every database,
# and tears everything down (containers AND volumes) on exit.
#
# MSSQL is opt-in: its image is amd64-only and SQL Server 2022 does not
# reliably start under amd64 emulation on arm64 Macs (needs an amd64 host).
# Set TESTENV_MSSQL=1 to include it:
#   TESTENV_MSSQL=1 bash smoke.sh
set -euo pipefail
cd "$(dirname "$0")"

COMPOSE=(podman compose)
EXPECTED="101 301 1" # accounts, issues, accounts with id above 2^53
HEALTH_BUDGET_S=180 # per-service health wait before failing (no hanging)

WITH_MSSQL=0
if [ "${TESTENV_MSSQL:-0}" = "1" ]; then
  WITH_MSSQL=1
fi

teardown() {
  echo "== teardown"
  if [ "$WITH_MSSQL" = "1" ]; then
    "${COMPOSE[@]}" --profile mssql down -v --remove-orphans
  else
    "${COMPOSE[@]}" down -v --remove-orphans
  fi
}
trap teardown EXIT

echo "== up"
# mssql-seed is run explicitly below so its exit code is checked.
if [ "$WITH_MSSQL" = "1" ]; then
  "${COMPOSE[@]}" --profile mssql up -d postgres mysql mssql supabase
else
  "${COMPOSE[@]}" up -d postgres mysql supabase
fi

state_of() {
  # Container state for a service, e.g. "running" or "exited".
  "${COMPOSE[@]}" ps -a --format '{{.State}}' "$1" 2>/dev/null | head -n 1
}

fail_if_exited() {
  # Fail immediately (with the last 20 log lines) if the service container
  # has exited, e.g. because its init/seed script failed.
  local name=$1
  if [ "$(state_of "$name")" = "exited" ]; then
    echo "FAIL $name: container exited" >&2
    "${COMPOSE[@]}" logs --tail 20 "$name" >&2 || true
    return 1
  fi
  return 0
}

wait_healthy() {
  # wait_healthy <name> <command...> — the command must exit 0 within budget.
  local name=$1
  shift
  local deadline=$((SECONDS + HEALTH_BUDGET_S))
  while [ "$SECONDS" -lt "$deadline" ]; do
    fail_if_exited "$name" || return 1
    if "$@" >/dev/null 2>&1; then
      echo "OK $name: healthy"
      return 0
    fi
    sleep 3
  done
  echo "FAIL $name: not healthy within ${HEALTH_BUDGET_S}s" >&2
  return 1
}

wait_count() {
  # wait_count <name> <command...> — the command must print "101 301 1".
  local name=$1
  shift
  local deadline=$((SECONDS + HEALTH_BUDGET_S))
  local out=""
  while [ "$SECONDS" -lt "$deadline" ]; do
    fail_if_exited "$name" || return 1
    if out=$("$@" 2>/dev/null); then
      if [ "$out" = "$EXPECTED" ]; then
        echo "OK $name: accounts=101 issues=301 big-id=1"
        return 0
      fi
    fi
    sleep 2
  done
  echo "FAIL $name: expected '$EXPECTED' within ${HEALTH_BUDGET_S}s, got '${out:-<no output>}'" >&2
  return 1
}

h_postgres() {
  "${COMPOSE[@]}" exec -T postgres pg_isready -U testenv -d testenv
}
h_supabase() {
  "${COMPOSE[@]}" exec -T supabase pg_isready -U testenv -d supabase
}
h_mysql() {
  "${COMPOSE[@]}" exec -T mysql mysqladmin ping -h localhost -u testenv -ptestenv --silent
}
h_mssql() {
  "${COMPOSE[@]}" exec -T mssql /opt/mssql-tools18/bin/sqlcmd -C -S localhost -U sa -P Testenv_local1 -Q "SELECT 1" -b -o /dev/null
}

q_postgres() {
  "${COMPOSE[@]}" exec -T postgres psql -U testenv -d testenv -At \
    -c "SELECT (SELECT count(*) FROM accounts) || ' ' || (SELECT count(*) FROM issues) || ' ' || (SELECT count(*) FROM accounts WHERE id = 9007199254740993);"
}
q_supabase() {
  "${COMPOSE[@]}" exec -T supabase psql -U testenv -d supabase -At \
    -c "SELECT (SELECT count(*) FROM accounts) || ' ' || (SELECT count(*) FROM issues) || ' ' || (SELECT count(*) FROM accounts WHERE id = 9007199254740993);"
}
q_mysql() {
  "${COMPOSE[@]}" exec -T mysql mysql -N -B -u testenv -ptestenv testenv \
    -e "SELECT CONCAT((SELECT count(*) FROM accounts), ' ', (SELECT count(*) FROM issues), ' ', (SELECT count(*) FROM accounts WHERE id = 9007199254740993));"
}
q_mssql() {
  "${COMPOSE[@]}" exec -T mssql /opt/mssql-tools18/bin/sqlcmd -C -S localhost -U sa -P Testenv_local1 -W -h -1 \
    -Q "SELECT CAST((SELECT COUNT(*) FROM accounts) AS NVARCHAR(20)) + ' ' + CAST((SELECT COUNT(*) FROM issues) AS NVARCHAR(20)) + ' ' + CAST((SELECT COUNT(*) FROM accounts WHERE id = 9007199254740993) AS NVARCHAR(20));"
}

echo "== verify"
wait_healthy postgres h_postgres
wait_count postgres q_postgres
wait_healthy supabase h_supabase
wait_count supabase q_supabase
wait_healthy mysql h_mysql
wait_count mysql q_mysql

if [ "$WITH_MSSQL" = "1" ]; then
  wait_healthy mssql h_mssql
  echo "== mssql seed (sidecar)"
  "${COMPOSE[@]}" --profile mssql run --rm mssql-seed
  wait_count mssql q_mssql
fi

echo "SMOKE OK"
