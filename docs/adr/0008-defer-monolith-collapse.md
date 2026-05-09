# ADR-0008: Monolith collapse is deferred to a separate plan

## Status

Deferred — 2026-05-07

## Context

See ADR-0006 for the architecture context. The /autoplan CEO and Eng phases independently flagged that the 13-service mesh is over-engineered for a single-developer commercial SaaS foundation. Both phases deferred the collapse to keep modernization scope bounded.

## Decision

Do not execute monolith collapse during this modernization. After modernization lands and one or two product features ship (e.g. worker login), reassess: did the microservice boundaries help or hurt feature velocity?

If they hurt, write a "monolith collapse" plan that:

1. Merges `account-server`, `company-server`, `email-server`, `sms-server`, `bot-server`, `whoami`, `ical`, `superpowers`, and `faraday` into one Go binary (e.g. `cmd/staffjoy-server/`).
2. Replaces internal gRPC calls with direct function calls.
3. Keeps `app/` and `myaccount/` frontends and `mysql` as separate compose services.
4. Removes `protobuf/`, `buf.yaml`, `buf.gen.yaml`, gRPC + grpc-gateway dependencies.
5. Deletes the per-service Dockerfile variants in favor of a single binary image.

Estimated effort: 5-8 days CC, 2 weeks human (per /autoplan CEO Section 1 "Implementation Alternatives").

## Consequences

- Modernization stays simple and ships fast.
- Product feature work begins on the existing decomposition; if it goes smoothly, monolith collapse may never happen.
- If it goes badly, the data point informs the collapse decision.
