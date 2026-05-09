# Changelog

All notable changes to this project are documented here. The format is loosely based on [Keep a Changelog](https://keepachangelog.com/), and the project follows the 4-digit `MAJOR.MINOR.PATCH.MICRO` versioning convention used by gstack.

## [1.0.0.0] - 2026-05-08

The modernized foundation. Replaces Vagrant + Bazel with Docker Compose + a Makefile. Bumps the runtime stack and rewires Faraday for path-based routing on a single host. Sets the base for a When I Work-style commercial product to be built on top.

### Added
- 10 architecture decision records in `docs/adr/` capturing every modernization choice (base on LandRover, drop Bazel, single `go.mod`, path-based Faraday routing, replace Vagrant, keep microservices for now, keep mailgun, defer monolith collapse, design tokens TBD, keep Webpack 5).
- `docker-compose.yml` brings up MySQL 8.4, mailhog, all 13 services + 2 React frontends, and Faraday on `localhost:8080` (`STAFFJOY_PORT` overridable).
- 3 multi-stage Dockerfile templates in `docker/`: generic Go service (distroless), Go service with migrations (alpine + golang-migrate), React frontend (node + nginx).
- Top-level `Makefile` with `make help` as default. Targets: `bootstrap`, `up`, `down`, `reset`, `rebuild`, `logs[.svc]`, `shell.svc`, `psql`, `status`, `doctor`, `build`, `test`, `lint`, `tidy`, `proto`, `images`.
- Faraday `/health` aggregator: probes every backend's `/health` in parallel and returns a JSON rollup. `make doctor` curls Faraday from the host.
- `tools/seed/` — Go seed tool that emits SQL inserts for a confirmed-and-active account (invite-only flow; no public signup planned).
- `tools/verify_test/` — internal password-verifier that mirrors the account-server's SELECT/Scan/CheckPasswordHash path.
- `mise.toml` pins Go 1.23.4, Node 20.18.0, pnpm 9.12.0, buf 1.45.0, golangci-lint 1.61.0, air, golang-migrate, protoc-gen-go(-grpc).
- `.editorconfig`, `lefthook.yml` (pre-commit hooks: gofumpt + golangci-lint + eslint + prettier + conventional commits), `.vscode/launch.json` Delve attach configs.
- `.env.schema` + `.env.example` (canonical env-var registry; entrypoints fail loud on missing required vars).
- Github Actions CI in `.github/workflows/ci.yml`: go test + vet, golangci-lint, frontend builds (matrix), `docker compose config` validation, `buf` proto checks (gated on `buf.yaml`).
- Validator regression suite in `company/server/helpers_test.go` (sanitizeDayOfWeek, validColor, validTimezone). Honest replacement for the originally-planned "scheduler golden test" — Staffjoy v2 has no scheduling algorithm, only a manual-shift CRUD layer.
- `docs/troubleshooting.md` with grep-able error strings (port collisions, MySQL boot, gRPC unreachable, charset migration, Apple Silicon emulation).
- Dev seed flow validated end-to-end:
  - Anonymous: `/whoami` → 403, `/app/` and `/myaccount/` → 307 to `/login/?return_to=...`
  - Login POST → 302 + JWT cookie scoped to current host (no `Domain=` attribute)
  - Authenticated: `/whoami` returns user info; `/app/` and `/myaccount/` serve the React bundles
  - Support flag (`UPDATE account SET support=1; relogin`) round-trips through DB → gRPC → JWT
  - `/logout/` clears the cookie; subsequent requests return 403/307 as anonymous

### Changed
- Go directive bumped from 1.18 to 1.23 (`go.mod` `go 1.23` + `toolchain go1.23.4`). Full repo compiles clean under Go 1.23 with the existing dependency set; no breakage observed.
- `golang-jwt/jwt/v4` → `golang-jwt/jwt/v5`. Sole call site is `crypto/sign.go`. `exp` claims switched to `jwt.NewNumericDate(...)` per v5 expectations. Both Parse calls now register signing methods via `jwt.WithValidMethods`.
- `go.rice` → `embed` for `errorpages/`, `account/api/`, `company/api/` (3 of 5 packages). `apidocs/` and `www/` still use `go.rice`; the rice→embed migration for those needs runtime validation against the compose stack and is tracked in TODOS.
- `keighl/mandrill` → `mailgun-go/v4` (inherited from LandRover; ADR-0007).
- `dgrijalva/jwt-go` → `golang-jwt/jwt` (inherited from LandRover; this PR moves to v5).
- Faraday refactored from subdomain routing (`*.staffjoy-v2.local`) to path-based routing on `localhost:8080`. New `Service.StripPrefix` flag + matched-prefix context propagation. Replaces the apex-redirect dance; introduces a longest-prefix-wins matcher pre-sorted at construction.
- Faraday `/health` is now an aggregator across all backends instead of returning `{"hello":"world"}`.
- `docker-compose.yml` service names match the gRPC `Endpoint` constants in code (`accountserver-service`, `companyserver-service`, …) so internal DNS resolution works.
- Cookie `Domain` attribute is left empty in dev (`config.Debug`) so cookies bind to the current host. Production behavior unchanged.
- Login + password-reset flows redirect to path-based URLs (`/app/`, `/myaccount/`, `/new-company/`) instead of `app.staffjoy-v2.local` etc. New `isValidReturnPath` validator replaces the old `isValidSub` subdomain validator.
- Webpack `publicPath` set to `/app/` (in `app/`) and `/myaccount/` (in `myaccount/`) so the bundle URLs route through Faraday.
- Frontend deps swapped: `node-sass@7` → `sass` (dart-sass), `raven-js@3` → `@sentry/browser@8`. Source switched from `Raven.config(...)` to `Sentry.init({dsn, environment})`.
- `dpapathanasiou/go-recaptcha` is bypassed in dev when `RECAPTCHA_PRIVATE` is unset (single call site in `www/reset.go`); production behavior unchanged.

### Removed
- Vagrantfile + entire `vagrant/` directory (12 provisioner scripts).
- `WORKSPACE`, 33 `BUILD` files, 50 `*.BUILD` third-party files, 4 `*.bzl` rules, `tools/build_rules/`, `third_party/go/`, `BUILD.ubuntu`.
- `docker/Dockerfile.ubuntu-trusty`, `docker/Dockerfile.ubuntu-xenial`, `docker/docker_pull.bzl`, `docker/Dockerfile.docker-nginx`.
- 9 Bazel-coupled CI scripts (`ci/build.sh`, `build-fmt.sh`, `test.sh`, `stage.sh`, `promote.sh`, `protobuf.sh`, `dev-build.sh`, `watch.sh`, `deploy-service.sh`).
- `glide.yaml` (orphan; modules supersede it).
- `modd.conf` (Vagrant + modd hot-reload; replaced by Docker Compose).
- The legacy `.github/workflows/ci-master.yaml` (Bazel-shaped CI; replaced by `ci.yml`).
- `ci/k8s/` kept as `ci/k8s/README.md`-flagged legacy reference — manifests don't match the current stack.

### Fixed (during QA)
- MySQL 8.4 dropped the `--default-authentication-plugin` flag; removed from compose `command:`.
- Service env-var name was `ACCOUNT_/COMPANY_MYSQL_CONFIG`; the Go services read `MYSQL_CONFIG` and append `?parseTime=true`. Compose now sets `MYSQL_CONFIG` in the plain-DSN form go-sql-driver expects.
- gRPC service DNS mismatch: code declared `Endpoint = "accountserver-service:1000"` but compose named the service `account-server`. Renamed services so internal DNS resolves.
- Login form was returning the same login page after correct creds because the `accountserver-service` hostname didn't resolve (above bug). Surfaced visibly only after that DNS fix.
- `/app/` and `/myaccount/` were blank pages: webpack `publicPath` defaulted to `/`, so the bundle was fetched from `/main-<hash>.bundle.js`, which Faraday's catch-all sent to `www-service` (404). Fixed by setting `publicPath: '/app/'` / `/myaccount/'` and adding `StripPrefix=true` to those routes.
- `@sentry/react@8` + pnpm strict layout broke webpack resolution (internal `@sentry/core` imports). Switched to `@sentry/browser` (the call site only uses `Sentry.init`) and added `.npmrc` with `shamefully-hoist=true` to keep transitive deps webpack-resolvable.
- Frontend Dockerfile only copied `${APP_PATH}/`, but webpack imports walk into `frontend_resources/` and `webassets/`. Restructured the build stage to mirror the repo layout (`/repo/${APP_PATH}` + `/repo/frontend_resources` + `/repo/webassets`).
- Cookie `Domain=staffjoy-v2.local` was hardcoded — browsers refused to set the cookie on `localhost`. Empty Domain in dev fixes it.

### Out of scope (tracked in TODOS.md)
- gogo/protobuf → vanilla `google.golang.org/protobuf` migration. Touches gorp ORM mappings + every timestamp call site. Estimated 5+ days CC; deferred to monolith-collapse plan or follow-up.
- `apidocs/` and `www/` rice → embed migration.
- Frontend test framework: Enzyme has no React 18 adapter; existing app/myaccount tests are silently skipped. Migration to `@testing-library/react` deferred.
- `make benchmark` DX measurement target.
- shadcn/ui + design-tokens workstream (ADR-0009).
- Production deploy plumbing (GKE manifests in `ci/k8s/` are legacy and don't match the current stack).
