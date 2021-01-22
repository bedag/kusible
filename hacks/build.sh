#!/bin/sh

export GO111MODULE=on
export CGO_ENABLED=0
APPNAME="kusible"
VERSION=${VERSION:-"development"}

echo "Building ${APPNAME}..."
# Details: https://golang.org/cmd/go/#hdr-Compile_packages_and_dependencies
GO_BUILD_CMD="go build -a -v -trimpath"
# Details: https://golang.org/cmd/link/
GO_BUILD_LDFLAGS="-s -w -X 'github.com/bedag/kusible/cmd.Version=${VERSION}'"

mkdir -p release
RELEASEDIR=$(readlink -f release)

echo ${GO_BUILD_LDFLAGS}

export GOARCH=amd64

OS=${1:-"linux"}
RACE=${2:-"false"}

export GOOS="${OS}"

case "${OS}" in
  "windows")
    EXT=".exe"
    ;;
  *)
    EXT=""
    ;;
esac

if [ "${RACE}" = "true" ]; then
  export CGO_ENABLED=1
  GO_BUILD_CMD="${GO_BUILD_CMD} -race"
  BINNAME="${APPNAME}-${GOOS}-${GOARCH}_race${EXT}"
else
  BINNAME="${APPNAME}-${GOOS}-${GOARCH}${EXT}"
fi

${GO_BUILD_CMD} -ldflags "${GO_BUILD_LDFLAGS}" -o "${RELEASEDIR}/${BINNAME}"