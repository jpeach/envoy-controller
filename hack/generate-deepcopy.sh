#! /usr/bin/env bash

set -o pipefail
set -o nounset
set -o errexit

readonly HERE=$(cd $(dirname "$0") && pwd)
readonly REPO=$(cd "${HERE}/.." && pwd)

cd "$REPO"

exec go run -mod=vendor sigs.k8s.io/controller-tools/cmd/controller-gen \
	object \
	paths=./api/...
