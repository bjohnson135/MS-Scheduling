# ADR-0006: Keep the 13-service architecture for the modernization; flag monolith collapse for later

## Status

Accepted — 2026-05-07

## Context

The original Staffjoy v2 architecture is 13 services: `account/{api,server}`, `company/{api,server}`, `email/server`, `sms/server`, `bot/server`, `whoami`, `ical`, `superpowers`, `faraday`, `www`, `app` (React), `myaccount` (React). They communicate via gRPC + grpc-gateway, fronted by Faraday.

For a multi-engineer team, this is sound. For a single developer building toward a When I Work replacement:

- Every gRPC service requires its own protobuf surface, separate codegen, and a gateway round-trip.
- Splitting by company vs. account doubles MySQL connection management, migrations, and ops overhead.
- A modular monolith (one Go binary, internal packages, single MySQL) ships faster, debugs faster, and decomposes back into services if/when scale demands it.

The /autoplan CEO and Eng phases both flagged this as the highest-impact strategic question. Both deferred it: doing a microservice → monolith refactor *in the same plan* as Vagrant → Docker doubles risk.

## Decision

Keep the 13-service architecture for this modernization. Containerize each service as-is. Refactor Faraday to path-based routing (ADR-0004) without changing service decomposition.

Open a separate plan ("monolith collapse") in TODOS.md to be executed *after* the modernization lands, ideally before significant new product features. Specifically: `account-server` + `company-server` + `email-server` + `sms-server` + `bot-server` + `whoami` + `ical` + `superpowers` + `faraday` could collapse into one Go binary with internal packages, with `app-service` and `myaccount-service` (frontends) and `mysql` remaining as separate compose services.

## Consequences

- Modernization stays focused; architecture risk stays bounded.
- Some Phase 0 (modernization) work — 13 separate Dockerfiles, healthcheck per service — gets thrown away when monolith collapse happens later.
- The cost is acceptable: each Dockerfile is ~30 lines from a shared template; `make` targets are parameterized; net throwaway is < 1 day of work.
- Product features (worker login, time clock) start before monolith collapse decision; if they're fast to ship on the existing decomposition, the collapse may never happen — and that is also fine.
