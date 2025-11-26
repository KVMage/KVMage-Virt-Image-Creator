#!/bin/bash

set -euo pipefail

REPO_URL="https://gitlab.com/kvmage/kvmage.git"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR" && git rev-parse --show-toplevel)"
SCRIPTS_DIR="${REPO_ROOT}/scripts"

echo "[*] Cloning repository..."
git clone "$REPO_URL"

echo "[*] Entering repo directory..."
cd "$REPO_ROOT"

echo "[*] Creating dist directory..."
mkdir -p dist

echo "[*] Running build.sh..."
bash $SCRIPTS_DIR/kvmage-build.sh

echo "[*] Running install.sh..."
bash $SCRIPTS_DIR/kvmage-install.sh

echo "[*] Cleaning up..."
cd ..
rm -rf "$REPO_ROOT"

echo "[*] Done."
