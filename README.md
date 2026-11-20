# mini-mes-go

A small Manufacturing Execution System (MES) service written in Go — a learning project.

I spent four years building and running a production MES for a high-volume tire plant (C#, Java, C++, Oracle).
One system I built there blocked wrong raw materials before they entered the mixer (barcode scan → validation → PLC interlock),
and another traced every in-scope raw material through the supply chain.

This project rebuilds simplified versions of those ideas with a stack I have **not used in production yet** —
Go, PostgreSQL, Kafka, Prometheus and Kubernetes — to learn it hands-on.

> **Status:** in progress (Day 1 of 7)

## What it does

- **Work orders & lots** — create a production order for a recipe, then production batches (lots).
- **Scan validation (mis-feed prevention)** — a scanned material is checked against the recipe: `ALLOW` or `BLOCK` with a reason.
- **Genealogy (traceability)** — list which materials went into a lot.

## Planned API

| Method | Path | Purpose |
|---|---|---|
| GET  | `/health` | Health check ✅ |
| POST | `/work-orders` | Create a work order |
| GET  | `/work-orders/{id}` | Get a work order |
| POST | `/lots` | Create a lot for a work order |
| POST | `/scans` | Validate a scanned material → `ALLOW` / `BLOCK` |
| GET  | `/lots/{id}/genealogy` | Materials used in a lot |

## Run locally

```bash
docker compose up -d          # PostgreSQL with schema and sample data
go run ./cmd/api              # http://localhost:8080
curl localhost:8080/health
```

## Roadmap

- [ ] Day 1 — Go basics, project skeleton, health check, PostgreSQL schema
- [ ] Day 2 — REST endpoints, scan validation rule, unit tests
- [ ] Day 3 — Dockerfile, run everything with Docker Compose
- [ ] Day 4 — Kafka: publish scan events, consumer stores genealogy
- [ ] Day 5 — Prometheus metrics, Grafana dashboard, structured logs
- [ ] Day 6 — deploy to Kubernetes (kind)
- [ ] Day 7 — failure experiments (Kafka down, pod killed) and write-up

## What I learned

See [docs/LEARNING_LOG.md](docs/LEARNING_LOG.md).
