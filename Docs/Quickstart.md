# Quickstart — kiada-go on Kubernetes

## Prerequisites

| Tool | Version | Purpose |
|------|---------|---------|
| Go | ≥ 1.21 | Build the application |
| Docker | ≥ 24 | Build & push container image |
| kind | ≥ 0.20 | Local Kubernetes cluster |
| kubectl | ≥ 1.28 | Manage the cluster |

Install kind on macOS:
```bash
brew install kind kubectl
```

---

## 1 — Clone and inspect

```bash
# The repo is already on your machine at:
ls Chapter02  Chapter05  Chapter11  kiada-go/
```

---

## 2 — Run the Go app locally (no Kubernetes needed)

```bash
cd kiada-go

# Run directly
go run .

# In a second terminal
curl http://localhost:8080/
curl http://localhost:8080/healthz/ready
curl http://localhost:8080/info
```

Expected output from `/`:
```
Hello from kiada-go v1.0!
Pod: <hostname> | Node: unknown
Pod IP: 0.0.0.0 | Node IP: 0.0.0.0
Status: Running locally
```

---

## 3 — Run tests

```bash
cd kiada-go
go test ./... -v
```

---

## 4 — Build the Docker image

```bash
cd kiada-go

# Build locally
docker build -t kiada-go:1.0 .

# Smoke-test the container
docker run --rm -p 8080:8080 kiada-go:1.0
curl http://localhost:8080/healthz/ready
```

---

## 5 — Deploy to a local kind cluster

```bash
# From the repo root
./scripts/launch.sh
```

This script will:
1. Build the Docker image
2. Create a `kind` cluster named `kiada` (if it doesn't exist)
3. Load the image directly into kind (no registry needed)
4. Apply all manifests in `kiada-go/k8s/` in order
5. Wait for the rollout and print the URL

---

## 6 — Verify the deployment

```bash
# List pods
kubectl get pods -n kiada

# Watch rollout
kubectl rollout status deployment/kiada-go -n kiada

# Test via NodePort
curl http://localhost:30880/
curl http://localhost:30880/healthz/ready
curl http://localhost:30880/info

# View logs
kubectl logs -n kiada -l app=kiada-go --tail=20

# Describe a pod to see env vars injected by Downward API (Ch 8)
kubectl describe pod -n kiada -l app=kiada-go | grep -A5 "Environment"
```

---

## 7 — Explore Kubernetes concepts from the book

### ConfigMap (Chapter 8)
```bash
kubectl get configmap kiada-go-config -n kiada -o yaml
```

### Service (Chapter 11)
```bash
kubectl get svc -n kiada
```

### Scaling (Chapter 14-15)
```bash
# Manual scale
kubectl scale deployment kiada-go -n kiada --replicas=5

# Watch scale
kubectl get pods -n kiada -w
```

### Rolling update (Chapter 15)
Edit [`kiada-go/k8s/02-deployment.yaml`](../kiada-go/k8s/02-deployment.yaml) — change the image tag to simulate a new version:
```bash
kubectl set image deployment/kiada-go kiada-go=kiada-go:1.1 -n kiada
kubectl rollout status deployment/kiada-go -n kiada
```

### Rollback (Chapter 15)
```bash
kubectl rollout undo deployment/kiada-go -n kiada
```

---

## 8 — Tear down

```bash
# Remove resources, keep the cluster
./scripts/shutdown.sh

# Remove resources AND delete the kind cluster
./scripts/shutdown.sh --delete-cluster
```

---

## Environment Variable Reference

Copy `.env.example` to `.env` for local development:
```bash
cp .env.example .env
# edit .env to set QUOTE_URL / QUIZ_URL if running the full kiada suite
```

## Architecture Overview

See [`Docs/Architecture.md`](Architecture.md) for the full Mermaid flowchart of:
- The book's 18-chapter structure
- The complete kiada application architecture
- The kiada-go component diagram
- Kubernetes resource relationships
