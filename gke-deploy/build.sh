#!/bin/bash

set -e

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null && pwd)
cd "${SCRIPT_DIR}"

goimports -w .
gofmt -w .
go get -d -v
go install -v
