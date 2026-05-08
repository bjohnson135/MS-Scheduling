# MS-Scheduling

A locally-runnable foundation for a modern workforce-management / shift-scheduling product, built on the Staffjoy v2 architecture.

> **Heritage:** this repo is forked from [LandRover/StaffjoyV2](https://github.com/LandRover/StaffjoyV2), which is itself a maintained fork of the original [Staffjoy/v2](https://github.com/Staffjoy/v2) (deprecated 2019). See [docs/adr/0001-base-on-landrover.md](docs/adr/0001-base-on-landrover.md) for why we started here.

## What's here

- **13 Go services** (account, company, faraday, www, etc.) — gRPC + grpc-gateway, MySQL-backed
- **2 React frontends** (`app/` for managers, `myaccount/` for users) — React 18 + Webpack 5
- **Faraday gateway** — single host, path-based routing on `localhost:8080`
- **Docker Compose** dev stack with mailhog for SMTP capture
- **Per-service multi-stage Dockerfiles**, distroless runtimes
- **GitHub Actions** CI scaffold (extend in follow-up)

## What's NOT here yet

This base modernizes Staffjoy v2; it does *not* yet deliver the things you'd need for a When I Work-style product. These are tracked in [TODOS.md](TODOS.md):

- Worker login (Staffjoy v2 explicitly never had it)
- Time clock / clock-in
- Worker mobile app (PWA recommended for v1)
- Team messaging (only one-way SMS bot exists today)
- Payroll integration
- Reporting
- Production deploy story
- Auto-scheduling solver

## Quickstart

Requirements: **[Docker](https://www.docker.com/products/docker-desktop) 24+** with `docker compose` v2, **[mise](https://mise.jdx.dev/)** (recommended) for Go/Node version pinning, plus `make` and `openssl`.

```bash
git clone https://github.com/bjohnson135/MS-Scheduling.git
cd MS-Scheduling
mise install            # installs Go 1.23.4, Node 20.18.0, pnpm, buf, etc.
make bootstrap          # copies .env, generates SIGNING_SECRET / CSRF_SECRET, checks port
make up                 # docker compose up -d --build (~5 min cold)
make doctor             # probe Faraday + every backend; non-zero exit on failure
open http://localhost:8080
```

If `make doctor` returns 200 across the board, you have a fully booted local Staffjoy stack. If a backend is unhealthy, `make logs.<service>` (e.g. `make logs.faraday`, `make logs.account-server`) tails its container.

To wipe state and start over: `make reset` (drops MySQL volume).

## Make targets

```
make help          # default; grouped command reference
make bootstrap     # one-time: .env, secrets, port check
make up            # docker compose up -d --build
make down          # docker compose down (volumes preserved)
make reset         # down -v, up -d --build, reseed
make rebuild       # rebuild every image, force-recreate
make logs          # tail every service
make logs.<svc>    # tail one service (e.g. logs.faraday)
make shell.<svc>   # exec /bin/sh in a service container
make psql          # mysql -uroot inside the mysql container
make status        # one-line per service: running / healthy / unhealthy
make doctor        # curl Faraday's aggregator; exits non-zero on failure
make build         # go build ./... + frontend bundles
make test          # go test -race -short ./...
make lint          # golangci-lint + eslint
make tidy          # go mod tidy
make proto         # buf generate (regenerate gRPC + grpc-gateway + swagger)
make images        # build every service image locally without compose up
```

## Architecture

```
HOST: localhost:8080  (single port, single binary frontline)
              │
              ▼
        ┌──────────┐
        │ faraday  │  Path-based routing (ADR-0004) + auth check
        │          │  + /health aggregator across every backend
        └─────┬────┘
   ┌──────────┼─────────┬──────────┬─────────┬───────────┬─────────┐
   ▼          ▼         ▼          ▼         ▼           ▼         ▼
 /www      /app/*   /myaccount  /api/v2/*  /api/v2/*  /whoami   /ical
 (Go)     (React)   (React)    /account/* /company/*  (Go)      (Go)
                                  │           │
                                  ▼           ▼
                            account-server  company-server
                                  │           │
                                  ▼           ▼
                              MySQL: account, company

Side services (profile-gated): email-server, sms-server, bot-server, superpowers
```

## Docs

- [docs/adr/](docs/adr/) — Architecture decision records (10 entries; the *why* behind every modernization choice)
- [TODOS.md](TODOS.md) — Deferred items, organized by priority
- [docs/troubleshooting.md](docs/troubleshooting.md) — Common failure modes keyed to fix commands

## Working with the codebase

### Adding a Go service

1. Create the service directory with a `main.go`.
2. Add an entry to `faraday/services/services.go` mapping its path prefix to its in-network hostname.
3. Add a service block to `docker-compose.yml` using the appropriate Dockerfile template (`docker/Dockerfile.go-service` for stateless services; `docker/Dockerfile.go-service-migrations` if it owns SQL migrations).
4. `make rebuild`.

### Adding a SQL migration

1. Drop a numbered `*.up.sql` and `*.down.sql` into `account/migrations/` or `company/migrations/`.
2. Restart the owning service (`docker compose restart account-server`) — the entrypoint runs pending migrations on every boot.

### Running tests selectively

```bash
go test -race -count=1 ./company/...
go test -race -count=1 -run TestSanitizeDayOfWeek ./company/server/
```

### Debugging a Go service in VS Code

`.vscode/launch.json` ships with Delve attach configs for the major services on ports 2345–2348.

## Origins + credits

The original Staffjoy v2 was authored by [@philipithomas](https://github.com/philipithomas), [@samdturner](https://github.com/samdturner), [@andhess](https://github.com/andhess), and contractors at Staffjoy circa 2016–2017. [@kootommy](https://github.com/kootommy) designed the application. The code was open-sourced when Staffjoy shut down. [@LandRover](https://github.com/LandRover) maintained the fork through 2025-06-23, modernizing Go modules, React 18, GitHub Actions, mailgun, and Ubuntu jammy.

This fork (May 2026) replaces Vagrant + Bazel with Docker Compose + a Makefile, refactors Faraday to path-based routing, bumps to Go 1.23 and Node 20, swaps `dgrijalva/jwt-go` → `golang-jwt/jwt/v5`, and migrates `go.rice` → `embed` (partial). See [docs/adr/](docs/adr/) for the full decision log.

## License

[MIT](LICENSE). Same as the upstream Staffjoy v2.
