# ADR-0003: Single go.mod at repo root; no workspaces

## Status

Accepted — 2026-05-07 (inherited from LandRover)

## Context

Three layouts were considered:

- **One `go.mod` per service** — true monorepo with independent module versions. High overhead for shared `helpers/`, `crypto/`, `middlewares/`, `healthcheck/` packages; every cross-service import needs a `replace` directive.
- **Single `go.mod` at the root** (LandRover's choice) — all services share a dependency graph. Simple. One `go mod tidy` covers everything.
- **Go workspace (`go.work`) with multiple modules** — middle ground; better for independent versioning but more moving parts than this repo justifies.

## Decision

Keep LandRover's single `go.mod` at the repo root with `module v2.staffjoy.com`.

Bump the `go` directive to 1.23 in W1 (currently 1.18). Run `go mod tidy` after each dependency change in W5.0.

Delete `glide.yaml` and `glide.lock` — they are orphaned but still on disk.

## Consequences

- Every service runs against the same dependency versions; major-version bumps affect all 13 services at once. Renovate must group by ecosystem to avoid PR sprawl.
- `vendor/` is not committed (modules handle this). Build hosts need internet on first build; CI caches `~/go/pkg/mod`.
- Future split (e.g. when one service becomes a public library) is straightforward: extract its directory, add a new `go.mod`, set up `replace` directive.
