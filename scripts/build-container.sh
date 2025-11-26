#!/bin/bash

docker build \
  --build-arg KVMAGE_VERSION="$(cat VERSION)" \
  --build-arg BUILD_DATE="$(date -u +'%Y-%m-%dT%H:%M:%SZ')" \
  -t kvmage:$(cat VERSION) .
