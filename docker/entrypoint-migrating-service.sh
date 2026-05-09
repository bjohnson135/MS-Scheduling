#!/bin/sh
# Entrypoint for Go services that ship SQL migrations.
#
# 1. Fail loud if the DSN env var is unset (per autoplan DX Pass 3).
# 2. Run any pending migrations against the DSN.
# 3. exec the service binary so PID 1 is the Go process and signals work.
set -e

if [ -z "${MIGRATE_DSN:-}" ]; then
  echo "ERROR: MIGRATE_DSN is required (set in .env). Refusing to boot." >&2
  exit 1
fi

# Convert mysql://user:pass@tcp(host:port)/db?... DSN format used by gorp into
# the format golang-migrate expects: mysql://user:pass@tcp(host:port)/db?...
# (it actually accepts both; this is just normalization for clarity).
MIGRATE_URL="${MIGRATE_DSN}"

echo "Running migrations against $(echo "$MIGRATE_URL" | sed 's|//[^:]*:[^@]*@|//***:***@|')"
migrate -path /app/migrations -database "$MIGRATE_URL" up || {
  echo "ERROR: migrations failed" >&2
  exit 1
}

echo "Migrations complete; starting service."
exec /app/service "$@"
