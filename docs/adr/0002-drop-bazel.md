# ADR-0002: Drop Bazel; replace with go modules + Make + per-service Dockerfiles

## Status

Accepted — 2026-05-07

## Context

Bazel was orchestrating five things in this repo:

1. Per-service Go binary build.
2. Protobuf code generation with consumer-aware paths.
3. Swagger / OpenAPI generation.
4. Image layering and reproducible builds.
5. Asset embedding (via `go-bindata`).

Bazel earned its keep at multi-team scale. For a single-developer revival of a 13-service mesh with no shared C++ / no cross-language build, it adds friction at every step: every contributor must install Bazel + Java, every service needs a `BUILD` file kept in sync with imports, every CI job pre-fetches Bazel dependencies (slow on cold cache).

## Decision

Replace Bazel piece-by-piece with the simplest tool that does each job:

- (1) → `go build ./...` orchestrated by a top-level `Makefile`.
- (2) → `buf` with `buf.gen.yaml`, generating into the existing import paths.
- (3) → `buf` with `buf-gen-swagger` plugin (see W3.5).
- (4) → multi-stage `Dockerfile`s per service (3 variants: go-server, go-server+migrations, frontend) using `golang:1.23-alpine` build stage and `gcr.io/distroless/static-debian12` runtime stage.
- (5) → `embed` (Go 1.16+ stdlib).

Delete: `WORKSPACE`, all `BUILD` files (33 of them), `tools/build_rules/`, `third_party/go/*.BUILD`, `BUILD.ubuntu`, `Dockerfile.ubuntu-trusty`, `Dockerfile.ubuntu-xenial`.

## Consequences

- New contributors do not need Bazel + Java.
- CI image builds rely on GitHub Actions cache for `~/.cache/go-build` and Docker buildx layer cache; Bazel's per-target cache is lost.
- `buf` does not currently match Bazel's exact codegen output paths; W3 must verify generated files end up where existing callers expect them, or move the callers.
- Loss of Bazel's reproducible-build guarantee; acceptable trade for the ergonomic win.
