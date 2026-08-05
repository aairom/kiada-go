#!/usr/bin/env bash
# ──────────────────────────────────────────────────────────────
# shutdown.sh — Gracefully remove kiada-go from the cluster
# Usage:  ./scripts/shutdown.sh [--delete-cluster]
# ──────────────────────────────────────────────────────────────
set -euo pipefail

DELETE_CLUSTER="${1:-}"

echo "==> Removing kiada-go Kubernetes resources from namespace 'kiada'"
kubectl delete -f kiada-go/k8s/03-service.yaml    --ignore-not-found
kubectl delete -f kiada-go/k8s/02-deployment.yaml --ignore-not-found
kubectl delete -f kiada-go/k8s/01-configmap.yaml  --ignore-not-found
kubectl delete -f kiada-go/k8s/00-namespace.yaml  --ignore-not-found

echo "==> Waiting for namespace 'kiada' to terminate..."
kubectl wait --for=delete namespace/kiada --timeout=60s 2>/dev/null || true

if [[ "${DELETE_CLUSTER}" == "--delete-cluster" ]]; then
  echo "==> Deleting kind cluster 'kiada'"
  kind delete cluster --name kiada
fi

echo ""
echo "kiada-go has been removed cleanly."
