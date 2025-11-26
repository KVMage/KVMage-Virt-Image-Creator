#!/bin/bash

set -euo pipefail

REPO_URL="https://gitlab.com/kvmage/kvmage.git"
REPO_NAME="kvmage"
SCRIPTS_DIR="scripts"

rm -rf "${REPO_NAME}"
echo "Cleaning residual artifacts"

echo "[*] Cloning repository..."
git clone "${REPO_URL}"

echo "[*] Entering repo directory..."
cd "${REPO_NAME}"

echo "[*] Creating dist directory..."
mkdir -p dist

echo "[*] Running build.sh..."
bash "${SCRIPTS_DIR}/kvmage-build.sh"

echo "[*] Running install.sh..."
bash "${SCRIPTS_DIR}/kvmage-install.sh"

echo "[*] Cleaning up..."
cd ..
rm -rf "${REPO_NAME}"

echo "[*] Done."
