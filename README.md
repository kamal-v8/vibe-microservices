# Pulse — Social Feed Microservices Platform

> A polyglot microservices skeleton built for DevOps practice. Minimal business logic, maximum infrastructure potential.

```
┌──────────┐     ┌──────────────┐     ┌─────────────────────┐
│  Client  │────▶│ User Service │     │ Notification Service│
│(curl/etc)│  │  │   Go / Gin   │     │  Node.js / Fastify  │
└──────────┘  │  │    :8081     │     │       :8083         │
              │  └──────────────┘     └──────────▲──────────┘
              │                                  │ Subscribe
              │  ┌──────────────┐     ┌──────────┴──────────┐
              └─▶│ Post Service │────▶│       Redis         │
                 │Python/FastAPI│     │  Cache + Pub/Sub    │
                 │    :8082     │     └─────────────────────┘
                 └──────┬───────┘
                        │
                 ┌──────▼───────┐
                 │  PostgreSQL  │
                 │  Users+Posts │
                 └──────────────┘
```

---

## Services

| Service | Language | Framework | Port | Purpose |
|---------|----------|-----------|------|---------|
| **user-service** | Go 1.22 | Gin | 8081 | User CRUD, validates users for post-service |
| **post-service** | Python 3.12 | FastAPI | 8082 | Post CRUD, publishes events to Redis |
| **notification-service** | Node.js 20 | Fastify | 8083 | Subscribes to Redis events, stores/serves notifications |

## Data Stores

| Store | Version | Purpose |
|-------|---------|---------|
| **PostgreSQL** | 16 | Primary data store for users and posts |
| **Redis** | 7 | Pub/Sub event bus + notification storage |

---

## Quick Start

```bash
# 1. Clone and enter the project
cd pulse/

# 2. Copy environment variables
cp .env.example .env

# 3. Build and start everything
make up

# 4. Verify all services are healthy
make health

# 5. Run the end-to-end smoke test
make test-flow
```

## Available Make Targets

```
make help          Show all available targets
make up            Build and start all services
make down          Stop all services
make build         Build images without starting
make logs          Tail logs from all services
make logs-user     Tail user-service logs
make logs-post     Tail post-service logs
make logs-notif    Tail notification-service logs
make health        Check health of all services
make ps            Show running containers
make restart       Restart all services
make test-flow     Run end-to-end smoke test
make clean         Remove everything (containers, volumes, images)
```

---

## API Reference

### User Service (:8081)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/v1/users` | Create a user |
| `GET` | `/api/v1/users` | List all users |
| `GET` | `/api/v1/users/:id` | Get user by ID |
| `PUT` | `/api/v1/users/:id` | Update user |
| `DELETE` | `/api/v1/users/:id` | Delete user |
| `GET` | `/api/v1/health` | Health check |

**Create User:**

```bash
curl -X POST http://localhost:8081/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"username":"kamal","email":"kamal@pulse.dev","bio":"DevOps Engineer"}'
```

### Post Service (:8082)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/v1/posts` | Create a post (validates user, publishes event) |
| `GET` | `/api/v1/posts` | List posts (`?user_id=` filter, `?limit=&offset=`) |
| `GET` | `/api/v1/posts/:id` | Get post by ID |
| `DELETE` | `/api/v1/posts/:id` | Delete post |
| `GET` | `/api/v1/health` | Health check |

**Create Post:**

```bash
curl -X POST http://localhost:8082/api/v1/posts \
  -H "Content-Type: application/json" \
  -d '{"user_id":"<UUID>","content":"Hello from Pulse! 🚀"}'
```

### Notification Service (:8083)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/notifications` | List all notifications |
| `GET` | `/api/v1/notifications/:userId` | Get notifications for a user |
| `GET` | `/api/v1/health` | Health check |

---

## Architecture Decisions

| Decision | Choice | Why |
|----------|--------|-----|
| Inter-service calls | HTTP (REST) | Simple, debuggable, works out of the box. Upgrade to gRPC as a DevOps exercise. |
| Event bus | Redis Pub/Sub | Lightweight, zero config. Swap for Kafka/NATS when you need persistence/replay. |
| Database | Shared PostgreSQL | Single instance, separate tables per service. Split into per-service DBs as you scale. |
| No ORM | Raw SQL | Keeps dependencies minimal, Dockerfiles lighter, and makes the DB layer transparent. |

## DevOps-Ready Hooks

These are baked into every service to make DevOps integration seamless:

- ✅ **Health checks** — `/api/v1/health` with dependency verification (DB ping, Redis ping)
- ✅ **Structured JSON logging** — Ready for ELK/Loki/Fluentd ingestion
- ✅ **Environment config** — 12-factor app, ready for K8s ConfigMaps/Secrets
- ✅ **Graceful shutdown** — Handles SIGTERM (K8s pod termination lifecycle)
- ✅ **Metrics stub** — `/metrics` endpoint ready for Prometheus exporters
- ✅ **Intentionally un-optimized Dockerfiles** — You'll improve these

## What You'll Add (DevOps Roadmap)

```
□ Multi-stage + Distroless Docker builds
□ GitHub Actions / Jenkins CI/CD pipelines
□ Kubernetes manifests (Deployments, Services, Ingress)
□ Helm charts with values.yaml
□ Terraform (Kind/EKS + RDS + ElastiCache)
□ ArgoCD GitOps (App-of-Apps pattern)
□ Prometheus + Grafana monitoring stack
□ ELK / Loki logging pipeline
□ Trivy / Snyk security scanning
□ Vault / Sealed Secrets for secret management
□ Service mesh (Istio / Linkerd) — optional
□ HPA autoscaling based on custom metrics
```

---

## Project Structure

```
pulse/
├── docker-compose.yaml          # Local orchestration
├── Makefile                     # Convenience commands
├── .env.example                 # Environment template
├── .gitignore
├── README.md
├── scripts/
│   └── init-db.sql              # PostgreSQL initialization
└── services/
    ├── user-service/            # Go (Gin) — :8081
    ├── post-service/            # Python (FastAPI) — :8082
    └── notification-service/    # Node.js (Fastify) — :8083
```

---

## License

This is a learning project. Use it however you want.
