# ADR-0004: Faraday uses path-based routing on localhost (Path A)

## Status

Accepted — 2026-05-07 (resolved at /autoplan F3.1 gate)

## Context

The original Faraday gateway routes by subdomain: `account.staffjoy-v2.local` → `accountapi-service`, `app.staffjoy-v2.local` → `app-service`, etc. Subdomains were resolved via `vagrant-hostmanager` writing `/etc/hosts` entries on the host. Under Docker Compose, host-level subdomain resolution is no longer available; two paths existed:

- **Path A — refactor Faraday to path-based routing** on `localhost:8080`. `/account/*` → `accountapi-service`, `/company/*` → `companyapi-service`, `/app/*` → `app-service`, `/myaccount/*` → `myaccount-service`, `/whoami` → `whoami-service`, `/ical/*` → `ical-service`. Touches `faraday/services.go` and any frontend code that builds absolute subdomain URLs.

- **Path B — keep subdomains, ship dnsmasq + wildcard cert** in compose. `*.staffjoy-v2.local` resolves via a `dnsmasq` container; each dev points their resolver at it. Closer to prod parity, less Go code touched, but ongoing per-dev DNS setup friction.

## Decision

**Path A.** Path-based routing on `localhost:8080`.

Rationale:
- The downstream goal is a When I Work-style product whose primary surface is a worker mobile app. Path-based APIs are what mobile clients want.
- No DNS configuration per developer; works on macOS, Linux, Windows/WSL, and behind any corporate VPN.
- Eliminates wildcard-cert generation step.

## Consequences

- `faraday/services.go` rewritten: subdomain matching → path prefix matching.
- All frontends (`app/`, `myaccount/`, `www/`) audited for absolute subdomain URL construction; URLs replaced with paths.
- Internal redirects (`Location:` headers) audited for the same.
- README quickstart works without `/etc/hosts` edits.
- "Production parity" is reduced: a real production deploy will likely use real subdomains again. Acceptable: production deploy is product-roadmap, not modernization scope.
