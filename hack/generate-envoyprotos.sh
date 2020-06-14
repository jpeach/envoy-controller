#! /usr/bin/env bash

set -o pipefail
set -o nounset
set -o errexit

readonly HERE=$(cd $(dirname "$0") && pwd)
readonly REPO=$(cd "${HERE}/.." && pwd)

readonly CONTROLPLANE="${REPO}/vendor/github.com/envoyproxy/go-control-plane"

exec > ${REPO}/pkg/xds/register.go

(
	echo package xds

	cat <<EOF
// Import all the Envoy API packages from go-control-plane for their
// side-effects. This causes all their protobuf types to be registered.

// The import list can be generated with the following script:
//
// $0
EOF

	echo "import ("

	go list -mod=readonly github.com/envoyproxy/go-control-plane/envoy/...  |
		sort | \
		awk '{printf "_ \"%s\"\n", $1}'

	echo ")"
) | gofmt
