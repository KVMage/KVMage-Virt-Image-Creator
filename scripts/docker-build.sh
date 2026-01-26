#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR" && git rev-parse --show-toplevel)"
VERSION_FILE="${REPO_ROOT}/VERSION"
KVMAGE_VERSION="$(< "$VERSION_FILE")"
BUILD_DATE="$(date -u +'%Y-%m-%dT%H:%M:%SZ')"
KVMAGE_BRANCH="${KVMAGE_BRANCH:-main}"

echo "[INFO] Building kvmage:${KVMAGE_VERSION} (branch: ${KVMAGE_BRANCH}) from ${REPO_ROOT}"

echo "[INFO] Local docker build"
echo "[INFO] Repo:    ${REPO_ROOT}"
echo "[INFO] Branch:  ${KVMAGE_BRANCH}"
echo "[INFO] Version: ${KVMAGE_VERSION}"

docker_build() {
  docker build \
    --progress=plain \
    --build-arg KVMAGE_VERSION="${KVMAGE_VERSION}" \
    --build-arg BUILD_DATE="${BUILD_DATE}" \
    --build-arg KVMAGE_BRANCH="${KVMAGE_BRANCH}" \
    -t "kvmage:${KVMAGE_VERSION}" \
    "${REPO_ROOT}"
}


docker_tag_latest() {
  docker tag "kvmage:${KVMAGE_VERSION}" "kvmage:latest"
}

if docker_build; then
  echo "[INFO] Docker build succeeded (no sudo)."
else
  echo "[WARN] Docker build failed without sudo. Retrying with sudo..."
  sudo -n true 2>/dev/null || echo "[INFO] sudo may prompt for a password..."
  sudo docker_build
fi

echo "[INFO] Tagging kvmage:${KVMAGE_VERSION} as latest"
if docker_tag_latest; then
  echo "[INFO] Docker tag succeeded (no sudo)."
else
  echo "[WARN] Docker tag failed without sudo. Retrying with sudo..."
  sudo docker_tag_latest
fi