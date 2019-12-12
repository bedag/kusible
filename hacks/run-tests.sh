#!/bin/sh

export GO111MODULE=on
# go test -race requires cgo
# export CGO_ENABLED=0

echo "Building kusible..."

GO_BUILD_CMD="go build"


mkdir -p coverage
COVERAGEDIR=$(readlink -f coverage)

go build main.go || exit 1

echo "mode: atomic" > "${COVERAGEDIR}/coverage.out"
PKGS=$(go list ./...)

fail=false
for pkg in $PKGS; do
  go test -race -coverprofile=profile.out -covermode=atomic $pkg
  if [ $? -ne 0 ]; then
    fail=true
  fi

  if [ -f profile.out ]; then
    cat profile.out | grep -v '^mode:' >> "${COVERAGEDIR}/coverage.out"
    rm profile.out
  fi
done

if [ "$fail" = true ]; then
  echo "Failure"
  exit 1
fi
