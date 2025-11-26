#!/bin/bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR" && git rev-parse --show-toplevel)"

DIST_DIR="${SCRIPT_DIR}/dist"
VERSION="$(cat "${REPO_ROOT}/VERSION")"

GOOS=linux  GOARCH=amd64  go build -ldflags "-X kvmage/cmd.Version=$VERSION" -o "$DIST_DIR/kvmage-linux-amd64" "$SCRIPT_DIR"
GOOS=linux  GOARCH=arm64  go build -ldflags "-X kvmage/cmd.Version=$VERSION" -o "$DIST_DIR/kvmage-linux-arm64" "$SCRIPT_DIR"
GOOS=darwin GOARCH=amd64  go build -ldflags "-X kvmage/cmd.Version=$VERSION" -o "$DIST_DIR/kvmage-darwin-amd64" "$SCRIPT_DIR"
GOOS=darwin GOARCH=arm64  go build -ldflags "-X kvmage/cmd.Version=$VERSION" -o "$DIST_DIR/kvmage-darwin-arm64" "$SCRIPT_DIR"
