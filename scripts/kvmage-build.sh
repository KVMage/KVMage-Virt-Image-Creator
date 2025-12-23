#!/bin/bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR" && git rev-parse --show-toplevel)"

DIST_DIR="${REPO_ROOT}/dist"
VERSION="$(cat "${REPO_ROOT}/VERSION")"

mkdir -p "$DIST_DIR"

GOOS=linux  GOARCH=amd64  go build -ldflags "-X kvmage/cmd.Version=$VERSION -X kvmage/cmd.RequirementsB64=$REQUIREMENTS" -o "$DIST_DIR/kvmage-linux-amd64" "$REPO_ROOT"
GOOS=linux  GOARCH=arm64  go build -ldflags "-X kvmage/cmd.Version=$VERSION -X kvmage/cmd.RequirementsB64=$REQUIREMENTS" -o "$DIST_DIR/kvmage-linux-arm64" "$REPO_ROOT"
GOOS=darwin GOARCH=amd64  go build -ldflags "-X kvmage/cmd.Version=$VERSION -X kvmage/cmd.RequirementsB64=$REQUIREMENTS" -o "$DIST_DIR/kvmage-darwin-amd64" "$REPO_ROOT"
GOOS=darwin GOARCH=arm64  go build -ldflags "-X kvmage/cmd.Version=$VERSION -X kvmage/cmd.RequirementsB64=$REQUIREMENTS" -o "$DIST_DIR/kvmage-darwin-arm64" "$REPO_ROOT"
