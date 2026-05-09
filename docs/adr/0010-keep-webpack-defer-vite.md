# ADR-0010: Keep Webpack 5; defer Vite migration

## Status

Accepted — 2026-05-07

## Context

The original Staffjoy v2 used Webpack 1. The /autoplan plan originally proposed migrating to Vite (a non-trivial rewrite of build config). LandRover, however, already migrated `app/` and `myaccount/` to Webpack 5 + Babel 7 + React 18 — modern, supported, and not actively painful.

Vite is faster and has nicer defaults, but Vite migration on top of an already-modern Webpack 5 setup is yak-shaving when there are no users and product features are queued.

## Decision

Keep Webpack 5. Do not migrate to Vite during modernization.

W6 frontend work scope is reduced to:

- `node-sass@7` → `sass` (dart-sass) — `node-sass` is unmaintained.
- `raven-js@3` → `@sentry/react` — Sentry's old SDK is dead.
- Remove `react-hot-loader` (use Webpack 5's React Refresh).
- Replace `enzyme` with `@testing-library/react` only as needed to keep tests green under React 18 (Enzyme has no React 18 adapter).

## Consequences

- W6 effort drops from ~5 days to ~1 day.
- Build is slightly slower than Vite would be; acceptable for a foundation.
- Future Vite migration remains an option; revisit once feature work needs faster HMR or smaller bundles.
- Storybook 6.4 (Webpack-based) keeps working.
