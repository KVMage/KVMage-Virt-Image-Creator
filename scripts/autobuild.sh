#!/bin/bash
set -euo pipefail

REPO_URL="https://gitlab.com/kvmage/kvmage.git"
REPO_NAME="kvmage"
KVMAGE_BRANCH="${KVMAGE_BRANCH:-main}"
SCRIPTS_DIR="scripts"

echo "[*] Cleaning residual artifacts"
rm -rf "${REPO_NAME}"

echo "[*] Cloning repository (branch: ${KVMAGE_BRANCH})..."
git clone --branch "${KVMAGE_BRANCH}" --single-branch "${REPO_URL}" "${REPO_NAME}"

echo "[*] Entering repo directory..."
cd "${REPO_NAME}"

echo "[*] Confirming branch..."
echo "[*] Branch: $(git rev-parse --abbrev-ref HEAD)"
echo "[*] Commit: $(git rev-parse --short HEAD)"

echo "[*] Running docker-build.sh..."
KVMAGE_BRANCH="${KVMAGE_BRANCH}" bash "${SCRIPTS_DIR}/docker-build.sh"

echo "[*] Cleaning up..."
cd ..
rm -rf "${REPO_NAME}"

echo "[*] Done."
