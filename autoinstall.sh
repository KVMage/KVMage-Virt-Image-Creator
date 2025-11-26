#!/bin/bash

set -euo pipefail

REPO_URL="https://gitlab.com/kvmage/kvmage.git"
REPO_ROOT="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"

echo "[*] Cloning repository..."
git clone "$REPO_URL"

echo "[*] Entering repo directory..."
cd "$REPO_ROOT"

echo "[*] Creating dist directory..."
mkdir -p dist

echo "[*] Running build.sh..."
bash build.sh

echo "[*] Running install.sh..."
bash install.sh

echo "[*] Cleaning up..."
cd ..
rm -rf "$REPO_ROOT"

echo "[*] Done."
