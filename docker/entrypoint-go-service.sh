#!/bin/sh
# Entrypoint for Go services that DON'T ship migrations. Validates required env.
# Distroless static doesn't have a shell; this script is for services that opt
# into Dockerfile.go-service-migrations or anywhere /bin/sh exists.
set -e

REQUIRED="${REQUIRED_ENV:-SIGNING_SECRET}"
for var in $REQUIRED; do
  eval "value=\${$var:-}"
  if [ -z "$value" ]; then
    echo "ERROR: $var is required (set in .env). Refusing to boot." >&2
    exit 1
  fi
done

exec /app/service "$@"
