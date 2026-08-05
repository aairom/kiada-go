#!/usr/bin/env bash
# ──────────────────────────────────────────────────────────────
# launch.sh — Build and deploy kiada-go to a local kind cluster
# Usage:  ./scripts/launch.sh [REGISTRY] [IMAGE_TAG]
# ──────────────────────────────────────────────────────────────
set -euo pipefail

REGISTRY="${1:-$(id -un)/kiada-go}"
IMAGE_TAG="${2:-1.0}"
IMAGE="${REGISTRY}:${IMAGE_TAG}"
CLUSTER_NAME="kiada"
PORT=30880

echo "==> Building kiada-go Docker image: ${IMAGE}"
docker build -t "${IMAGE}" ./kiada-go

echo "==> Loading image into kind cluster '${CLUSTER_NAME}'"
if ! kind get clusters 2>/dev/null | grep -q "^${CLUSTER_NAME}$"; then
  echo "==> Creating kind cluster '${CLUSTER_NAME}'"
  kind create cluster --name "${CLUSTER_NAME}"
fi
kind load docker-image "${IMAGE}" --name "${CLUSTER_NAME}"

echo "==> Patching deployment image reference"
# Replace the ${REGISTRY}/kiada-go:1.0 placeholder with the real image
sed "s|\${REGISTRY}/kiada-go:1.0|${IMAGE}|g" \
  kiada-go/k8s/02-deployment.yaml > /tmp/kiada-go-deploy-patched.yaml

echo "==> Applying Kubernetes manifests"
kubectl apply -f kiada-go/k8s/00-namespace.yaml
kubectl apply -f kiada-go/k8s/01-configmap.yaml
kubectl apply -f /tmp/kiada-go-deploy-patched.yaml
kubectl apply -f kiada-go/k8s/03-service.yaml

echo "==> Waiting for rollout to complete..."
kubectl rollout status deployment/kiada-go -n kiada --timeout=90s

URL="http://localhost:${PORT}"
echo ""
echo "────────────────────────────────────────────"
echo "  kiada-go is running!"
echo "  NodePort URL : ${URL}"
echo "  Health check : ${URL}/healthz/ready"
echo "  Pod info     : ${URL}/info"
echo "────────────────────────────────────────────"
echo ""
echo "To watch pods:"
echo "  kubectl get pods -n kiada -w"
