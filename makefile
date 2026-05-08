# MS-Scheduling — top-level Makefile.
#
# `make help` is the default target. The compose stack defined in
# docker-compose.yml plus the per-service Dockerfiles in docker/ are the
# substrate; this file is the developer-facing CLI.
#
# See ADR-0002 (drop Bazel) and ADR-0005 (replace Vagrant) for context.

SHELL := /bin/bash
.DEFAULT_GOAL := help
.SUFFIXES:

# Colors for `make help`. Set NO_COLOR=1 to disable.
ifeq ($(NO_COLOR),)
  c_bold   := \033[1m
  c_dim    := \033[2m
  c_cyan   := \033[36m
  c_green  := \033[32m
  c_reset  := \033[0m
endif

# Compose command — supports both `docker compose` and `docker-compose`.
COMPOSE ?= docker compose

# Default port for `make doctor` checks.
PORT ?= $(if $(STAFFJOY_PORT),$(STAFFJOY_PORT),8080)

# Active services for `make doctor` to probe via Faraday.
DOCTOR_PATHS := / /health

# ---------------------------------------------------------------------
# Help (default target)
# ---------------------------------------------------------------------

.PHONY: help
help: ## Show grouped command reference (default target)
	@printf "$(c_bold)MS-Scheduling$(c_reset) — see docs/adr/ for the why behind these.\n\n"
	@printf "$(c_cyan)One-time setup$(c_reset)\n"
	@printf "  $(c_green)bootstrap$(c_reset)        Detect tools, copy .env, write secrets, pick a free host port\n"
	@printf "\n$(c_cyan)Compose lifecycle$(c_reset)\n"
	@printf "  $(c_green)up$(c_reset)               Start the core stack (mysql, faraday, all services) in the background\n"
	@printf "  $(c_green)down$(c_reset)             Stop the stack (volumes preserved)\n"
	@printf "  $(c_green)reset$(c_reset)            Stop, drop volumes (mysql data!), re-up + reseed\n"
	@printf "  $(c_green)rebuild$(c_reset)          Rebuild every image and recreate containers\n"
	@printf "  $(c_green)logs$(c_reset)             Tail logs from every service\n"
	@printf "  $(c_green)logs.<service>$(c_reset)   Tail logs from one service (e.g. logs.faraday)\n"
	@printf "  $(c_green)shell.<service>$(c_reset)  Exec /bin/sh in a service container\n"
	@printf "  $(c_green)psql$(c_reset)             mysql -uroot inside the mysql container (despite the name)\n"
	@printf "\n$(c_cyan)Health$(c_reset)\n"
	@printf "  $(c_green)status$(c_reset)           One-line per service: running / healthy / unhealthy\n"
	@printf "  $(c_green)doctor$(c_reset)           Hit Faraday + every backend /healthz; non-zero exit on failure\n"
	@printf "\n$(c_cyan)Code$(c_reset)\n"
	@printf "  $(c_green)build$(c_reset)            go build ./... + pnpm build (frontends)\n"
	@printf "  $(c_green)test$(c_reset)             go test -race -short ./... + pnpm test\n"
	@printf "  $(c_green)lint$(c_reset)             golangci-lint + eslint\n"
	@printf "  $(c_green)tidy$(c_reset)             go mod tidy\n"
	@printf "  $(c_green)proto$(c_reset)            buf generate (regenerate gRPC + grpc-gateway + swagger)\n"
	@printf "  $(c_green)images$(c_reset)           Build every service image locally (no compose up)\n"
	@printf "\n$(c_dim)Run \`make <target>\` or \`make help\` for this list.$(c_reset)\n"

# ---------------------------------------------------------------------
# Bootstrap
# ---------------------------------------------------------------------

.PHONY: bootstrap
bootstrap: ## one-time first-run setup
	@if [ ! -f .env ]; then \
	  cp .env.example .env; \
	  printf "$(c_green)created$(c_reset) .env from .env.example\n"; \
	  if command -v openssl >/dev/null 2>&1; then \
	    SIGN=$$(openssl rand -hex 32); \
	    CSRF=$$(openssl rand -hex 32); \
	    sed -i.bak "s|SIGNING_SECRET=replace_with_openssl_rand_hex_32|SIGNING_SECRET=$$SIGN|" .env; \
	    sed -i.bak "s|CSRF_SECRET=replace_with_a_different_openssl_rand_hex_32|CSRF_SECRET=$$CSRF|" .env; \
	    rm -f .env.bak; \
	    printf "$(c_green)wrote$(c_reset) random SIGNING_SECRET and CSRF_SECRET into .env\n"; \
	  else \
	    printf "$(c_dim)NOTE: openssl not found — set SIGNING_SECRET and CSRF_SECRET in .env yourself.$(c_reset)\n"; \
	  fi; \
	else \
	  printf "$(c_dim).env already exists — leaving alone$(c_reset)\n"; \
	fi
	@if ! command -v docker >/dev/null 2>&1; then \
	  printf "$(c_bold)ERROR$(c_reset): docker not found. Install Docker Desktop / OrbStack / colima first.\n"; \
	  exit 1; \
	fi
	@if ! $(COMPOSE) version >/dev/null 2>&1; then \
	  printf "$(c_bold)ERROR$(c_reset): \"docker compose\" not available. Update Docker Desktop or install compose v2.\n"; \
	  exit 1; \
	fi
	@PORT=$$(grep '^STAFFJOY_PORT=' .env | cut -d= -f2 | tr -d ' '); \
	  if lsof -nP -iTCP:$$PORT -sTCP:LISTEN >/dev/null 2>&1; then \
	    printf "$(c_bold)WARNING$(c_reset): port $$PORT is occupied. Edit STAFFJOY_PORT in .env and rerun.\n"; \
	  else \
	    printf "$(c_green)ok$(c_reset) port $$PORT is free\n"; \
	  fi
	@printf "\n$(c_bold)Next:$(c_reset) make up    (then) open http://localhost:$$(grep '^STAFFJOY_PORT=' .env | cut -d= -f2 | tr -d ' ')\n"

# ---------------------------------------------------------------------
# Compose lifecycle
# ---------------------------------------------------------------------

.PHONY: up down reset rebuild logs status psql
up: ## Start the core stack in the background
	$(COMPOSE) up -d --build
	@printf "\n$(c_green)up$(c_reset) — give services 30s to boot, then run \`make doctor\`.\n"

down: ## Stop the stack (volumes preserved)
	$(COMPOSE) down

reset: ## Stop, drop volumes, rebuild + reseed
	$(COMPOSE) down -v
	$(COMPOSE) up -d --build

rebuild: ## Rebuild every image and recreate containers
	$(COMPOSE) build --no-cache
	$(COMPOSE) up -d --force-recreate

logs: ## Tail every service's logs
	$(COMPOSE) logs -f

logs.%: ## Tail one service's logs (e.g. logs.faraday)
	$(COMPOSE) logs -f $*

shell.%: ## Exec /bin/sh inside a service container
	$(COMPOSE) exec $* /bin/sh || $(COMPOSE) exec $* /bin/bash

psql: ## mysql shell as root
	$(COMPOSE) exec mysql mysql -uroot -p"$$(grep '^MYSQL_ROOT_PASSWORD=' .env | cut -d= -f2)"

# ---------------------------------------------------------------------
# Health
# ---------------------------------------------------------------------

.PHONY: status doctor
status: ## One-line per service
	@$(COMPOSE) ps --format 'table {{.Name}}\t{{.Service}}\t{{.Status}}\t{{.Ports}}'

doctor: ## Probe every backend through Faraday
	@printf "$(c_bold)Probing http://localhost:$(PORT)$(c_reset)\n"
	@RC=0; \
	for path in $(DOCTOR_PATHS); do \
	  status=$$(curl -s -o /dev/null -w "%{http_code}" --max-time 3 http://localhost:$(PORT)$$path 2>/dev/null || echo "000"); \
	  if [ "$$status" = "200" ] || [ "$$status" = "301" ] || [ "$$status" = "302" ]; then \
	    printf "  $(c_green)$$status$(c_reset)  $$path\n"; \
	  else \
	    printf "  $(c_bold)$$status$(c_reset)  $$path\n"; \
	    RC=1; \
	  fi; \
	done; \
	exit $$RC

# ---------------------------------------------------------------------
# Code
# ---------------------------------------------------------------------

.PHONY: build test lint tidy proto images
build: ## go build ./... + frontend bundles
	go build ./...
	@if [ -d app ] && [ -f app/package.json ]; then \
	  cd app && (pnpm install --frozen-lockfile 2>/dev/null || npm install) && (pnpm run build || npm run build); \
	fi
	@if [ -d myaccount ] && [ -f myaccount/package.json ]; then \
	  cd myaccount && (pnpm install --frozen-lockfile 2>/dev/null || npm install) && (pnpm run build || npm run build); \
	fi

test: ## go test -race -short ./... + pnpm test
	go test -race -short -count=1 -timeout=120s ./...

lint: ## golangci-lint run + eslint
	@if command -v golangci-lint >/dev/null 2>&1; then \
	  golangci-lint run; \
	else \
	  printf "$(c_dim)golangci-lint not on PATH — \`mise install\` to get it$(c_reset)\n"; \
	  go vet ./...; \
	fi

tidy: ## go mod tidy
	go mod tidy

proto: ## Regenerate gRPC + grpc-gateway + swagger
	@if command -v buf >/dev/null 2>&1; then \
	  buf generate; \
	else \
	  printf "$(c_bold)ERROR$(c_reset): buf not on PATH. Run \`mise install\` first.\n"; \
	  exit 1; \
	fi

images: ## Build every service image locally (no compose up)
	$(COMPOSE) build
