# ADR-0001: Base modernization on the LandRover/StaffjoyV2 fork

## Status

Accepted — 2026-05-07

## Context

The original `Staffjoy/v2` repository (last meaningful commit ~2017) is explicitly deprecated. Its README redirects to `LandRover/StaffjoyV2`, an actively-maintained fork last touched 2025-06-23. LandRover has already shipped:

- Go modules (`go.mod`, Go 1.18.2)
- React 18 in `app/` and `myaccount/` (with Webpack 5)
- GitHub Actions CI (`.github/workflows/ci-master.yaml`)
- `mailgun-go/v4` replacing the discontinued `keighl/mandrill`
- `golang-jwt/jwt/v4` replacing the deprecated `dgrijalva/jwt-go`
- `sirupsen/logrus` (lowercase) replacing the renamed `Sirupsen/logrus`
- `russross/blackfriday/v2`
- Ubuntu jammy LTS in Vagrant
- Bazel 5.1.1
- Minikube 1.25.2

LandRover does NOT replace Vagrant with Docker, does NOT remove Bazel, and is on Go 1.18 / Node 16 (both EOL).

Starting modernization from upstream `Staffjoy/v2` instead of LandRover would re-do all of the above by hand — easily 5+ engineering days for a single developer.

## Decision

Use `bjohnson135/MS-Scheduling` (a personal fork of LandRover) as the base. Inherit LandRover's modernization. Layer Vagrant→Docker, Go 1.23, Node 20, gogo→vanilla protobuf, jwt v4→v5, bindata→embed, raven-go→sentry-go, design-tokens, and CI refactor on top.

The single LandRover commit ahead of MS-Scheduling (`3ed7f5c WSL autostart`) is a Vagrant-era WSL fix that becomes irrelevant after Vagrant is removed. Do not merge.

## Consequences

- Saves ~5 days of mechanical re-do work.
- Diverges from LandRover; sending PRs back upstream becomes harder.
- LandRover's own pace (slow, ~yearly bumps) does not block our cadence.
- We lose access to LandRover's incoming security patches automatically; Renovate (W4) covers most of this.
