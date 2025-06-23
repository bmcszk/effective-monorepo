# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture Overview

This is a microservices monorepo implementing a message-driven chess game system with Go services:

- **Producer Service** (`services/producer/`): HTTP REST API (port 8080) that accepts chess moves via POST `/moves` and publishes to RabbitMQ
- **Consumer Service** (`services/consumer/`): Message consumer that processes chess moves, validates game logic, and stores state in etcd
- **Shared Libraries** (`pkg/`): Common code including queue abstraction (`pkg/queue/`) and data models (`pkg/model/`)

The architecture uses RabbitMQ for async message passing and etcd for distributed state storage, deployed via Kubernetes with Helm charts.

## Development Commands

**Local Development Setup:**
```bash
make install-all-deps    # Install kind, helm, tilt dependencies
make kind-up             # Create local Kubernetes cluster
make tilt-up             # Build and deploy all services locally
make tilt-down           # Tear down local environment
make kind-down           # Delete local cluster
```

**Testing:**
```bash
make check               # Run linting, type checking and unit tests
```

**Docker Development:**
```bash
docker compose up        # Run services with RabbitMQ/etcd locally
```

**Cleanup Sidecar:**
```bash
# Enable cleanup during deployment (development only)
helm upgrade --install monorepo ./helm --set cleanup.enabled=true

# Deploy with selective cleanup
helm upgrade --install monorepo ./helm --set cleanup.enabled=true --set cleanup.targets.database=false
```

**API Testing:**
Use `test.rest` file for HTTP client testing of the producer API endpoints.

## Key Technologies

- **Language**: Go 1.23 with modules and vendoring
- **Message Queue**: RabbitMQ via ThreeDotsLabs/watermill
- **Data Store**: etcd distributed key-value store
- **Orchestration**: Kubernetes + Helm + Tilt for local development
- **Testing**: Standard Go testing + E2E tests in `/e2e/`

## Development Workflow

1. Use Tilt for rapid local development with hot reloading
2. E2E tests use BDD-style test blocks for integration testing
3. Services communicate exclusively via message queues (no direct HTTP between services)
4. Chess game state is maintained in etcd with proper move validation
5. Both services are containerized with multi-stage Docker builds for production

## Important Notes

- The main branch is `master` (not `main`)
- Services must be deployed together due to shared message queue topics
- Local development requires Kubernetes (via kind) - not just Docker Compose
- Message queue topics are configured via Helm values for different environments
