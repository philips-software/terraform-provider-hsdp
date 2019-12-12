#!/bin/sh

set -x

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
cd $SCRIPTPATH

GIT_COMMIT=$(git rev-parse --short HEAD)
go mod download

rm -rf build/ && mkdir -p build
go build -ldflags "-X main.commit=${GIT_COMMIT}" -o build/terraform-provider-hsdp .
