# ADR-0009: Adopt a modern design token system in W6.5

## Status

Accepted — 2026-05-07 (resolved at /autoplan W6.5 gate)

## Context

The existing `app/` (manager UI) uses Material Design Lite (`react-mdl@1.7.2`), now unmaintained and visually dated. The existing `myaccount/` uses RMWC. The modernization plan originally treated design tokens as out-of-scope.

The /autoplan Phase 2 (design) review scored the design system gap **2/10 (CRITICAL)**: shipping product features (worker login, time clock, scheduling UI) on Material Design Lite means rebuilding the visual layer 1-3 features in. The cost grows the longer it is deferred.

## Decision

Add **W6.5 — Design tokens workstream** to the modernization. Adopt one of:

- **shadcn/ui** (Radix primitives + Tailwind CSS) — most mature ecosystem, copy-paste components, easy to customize, web-only.
- **Tamagui** — cross-platform Radix-shaped library that compiles to React Native. Best long-term fit if a React Native worker app becomes part of the product roadmap.

Default recommendation: **shadcn/ui**. Reasons: smaller dependency surface, more contributors, easier to evaluate locally; the eventual mobile worker app is more likely to be a PWA + Capacitor than a React Native build, in which case shadcn applies directly.

Establish:

- `webassets/tokens.css` with CSS custom properties for color, spacing, typography, radii, shadows.
- `DESIGN.md` at repo root documenting the tokens, the chosen component library, and the contribution rules.
- `.storybook/` configuration that targets shadcn components (Storybook 6.4 already present in the base; reuse).
- Migration of one canonical screen in `app/` (the schedule view) to validate the tokens at product feature time.

## Consequences

- ~2 days CC of upfront work; saves 1-3 features' worth of rework later.
- Future product features ship on a modern, accessible component library out of the gate.
- Worker mobile app (when it lands) can re-use the tokens.
- The existing `frontend_resources/scss` and `webassets/main.scss` become legacy; they can stay for the modernization milestone and get displaced as components are rebuilt.
