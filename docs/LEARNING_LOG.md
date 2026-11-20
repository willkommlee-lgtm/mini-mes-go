# Learning log

One short entry per day, in English. Write it as if explaining to an interviewer.

## Day 1 — Go basics & skeleton
- **What I built:** Project skeleton with a health check, PostgreSQL running in Docker with schema and sample data.
- **New concepts:** structs and methods, maps, pointers, interfaces, goroutines and channels, graceful shutdown.
- **Compared to C# / Java:** Errors are return values, not exceptions. Interfaces are satisfied implicitly — no `implements`. Structs are copied by default; use a pointer to change the original.
- **Key insight:** Validate synchronously, send the "allow" signal first, and hand off recording to a background worker. Records may be delayed but must not be lost — that is why I will use Kafka on Day 4.

## Day 2 — REST API & scan validation
- **What I built:**
- **New concepts:**
- **Design decision and why:**

## Day 3 — Docker Compose
## Day 4 — Kafka
- **Why an event stream here instead of a direct DB write?**
- **What happens if a message is delivered twice?**

## Day 5 — Metrics & logs
## Day 6 — Kubernetes
- **What happens when a pod is killed?**

## Day 7 — Failure experiments
