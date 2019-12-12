#!/bin/sh

export GO111MODULE=on
export CGO_ENABLED=0

echo "Building kusible..."
# Details: https://golang.org/cmd/go/#hdr-Compile_packages_and_dependencies
GO_BUILD_CMD="go build -a -v -trimpath"
# Details: https://golang.org/cmd/link/
GO_BUILD_LDFLAGS="-s -w"

mkdir -p release
RELEASEDIR=$(readlink -f release)

${GO_BUILD_CMD} -ldflags "${GO_BUILD_LDFLAGS}" -o "${RELEASEDIR}/kusible"