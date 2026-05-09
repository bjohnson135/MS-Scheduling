# TODOS

Items deferred during modernization. Captured here so they don't get lost. Roughly ordered by priority for the When-I-Work-replacement product roadmap.

---

## Modernization tech debt (would-be-finished if scope had been larger)

### gogo/protobuf → google.golang.org/protobuf migration

**Why deferred:** the original Staffjoy v2 protos rely heavily on gogo extensions: `(gogoproto.stdtime) = true` makes `google.protobuf.Timestamp` round-trip as `time.Time`; `(gogoproto.nullable) = false` makes message fields values not pointers; `(gogoproto.moretags) = "db:\"colname\""` injects DB struct tags consumed by `gorp`. Vanilla `protoc-gen-go` has no equivalent for any of these. Migrating means:

1. Rewrite `protobuf/account.proto` and `protobuf/company.proto` to remove gogo extensions.
2. Generate via `buf` with vanilla protoc-gen-go (already wired in W3).
3. Replace every `time.Time` access on a Timestamp field with `pb.Field.AsTime()` and back via `timestamppb.New(t)` — touches scheduler, shift bounds, audit log.
4. Replace every nullable=false-implied value access with a pointer nil-check or default — touches every gorp model.
5. Move `db:"colname"` tags out of protobuf and into a separate "DB struct" pattern (or use `protoc-go-inject-tag` post-processor to keep them in `.proto`).

Effort: ~5 days CC. Recommendation: do this as part of the monolith-collapse plan (TODOS below) since both touch the model layer.

Files affected: `account/account.pb.go`, `company/company.pb.go`, `account/server/*.go`, `company/server/*.go`, `email/server/*.go`, `sms/server/*.go`, `bot/server/*.go`.

### rice → embed in `apidocs/` and `www/`

**Why deferred:** `apidocs/` uses 4 `rice.Box` vars (jsBox, cssBox, imagesBox, templatesBox) plus `rice.Box.Walk`. `www/` uses 8 boxes plus walks. Migration is straightforward but ~250 lines of changes that need runtime validation — easier to do once `docker compose up` exists (post-W2). Pattern: each `rice.MustFindBox("foo")` → `fs.Sub(assetsFS, "foo")`; each `box.String("path")` → `fs.ReadFile(subFS, "path")`; each `box.Walk` → `fs.WalkDir`.

After this lands, drop `github.com/GeertJohan/go.rice` from `go.mod`.

### apidocs/ / www/: validate that `make doctor` healthchecks confirm the migrated assets actually serve

Once compose is up: check `/swagger.json`, `/ui/`, `/`, `/login` and the sass-built CSS all load.

### Frontend test migration: enzyme → @testing-library/react

`app/` and `myaccount/` use Enzyme, which has no React 18 adapter. CI tests will silently skip Enzyme suites until rewritten. Scope: ~50-100 component tests across both apps.

### node-sass → dart-sass

`node-sass@7` is unmaintained; `sass` (dart-sass) is the canonical replacement. One-line dep swap; SCSS files compile identically. Do during W6.

### raven-js → @sentry/react; getsentry/raven-go → getsentry/sentry-go

Old SDKs are dead. Direct replacements exist. `getsentry/sentry-go` is in W5.

---

## Architecture follow-ups

### Monolith collapse plan

ADR-0006 + ADR-0008 deferred this. Once one or two product features ship on the existing decomposition, evaluate: did microservice boundaries help or hurt feature velocity? If hurt, collapse `account-server`, `company-server`, `email-server`, `sms-server`, `bot-server`, `whoami`, `ical`, `superpowers`, `faraday` into one Go binary with internal packages. Estimated 5-8 days CC.

### Auth provider externalization

JWT-signed sessions today. For commercial product, consider Auth0 or Clerk — gives social login, passkeys, MFA, audit log, all hosted. Trade: vendor lock-in vs. ~3 weeks of building these features ourselves.

### sqlc / sqlx replacing gorp

`gorp` works but is sleepy. `sqlc` generates type-safe Go from SQL files (best dev experience for a small team); `sqlx` is gorp's spiritual successor without the ORM cost. Touches every DB call in `account/server/`, `company/server/`, `auditlog/`. Defer until scaling pressure or until the gogo migration forces a model rewrite anyway.

### Production deploy story

Out of modernization scope. Options when ready: Fly.io (smallest ops surface), Railway (simplest UX), GKE (closest to original Staffjoy v2 prod), Render/Heroku-style (mid-tier). Recommend Fly.io for Phase 1 of going-live.

### MySQL → Postgres

Low priority. Existing MySQL migrations and `gorp` work. Postgres has nicer JSON, better text search, native UUID. Revisit only if a feature actively wants Postgres-only capability.

---

## Product features (the When I Work delta)

The original Staffjoy v2 explicitly does NOT have these. Each is a distinct feature with its own design + plan.

### Worker login

Workers can log in to view their own schedule, request shift swaps, set availability. Requires: worker auth flow (passwordless via SMS magic link is the modern default), worker-only routes in `app/` or a new worker-app, gRPC permission scopes for "worker reads own data only".

### Time clock / clock-in

Workers clock in / out at start of shift. Geofencing optional (Phase 2). Schema additions: `time_punches` table joining user × shift × punched_at. Manager UI to see who's late / absent.

### Worker mobile (PWA preferred over React Native)

Per /autoplan Phase 2 Pass 6. Worker UX is *the* product surface. Recommend PWA + Capacitor for v1 (web stack already in place); React Native for v2 (push, biometric) if usage shows demand.

### Team messaging beyond the existing one-way SMS bot

`bot/server/` sends one-way notifications. When-I-Work-style team chat is bi-directional, rooms, threads, mentions. Significant scope: own service, own data store (likely Postgres or Redis-backed).

### Payroll integration

Export approved hours to QuickBooks / Gusto / ADP. Each is a separate integration. v1: CSV export only.

### Reporting

Hours per worker, late minutes, no-shows, overtime warnings. New service `reports/` reading from `company` + `time_punches` DBs. Generate weekly via cron.

### Auto-scheduling algorithm

The original Staffjoy v1 had an actual scheduling solver (constraint optimization over employee availability + business demand). v2 dropped it in favor of manual scheduling. Re-introducing this is a significant feature: design the constraint inputs, pick a solver (OR-tools? custom?), build the manager UX for "solve this week."

When this lands, the `/autoplan` Phase 3 "scheduler golden-file regression test" becomes meaningful: pin baseline outputs against representative inputs. Until then, `company/server/helpers_test.go` is the closest equivalent.

---

## DX / observability gaps to revisit

### Sentry DSN real config

Modernization keeps Sentry off in dev. Wire real DSN in production deploy story.

### `make benchmark` (DX measurement)

/autoplan Phase 3.5 Pass 8 — deferred. Reconsider when contractors land on the project.

### CSRF / session secret rotation

Old keys may be in commit history. Rotate before any production deploy. Document in `docs/upgrading.md`.

### Multi-arch image publish

CI matrix should publish `linux/amd64` + `linux/arm64` to ghcr.io once W4 lands.

---

## Frontend modernization (lower priority)

### Vite migration

Webpack 5 works (ADR-0010). Consider Vite once HMR speed or bundle size becomes a real annoyance.

### Storybook component coverage

Storybook 6.4 already present. No stories committed. Once design tokens land (W6.5), seed stories for the canonical primitives (Button, Input, Modal, Schedule view).

### Astro for `www/` marketing site

Current Gulp pipeline only builds SCSS. Astro gives modern tooling at zero runtime cost. Revisit once the marketing site needs real content beyond what's there today.

### Intercom integration

Currently a no-op stub. Only wire when the product has a customer-facing support story.
