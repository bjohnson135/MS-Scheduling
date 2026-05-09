# ADR-0007: Use mailgun-go for transactional email

## Status

Accepted — 2026-05-07 (inherited from LandRover)

## Context

The original Staffjoy v2 used `keighl/mandrill` for transactional email. Mandrill (Mailchimp Transactional) was discontinued as a standalone product in 2016 and rebranded with restricted access; the API key story is also complicated.

LandRover already migrated to `mailgun-go/v4`. The /autoplan CEO phase considered Resend, SendGrid, and Postmark as alternatives. Resend has a generous free tier and a clean Go SDK; SendGrid has the largest install base; Postmark has the best deliverability story.

## Decision

Keep `mailgun-go/v4` (LandRover's choice). Do not migrate.

Rationale: it works, the integration code already exists in `email/server/`, and the user has no email volume to ship yet. Migrating to a different provider before having users would be premature optimization. If deliverability or pricing becomes an issue once users land, swap is one file: `email/server/transport.go` becomes an interface; the mailgun implementation becomes one of multiple drivers.

For local dev, mailgun is replaced by `mailhog` (catches all SMTP, exposes a web UI at `localhost:8025`). Real mailgun keys are never required for development.

## Consequences

- One less migration to do in modernization.
- mailgun's pricing is acceptable for the scale this product will have at launch (low transactional volume).
- The `transport.go` interface refactor is deferred to product roadmap; modernization keeps the direct mailgun-go calls.
