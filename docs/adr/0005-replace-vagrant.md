# ADR-0005: Replace Vagrant + VirtualBox + Minikube with Docker Compose

## Status

Accepted — 2026-05-07

## Context

The dev environment as inherited from LandRover requires:

- Vagrant + VirtualBox (or Hyper-V/WSL on Windows).
- A 6 GB / 2 vCPU Ubuntu jammy VM provisioned by `vagrant/provision.sh` (~10 minutes first run).
- Minikube inside the VM running a single-node Kubernetes cluster.
- Bazel 5.1.1 + Go 1.18 + Node 16 + Docker + nginx + MySQL 5.x + grpc tools, all installed by `vagrant/*.sh` scripts.
- Unison file sync between host and VM.

This stack does not run on Apple Silicon without significant pain: VirtualBox arm64 support is experimental and unreliable. WSL2 on Windows works but adds complexity. macOS users on Intel are dwindling.

## Decision

Replace the entire VM-plus-Minikube setup with Docker Compose. Deliverables:

- `docker-compose.yml` at repo root, defining: `mysql:8.4-oraclelinux9`, `mailhog`, every Go service, every frontend service, `faraday`.
- Per-service multi-stage `Dockerfile` (3 variants — see ADR-0002).
- `make bootstrap`, `make up`, `make down`, `make doctor`, `make status`, `make seed`, `make reset`, `make logs.<service>`, `make shell.<service>` (see ADR-0009 placeholder for full Make surface).
- Local file mounts for hot-reload via `air` (Go) and Webpack/Vite dev server (frontend).
- Delete: `Vagrantfile`, `vagrant/`, `vagrant/provision.sh`, `vagrant/dev-watch.sh`, `vagrant/unison.sh`, all `.sh` scripts that target the VM provisioner.

## Consequences

- TTHW (Time to Hello World) target: ~6 minutes cold, ~2 minutes warm.
- Apple Silicon native (`mysql:8.4-oraclelinux9` has arm64 builds; distroless static is multi-arch).
- Loss of Kubernetes parity for local dev (Minikube ran `kubectl` workflows). Production-style YAML manifests will need a separate `kind`/`k3d` profile if/when desired; out of this scope.
- Loss of Bazel-orchestrated image building (see ADR-0002).
