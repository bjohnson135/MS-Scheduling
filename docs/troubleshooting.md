# Troubleshooting

Common failure modes and the fixes that work. Keyed for grep — search for the literal error string you saw in the logs.

## `Port 8080 already in use`

Faraday binds to `STAFFJOY_PORT` (default 8080) on the host. Edit `.env`:

```bash
STAFFJOY_PORT=8081
```

Then `make down && make up`. `make doctor` will probe the new port automatically (`make doctor PORT=8081`).

## `SIGNING_SECRET is required; refusing to boot`

The Faraday container's entrypoint and the JWT signer in `crypto/sign.go` both refuse to start without a non-empty `SIGNING_SECRET` in the environment. `make bootstrap` writes a random one into `.env` on first run.

```bash
# Regenerate manually:
SIGN=$(openssl rand -hex 32)
sed -i.bak "s|^SIGNING_SECRET=.*|SIGNING_SECRET=$SIGN|" .env
rm -f .env.bak
make down && make up
```

## `connection refused on tcp 3306` from `account-server` or `company-server`

MySQL hasn't finished its first-boot init. Compose `depends_on: condition: service_healthy` should prevent this, but if you see it:

```bash
make logs.mysql
```

Wait for `[Server] /usr/sbin/mysqld: ready for connections.` Then `make up` again — restarting the dependent service is enough.

## `Error 1273 (HY000): Unknown collation: 'utf8mb4_0900_ai_ci'`

Means you're hitting an old MySQL image (5.x or earlier 8.x). Confirm `docker-compose.yml` pins `mysql:8.4-oraclelinux9`. If you customized it: `make reset` to drop the volume, then `make up`.

## `gRPC: connection unavailable` between gateway and server

A gateway (e.g. `accountapi-service`) can't reach its `accountserver` backend. Likely the server is crashlooping. `make doctor` flags it. Then:

```bash
make logs.account-server
```

Common cause: migrations failed. The migration entrypoint logs the exact `migrate up` output before refusing to boot the service binary.

## `migrate: up failed`

A migration didn't apply cleanly to the existing schema. Most often when switching MySQL major versions. To wipe the database volume and re-run from scratch:

```bash
make reset
```

If you need to debug the migration before resetting:

```bash
make psql
mysql> use account;
mysql> SELECT * FROM schema_migrations;
mysql> exit
```

## `make doctor` reports `503` from `/health`

Faraday returned a degraded rollup. Inspect which backends are unhealthy:

```bash
curl -s http://localhost:8080/health | jq
```

The JSON body shows per-service status. Then `make logs.<service>` for the offender.

## Apple Silicon: containers boot in 5+ minutes

If you see `mysql:5.7` (or anything pre-`8.0.32`) in your `docker-compose.yml`, those have no native arm64 builds and Docker emulates via QEMU. The shipped compose pins `mysql:8.4-oraclelinux9` which is multi-arch native. If you customized it, revert.

## `pnpm: command not found`

`mise install` from the repo root reads `mise.toml` and installs pnpm 9.12.0 alongside Node 20.18.0. If you don't use mise, install pnpm manually:

```bash
corepack enable
corepack prepare pnpm@9.12.0 --activate
```

## CSRF token rejected

The CSRF secret rotated between the cookie and the validating server. Most often after running `make reset` or rotating `.env`. Clear cookies for `localhost:8080` and reload.

## Frontend bundle stale after rebuild

Browser caching. Hash-named bundle assets are `Cache-Control: immutable` for a year. The entry HTML is `no-store`, so a hard refresh (`Cmd+Shift+R` / `Ctrl+Shift+R`) is enough.

## `docker-compose: command not found`

You're on Compose v1. The Makefile defaults to `docker compose` (v2). If your tooling provides v1 only:

```bash
make COMPOSE='docker-compose' up
```

Better: upgrade Docker Desktop or install Compose v2.

## I want to start over completely

```bash
make down
docker volume rm $(docker volume ls -q | grep ms-scheduling) 2>/dev/null
docker image rm $(docker image ls --format '{{.Repository}}:{{.Tag}}' | grep ms-scheduling) 2>/dev/null
rm .env
make bootstrap
make up
```
