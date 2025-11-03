# Platform Service

## Project Summary

Platform Service — a simple HTTP API that aggregates data, caches responses, exposes metrics, auto-scales behind a load balancer, and is deployable via Terraform + container registry.

## Tech Stack

### Service Core

* **Go** - HTTP server, concurrency, client libs

### Glue / Automation / Tests

* **Python** - Integration tests, deployment helpers, warmup scripts

### Cache

* **Redis** - Cache-aside + TTL

### Packaging

* **Docker** - Multi-stage builds

### CI/CD

* **GitHub Actions** - Build, test, push

### Infrastructure as Code

* **Terraform** - Provision one small VM or managed container service + Redis
* **Alternative**: Render/Fly for simplicity

### Observability

* **Prometheus** client (Go) - Expose `/metrics` endpoint
* **Grafana** - Simple dashboard or use hosted metrics

### Optional

* **Nginx** or **HAProxy** - For local load balancer tests